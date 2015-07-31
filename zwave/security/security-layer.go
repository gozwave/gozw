package security

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/bjyoungblood/gozw/zwave/command-class"
)

const (
	InternalNonceTTL    = time.Second * 15
	ExternalNonceTTL    = time.Second * 10
	NonceRequestTimeout = time.Second * 10
)

type ISecurityLayer interface {
	DecryptMessage(data *commandclass.SecurityMessageEncapsulation) ([]byte, error)
	EncapsulateMessage(payload []byte, senderNonce []byte, receiverNonce []byte, srcNode byte, dstNode byte, inclusionMode bool) []byte
	GenerateInternalNonce() (Nonce, error)
	GetExternalNonce(key byte) (Nonce, error)
	ReceiveNonce(fromNode uint8, data *commandclass.SecurityNonceReport)
	WaitForExternalNonce(nodeId uint8) (Nonce, error)
}

type SecurityLayer struct {
	// internal nonce table is keyed by the first byte of the nonce
	internalNonceTable *NonceTable

	// external nonce table is keyed by the node id
	externalNonceTable *NonceTable

	// maps node id to channel
	waitForNonce map[byte]chan bool
	waitMapLock  *sync.Mutex
}

func NewSecurityLayer() *SecurityLayer {
	securityLayer := &SecurityLayer{
		internalNonceTable: NewNonceTable(),
		externalNonceTable: NewNonceTable(),

		waitForNonce: map[byte]chan bool{},
		waitMapLock:  &sync.Mutex{},
	}

	return securityLayer
}

func (s *SecurityLayer) EncapsulateMessage(
	payload []byte,
	senderNonce []byte,
	receiverNonce []byte,
	srcNode byte,
	dstNode byte,
	inclusionMode bool,
) []byte {

	var encKey, authKey []byte
	if inclusionMode {
		encKey = InclusionEncKey
		authKey = InclusionAuthKey
	} else {
		encKey = NetworkEncKey
		authKey = NetworkAuthKey
	}

	iv := append(senderNonce, receiverNonce...)

	encryptedPayload := CryptMessage(payload, iv, encKey)

	authDataBuf := append(iv, commandclass.CommandSecurityMessageEncapsulation) // @todo CC should be determined by sequencing
	authDataBuf = append(authDataBuf, srcNode)                                  // sender node
	authDataBuf = append(authDataBuf, dstNode)                                  // receiver node
	authDataBuf = append(authDataBuf, uint8(len(encryptedPayload)))
	authDataBuf = append(authDataBuf, encryptedPayload...)

	hmac := CalculateHMAC(authDataBuf, authKey)

	return commandclass.NewSecurityMessageEncapsulation(
		senderNonce,
		encryptedPayload,
		hmac,
		receiverNonce[0],
	)
}

// @todo verify message hmac
func (s *SecurityLayer) DecryptMessage(data *commandclass.SecurityMessageEncapsulation) ([]byte, error) {
	receiverNonce, err := s.internalNonceTable.Get(data.ReceiverNonceId)
	if err != nil {
		return nil, err
	}

	senderNonce := make([]byte, 8)
	copy(senderNonce, data.SenderNonce)
	iv := append(senderNonce, receiverNonce...)

	pl := make([]byte, len(data.EncryptedPayload))
	copy(pl, data.EncryptedPayload)
	decryptedPayload := CryptMessage(pl, iv, NetworkEncKey)

	return decryptedPayload[1:], nil
}

// GenerateInternalNonce returns a new internal nonce and stores it in the
// internal nonce table.
//
// NOTE: The Z-Wave docs are not very clear on this, but the "receiver nonce id"
// is simply the first byte of the nonce (which must be unique among all of the
// active internal nonces)
func (s *SecurityLayer) GenerateInternalNonce() (Nonce, error) {
	return s.internalNonceTable.Generate(InternalNonceTTL)
}

func (s *SecurityLayer) GetExternalNonce(key byte) (Nonce, error) {
	return s.externalNonceTable.Get(key)
}

// ReceiveNonce stores the received nonce in the external nonce table. Additionally,
// it sets a timeout on the nonce (after which the nonce will be deleted from the
// nonce table) and notifies any goroutine that may be waiting for a nonce from
// the given node
func (s *SecurityLayer) ReceiveNonce(fromNode uint8, data *commandclass.SecurityNonceReport) {
	s.externalNonceTable.Set(fromNode, data.Nonce, ExternalNonceTTL)

	// if there is no matching channel in the waitForNonce map, then apparently we
	// either fetched the nonce for no reason, some node just randomly gave us one,
	// or whatever process requested the nonce timed out already. in any case, we've
	// stored the nonce, so it'll be valid for now
	if ch, ok := s.waitForNonce[fromNode]; ok {

		// perform a non-blocking send. it's possible that some process asked for the
		// nonce, but for whatever reason didn't bother to listen on the channel. in
		// any case, we never want to block here
		select {
		case ch <- true:
		default:
		}

		// closing the channel will unblock anything that is currently listening (in
		// the case of multiple listeners or future listeners). this is especially
		// important if two goroutines are waiting on the channel value, since we only
		// emit one time. this also helps us in the case that some rogue goroutine has
		// a reference to this channel that it hasn't blocked on yet, but will at some
		// point
		close(ch)

		// delete the channel from the map; a new one will be created when we request
		// another nonce from the node
		delete(s.waitForNonce, fromNode)
	}
}

func (s *SecurityLayer) WaitForExternalNonce(nodeId uint8) (Nonce, error) {
	var waitChan chan bool
	var ok bool

	// Get the wait channel, creating it if it doesn't exist (note the !ok condition)
	s.waitMapLock.Lock()
	if waitChan, ok = s.waitForNonce[nodeId]; !ok {
		waitChan = make(chan bool)
		s.waitForNonce[nodeId] = waitChan
	}
	s.waitMapLock.Unlock()
	runtime.Gosched()

	select {
	case <-waitChan:
	case <-time.After(NonceRequestTimeout):
		return nil, errors.New("nonce timeout")
	}

	nonce, err := s.externalNonceTable.Get(nodeId)
	if err == nil {
		return nonce, nil
	}

	return nil, errors.New("Failed to get external nonce")
}
