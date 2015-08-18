package thermostatsetpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalGet(t *testing.T) {
	get := Get{}

	err := get.UnmarshalBinary([]byte{0x00, 0x00, 0xFF})

	assert.NoError(t, err)
	assert.EqualValues(t, 0x0F, get.Level.SetpointType)

	err = get.UnmarshalBinary([]byte{0x00, 0x00, 0x01})

	assert.NoError(t, err)
	assert.EqualValues(t, 0x01, get.Level.SetpointType)
}

func TestUnmarshalGetHandlesBadPayloads(t *testing.T) {
	get := Get{}

	assert.Error(t, get.UnmarshalBinary([]byte{}))
	assert.Error(t, get.UnmarshalBinary(nil))
}

func TestMarshalGet(t *testing.T) {
	get := Get{
		Level: struct{ SetpointType byte }{SetpointType: 0xFF},
	}

	data, err := get.MarshalBinary()

	assert.NoError(t, err)
	assert.EqualValues(t, 0x0F, data[2])
}
