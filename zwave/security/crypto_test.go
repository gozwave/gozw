package security

import (
	"testing"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/stretchr/testify/assert"
)

var testMessagePlaintext = []byte{0x00, 0x98, 0x06, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
var testMessageCiphertext = []byte{0x4A, 0xD6, 0xB1, 0x33, 0xB8, 0xFA, 0x0F, 0x2E, 0x0A, 0xEB, 0x86, 0x87, 0x7B, 0xB2, 0xDF, 0x11, 0x13, 0x4E, 0xB4}
var testEncryptionKey = EncryptEBS(InclusionKey, EncryptPassword)
var testAuthKey = EncryptEBS(InclusionKey, AuthPassword)

var testIV = []byte{
	0xAA, 0xAA, 0xAA, 0xAA,
	0xAA, 0xAA, 0xAA, 0xAA,
	0x4F, 0x46, 0x61, 0x76,
	0x01, 0x48, 0x7B, 0xFF,
}

var authIV = []byte{
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
}

func TestGenerateNonceReturns8ByteSlice(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 8, len(GenerateNonce()))
}

func TestEncryptEBS(t *testing.T) {
	t.Parallel()

	encryptedKey := EncryptEBS(InclusionKey, EncryptPassword)
	assert.NotEqual(t, InclusionKey, encryptedKey)
	assert.Len(t, encryptedKey, 16)
}

func TestCryptMessage(t *testing.T) {
	t.Parallel()

	encryptedMessage := CryptMessage(testMessagePlaintext, testIV, testEncryptionKey)

	assert.Equal(t, testMessageCiphertext, encryptedMessage)
	assert.NotEqual(t, testMessagePlaintext, encryptedMessage)

	decryptedMessage := CryptMessage(testMessageCiphertext, testIV, testEncryptionKey)

	assert.Equal(t, testMessagePlaintext, decryptedMessage)
}

func TestCalculateHMAC(t *testing.T) {
	t.Parallel()

	authData := []byte{
		commandclass.CommandSecurityMessageEncapsulation,
		1, // sender node id
		3, // receiver node id
		uint8(len(testMessageCiphertext)),
	}

	authData = append(testIV, authData...)

	authData = append(authData, testMessageCiphertext...)

	hmac := CalculateHMAC(authData, authIV, testAuthKey)

	expectedHmac := []byte{
		0xF0, 0xDE, 0x2E, 0xB2,
		0x51, 0xEA, 0x0F, 0xF7,
	}

	assert.Equal(t, expectedHmac, hmac)
}

func TestPadPayloadToBlockSize(t *testing.T) {
	t.Parallel()

	short := []byte{0, 1, 2, 3, 4}
	message := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	assert.Len(t, padPayloadToBlockSize(message), 16)
	assert.Len(t, padPayloadToBlockSize(short), 16)
}
