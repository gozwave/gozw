package security

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/helioslabs/gozw/zwave/command-class"
)

const (
	internalNonceTTL    = time.Second * 15
	externalNonceTTL    = time.Second * 10
	nonceRequestTimeout = time.Second * 10
)

type ILayer interface {
	DecryptMessage(data *commandclass.SecurityMessageEncapsulation) ([]byte, error)
	EncapsulateMessage(payload []byte, senderNonce []byte, receiverNonce []byte, srcNode byte, dstNode byte, inclusionMode bool) []byte
	GenerateInternalNonce() (Nonce, error)
	GetExternalNonce(key byte) (Nonce, error)
	ReceiveNonce(fromNode byte, data *commandclass.SecurityNonceReport)
	WaitForExternalNonce(nodeID byte) (Nonce, error)
}

type Layer struct {
	networkKey []byte

	networkEncKey  []byte
	networkAuthKey []byte

	// internal nonce table is keyed by the first byte of the nonce
	internalNonceTable *NonceTable

	// external nonce table is keyed by the node id
	externalNonceTable *NonceTable

	// maps node id to channel
	waitForNonce map[byte]chan bool
	waitMapLock  *sync.Mutex
}

func NewLayer(networkKey []byte) *Layer {
	securityLayer := &Layer{
		networkKey: networkKey,

		networkEncKey:  EncryptEBS(networkKey, encryptPassword),
		networkAuthKey: EncryptEBS(networkKey, authPassword),

		internalNonceTable: NewNonceTable(),
		externalNonceTable: NewNonceTable(),

		waitForNonce: map[byte]chan bool{},
		waitMapLock:  &sync.Mutex{},
	}

	return securityLayer
}

func (s *Layer) EncapsulateMessage(
	payload []byte,
	senderNonce []byte,
	receiverNonce []byte,
	srcNode byte,
	dstNode byte,
	inclusionMode bool,
) []byte {

	var encKey, authKey []byte
	if inclusionMode {
		encKey = inclusionEncKey
		authKey = inclusionAuthKey
	} else {
		encKey = s.networkEncKey
		authKey = s.networkAuthKey
	}

	iv := append(senderNonce, receiverNonce...)

	encryptedPayload := CryptMessage(payload, iv, encKey)

	authDataBuf := append(iv, commandclass.CommandSecurityMessageEncapsulation) // @todo CC should be determined by sequencing
	authDataBuf = append(authDataBuf, srcNode)                                  // sender node
	authDataBuf = append(authDataBuf, dstNode)                                  // receiver node
	authDataBuf = append(authDataBuf, byte(len(encryptedPayload)))
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
func (s *Layer) DecryptMessage(data *commandclass.SecurityMessageEncapsulation) ([]byte, error) {
	receiverNonce, err := s.internalNonceTable.Get(data.ReceiverNonceID)
	if err != nil {
		return nil, err
	}

	senderNonce := make([]byte, 8)
	copy(senderNonce, data.SenderNonce)
	iv := append(senderNonce, receiverNonce...)

	pl := make([]byte, len(data.EncryptedPayload))
	copy(pl, data.EncryptedPayload)
	decryptedPayload := CryptMessage(pl, iv, s.networkEncKey)

	return decryptedPayload[1:], nil
}

// GenerateInternalNonce returns a new internal nonce and stores it in the
// internal nonce table.
//
// NOTE: The Z-Wave docs are not very clear on this, but the "receiver nonce id"
// is simply the first byte of the nonce (which must be unique among all of the
// active internal nonces)
func (s *Layer) GenerateInternalNonce() (Nonce, error) {
	return s.internalNonceTable.Generate(internalNonceTTL)
}

func (s *Layer) GetExternalNonce(key byte) (Nonce, error) {
	return s.externalNonceTable.Get(key)
}

// ReceiveNonce stores the received nonce in the external nonce table. Additionally,
// it sets a timeout on the nonce (after which the nonce will be deleted from the
// nonce table) and notifies any goroutine that may be waiting for a nonce from
// the given node
func (s *Layer) ReceiveNonce(fromNode byte, data *commandclass.SecurityNonceReport) {
	s.externalNonceTable.Set(fromNode, data.Nonce, externalNonceTTL)

	// if there is no matching channel in the waitForNonce map, then apparently we
	// either fetched the nonce for no reason, some node just randomly gave us one,
	// or whatever process requested the nonce timed out already. in any case, we've
	// stored the nonce, so it'll be valid for now
	s.waitMapLock.Lock()
	ch, ok := s.waitForNonce[fromNode]
	s.waitMapLock.Unlock()
	runtime.Gosched()

	if ok {

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
		s.waitMapLock.Lock()
		delete(s.waitForNonce, fromNode)
		s.waitMapLock.Unlock()
		runtime.Gosched()
	}
}

func (s *Layer) WaitForExternalNonce(nodeID byte) (Nonce, error) {
	var waitChan chan bool
	var ok bool

	// Get the wait channel, creating it if it doesn't exist (note the !ok condition)
	s.waitMapLock.Lock()
	if waitChan, ok = s.waitForNonce[nodeID]; !ok {
		waitChan = make(chan bool)
		s.waitForNonce[nodeID] = waitChan
	}
	s.waitMapLock.Unlock()
	runtime.Gosched()

	select {
	case <-waitChan:
	case <-time.After(nonceRequestTimeout):
		return nil, errors.New("nonce timeout")
	}

	nonce, err := s.externalNonceTable.Get(nodeID)
	if err == nil {
		return nonce, nil
	}

	return nil, errors.New("Failed to get external nonce")
}
