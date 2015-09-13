package security

import (
	"errors"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/comail/colog"
	"github.com/helioslabs/gozw/command-class/security"
	"github.com/helioslabs/gozw/serial-api"
)

const (
	internalNonceTTL    = time.Second * 15
	externalNonceTTL    = time.Second * 10
	nonceRequestTimeout = time.Second * 10
)

type ILayer interface {
	DecryptMessage(cmd serialapi.ApplicationCommand, inclusionMode bool) ([]byte, error)
	EncapsulateMessage(srcNode byte, dstNode byte, commandID security.CommandID, senderNonce []byte, receiverNonce []byte, payload []byte, inclusionMode bool) (*EncryptedMessage, error)
	GenerateInternalNonce() (Nonce, error)
	GetExternalNonce(key byte) (Nonce, error)
	ReceiveNonce(fromNode byte, report security.NonceReport)
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

	logger *log.Logger
}

func NewLayer(networkKey []byte) *Layer {
	securityLogger := colog.NewCoLog(os.Stdout, "security ", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	securityLogger.ParseFields(true)

	securityLayer := &Layer{
		networkKey: networkKey,

		networkEncKey:  EncryptEBS(networkKey, encryptPassword),
		networkAuthKey: EncryptEBS(networkKey, authPassword),

		internalNonceTable: NewNonceTable(),
		externalNonceTable: NewNonceTable(),

		waitForNonce: map[byte]chan bool{},
		waitMapLock:  &sync.Mutex{},

		logger: securityLogger.NewLogger(),
	}

	return securityLayer
}

func (s *Layer) EncapsulateMessage(
	srcNode byte,
	dstNode byte,
	commandID security.CommandID,
	senderNonce []byte,
	receiverNonce []byte,
	payload []byte,
	inclusionMode bool,
) (*EncryptedMessage, error) {

	var encKey, authKey []byte
	if inclusionMode {
		s.logger.Print("debug: encrypting message using inclusion encryption")
		encKey = inclusionEncKey
		authKey = inclusionAuthKey
	} else {
		s.logger.Print("debug: encrypting message using network encryption")
		encKey = s.networkEncKey
		authKey = s.networkAuthKey
	}

	iv := append(senderNonce, receiverNonce...)

	encryptedPayload := CryptMessage(payload, iv, encKey)

	authDataBuf := append(iv, byte(security.CommandMessageEncapsulation))
	authDataBuf = append(authDataBuf, srcNode) // sender node
	authDataBuf = append(authDataBuf, dstNode) // receiver node
	authDataBuf = append(authDataBuf, byte(len(encryptedPayload)))
	authDataBuf = append(authDataBuf, encryptedPayload...)

	hmac := CalculateHMAC(authDataBuf, authKey)

	return &EncryptedMessage{
		SenderNonce:      senderNonce,
		EncryptedPayload: encryptedPayload,
		ReceiverNonceID:  receiverNonce[0],
		HMAC:             hmac,
	}, nil
}

// @todo verify message hmac
func (s *Layer) DecryptMessage(cmd serialapi.ApplicationCommand, inclusionMode bool) ([]byte, error) {
	var encKey /*, authKey*/ []byte
	if inclusionMode {
		s.logger.Print("debug: decrypting message using inclusion encryption")
		encKey = inclusionEncKey
		// authKey = inclusionAuthKey
	} else {
		s.logger.Print("debug: decrypting message using network encryption")
		encKey = s.networkEncKey
		// authKey = s.networkAuthKey
	}

	message := EncryptedMessage{}
	err := message.UnmarshalBinary(cmd.CommandData)
	if err != nil {
		return nil, err
	}

	receiverNonce, err := s.internalNonceTable.Get(message.ReceiverNonceID)
	if err != nil {
		return nil, err
	}

	senderNonce := make([]byte, 8)
	copy(senderNonce, message.SenderNonce)
	iv := append(senderNonce, receiverNonce...)

	return CryptMessage(message.EncryptedPayload, iv, encKey), nil
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
func (s *Layer) ReceiveNonce(fromNode byte, report security.NonceReport) {
	s.externalNonceTable.Set(fromNode, report.NonceByte, externalNonceTTL)

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
			s.logger.Print("debug: received nonce and notified waiting channel")
		default:
			s.logger.Print("debug: no channel was waiting for received nonce")
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
		s.logger.Printf("warn: external nonce timeout node=%d", nodeID)
		return nil, errors.New("nonce timeout")
	}

	nonce, err := s.externalNonceTable.Get(nodeID)
	if err == nil {
		return nonce, nil
	}

	return nil, errors.New("Failed to get external nonce")
}
