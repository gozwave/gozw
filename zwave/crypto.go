package zwave

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

var AuthPassword = []byte{
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
	0x55,
}

var EncryptPassword = []byte{
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
	0xAA,
}

var InclusionKey = []byte{
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
	0x00,
}

// @todo don't be terrible
var NetworkKey = []byte{
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

var InclusionEncKey = EncryptEBS(InclusionKey, EncryptPassword)
var InclusionAuthKey = EncryptEBS(InclusionKey, AuthPassword)

var NetworkEncKey = EncryptEBS(NetworkKey, EncryptPassword)
var NetworkAuthKey = EncryptEBS(NetworkKey, AuthPassword)

func GenerateNonce() []byte {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		// @todo
		panic(err)
	}

	return buf
}

func EncryptEBS(key []byte, message []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	out := make([]byte, aes.BlockSize)
	block.Encrypt(out, message)
	return out
}

func CryptMessage(input, iv, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	output := make([]byte, len(input))

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(output, input)

	return output
}

func CalculateHMAC(payload, iv, key []byte) []byte {
	input := padPayloadToBlockSize(payload)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	output := make([]byte, len(input))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(output, input)

	totalBlocks := (len(output) / 16)
	if len(output)%16 == 0 {
		totalBlocks -= 1
	}
	lastBlockOffset := totalBlocks * 16

	return output[lastBlockOffset : lastBlockOffset+8]
}

func padPayloadToBlockSize(message []byte) []byte {
	// pad the message with null bytes until it is the correct size
	for len(message)%aes.BlockSize != 0 {
		message = append(message, 0)
	}

	return message
}
