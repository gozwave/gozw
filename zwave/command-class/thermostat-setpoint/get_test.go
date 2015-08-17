package thermostatsetpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalGet(t *testing.T) {
	get := ThermostatSetpointGet{}

	err := get.UnmarshalBinary([]byte{0xFF})

	assert.NoError(t, err)
	assert.EqualValues(t, 0x0F, get.Level.SetpointType)

	err = get.UnmarshalBinary([]byte{0x01})

	assert.NoError(t, err)
	assert.EqualValues(t, 0x01, get.Level.SetpointType)
}

func TestUnmarshalGetHandlesBadPayloads(t *testing.T) {
	get := ThermostatSetpointGet{}

	assert.Error(t, get.UnmarshalBinary([]byte{}))
	assert.Error(t, get.UnmarshalBinary(nil))
}

func TestMarshalGet(t *testing.T) {
	get := ThermostatSetpointGet{
		Level: struct{ SetpointType byte }{SetpointType: 0xFF},
	}

	data, err := get.MarshalBinary()

	assert.NoError(t, err)
	assert.EqualValues(t, 0x0F, data[0])
}
