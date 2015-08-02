package security

import (
	"crypto/aes"
	"crypto/cipher"
)

var authIV = []byte{
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
}

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

var InclusionEncKey = EncryptEBS(InclusionKey, EncryptPassword)
var InclusionAuthKey = EncryptEBS(InclusionKey, AuthPassword)

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

func CalculateHMAC(payload, key []byte) []byte {
	input := padPayloadToBlockSize(payload)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	output := make([]byte, len(input))

	mode := cipher.NewCBCEncrypter(block, authIV)
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
