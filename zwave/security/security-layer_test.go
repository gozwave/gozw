package security

import (
	"testing"
	"time"

	"github.com/helioslabs/gozw/zwave/command-class/security"
	"github.com/stretchr/testify/assert"
)

var testNetworkKey = []byte{
	0x01,
	0x02,
	0x03,
	0x04,
	0x05,
	0x06,
	0x07,
	0x08,
	0x09,
	0x0A,
	0x0B,
	0x0C,
	0x0D,
	0x0E,
	0x0F,
	0x10,
}

func TestSecurityLayerNonceGeneration(t *testing.T) {
	securityLayer := NewLayer(testNetworkKey)
	nonce, err := securityLayer.GenerateInternalNonce()

	assert.NoError(t, err)
	assert.Len(t, nonce, 8)

	savedNonce, err := securityLayer.internalNonceTable.Get(nonce[0])

	assert.NoError(t, err)
	assert.EqualValues(t, savedNonce, nonce)
}

func TestSecurityLayerGetExternalNonce(t *testing.T) {
	securityLayer := NewLayer(testNetworkKey)

	nonce, err := securityLayer.GetExternalNonce(1)
	assert.Error(t, err)
	assert.Nil(t, nonce)

	receivedNonce := []byte{0x98, 0xe4, 0x1b, 0x30, 0x84, 0x33, 0xf4, 0x3f}

	securityLayer.ReceiveNonce(1, security.NonceReport{
		NonceByte: receivedNonce,
	})

	nonce, err = securityLayer.GetExternalNonce(1)
	assert.NoError(t, err)
	assert.EqualValues(t, receivedNonce, nonce)
}

func TestSecurityLayerWaitForExternalNonce(t *testing.T) {
	securityLayer := NewLayer(testNetworkKey)

	done := make(chan bool)

	receivedNonce := []byte{0x98, 0xe4, 0x1b, 0x30, 0x84, 0x33, 0xf4, 0x3f}

	go func() {
		nonce, err := securityLayer.WaitForExternalNonce(1)
		assert.NoError(t, err)
		assert.EqualValues(t, receivedNonce, nonce)

		done <- true
	}()

	time.Sleep(time.Millisecond * 50)

	securityLayer.ReceiveNonce(1, security.NonceReport{
		NonceByte: receivedNonce,
	})

	<-done
}
