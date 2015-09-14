package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageEncapsulationUnmarshal(t *testing.T) {
	message := MessageEncapsulation{}

	err := message.UnmarshalBinary([]byte{
		0x98, 0x81, // command class and command ids
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, // 8-byte IV
		0xFF,                                           // 1-byte properties
		0xFE,                                           // 1-byte CCID
		0xFD,                                           // 1-byte CID
		0xFC,                                           // 1-byte command payload
		0xAA,                                           // 1-byte receiver nonce identifier
		0xBB, 0xCC, 0xDD, 0xFF, 0x88, 0x99, 0x77, 0x44, // 8-byte HMAC
	})

	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, message.InitializationVectorByte)
	assert.EqualValues(t, 0xFE, message.CommandClassIdentifier)
	assert.EqualValues(t, 0xFD, message.CommandIdentifier)
	assert.EqualValues(t, []byte{0xFC}, message.CommandByte)
	assert.EqualValues(t, 0xAA, message.ReceiversNonceIdentifier)
	assert.EqualValues(t, []byte{0xBB, 0xCC, 0xDD, 0xFF, 0x88, 0x99, 0x77, 0x44}, message.MessageAuthenticationCodeByte)
}
