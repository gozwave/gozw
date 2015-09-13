package application

import "crypto/rand"

// GenerateNetworkKey generates a 16-byte encryption key using Go's crypto/rand
// package, which is cryptographically secure.
func GenerateNetworkKey() []byte {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}
