package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratingNetworkKey(t *testing.T) {
	key := GenerateNetworkKey()
	assert.Len(t, key, 16)
}
