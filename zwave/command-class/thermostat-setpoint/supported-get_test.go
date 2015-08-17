package thermostatsetpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSupportedGet(t *testing.T) {
	supportedGet := SupportedGet{}

	assert.NoError(t, supportedGet.UnmarshalBinary(nil))
	assert.NoError(t, supportedGet.UnmarshalBinary([]byte{}))
}

func TestMarshalSupportedGet(t *testing.T) {
	set := SupportedGet{}

	data, err := set.MarshalBinary()

	assert.NoError(t, err)
	assert.EqualValues(t, []byte{}, data)
}
