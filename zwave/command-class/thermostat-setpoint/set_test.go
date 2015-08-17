package thermostatsetpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSet(t *testing.T) {
	set := ThermostatSetpointSet{}

	err := set.UnmarshalBinary([]byte{
		0xFF,
		0x71, // size, scale, precision = 1
		0xAA,
	})

	assert.NoError(t, err)
	assert.EqualValues(t, 0x0F, set.Level.SetpointType)
	assert.EqualValues(t, 1, set.Level2.Size)
	assert.EqualValues(t, 2, set.Level2.Scale)
	assert.EqualValues(t, 3, set.Level2.Precision)
	assert.EqualValues(t, []byte{0xAA}, set.Value)
}

func TestUnmarshalSetHandlesBadPayloads(t *testing.T) {
	set := ThermostatSetpointSet{}

	assert.Error(t, set.UnmarshalBinary([]byte{}))
	assert.Error(t, set.UnmarshalBinary([]byte{0x00}))
	assert.Error(t, set.UnmarshalBinary([]byte{0x00, 0x00}))
	assert.Error(t, set.UnmarshalBinary(nil))
}

func TestMarshalSet(t *testing.T) {
	set := ThermostatSetpointSet{
		Level:  struct{ SetpointType byte }{SetpointType: 0x01},
		Level2: struct{ Size, Scale, Precision byte }{0x01, 0x02, 0x03},
		Value:  []byte{0xAA},
	}

	data, err := set.MarshalBinary()

	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x01, 0x71, 0xAA}, data)
}
