package thermostatsetpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalReport(t *testing.T) {
	report := Report{}

	err := report.UnmarshalBinary([]byte{
		0x00,
		0x00,
		0xFF,
		0x71, // size, scale, precision = 1
		0xAA,
	})

	assert.NoError(t, err)
	assert.EqualValues(t, 0x0F, report.Level.SetpointType)
	assert.EqualValues(t, 1, report.Level2.Size)
	assert.EqualValues(t, 2, report.Level2.Scale)
	assert.EqualValues(t, 3, report.Level2.Precision)
	assert.EqualValues(t, []byte{0xAA}, report.Value)
}

func TestUnmarshalReportHandlesBadPayloads(t *testing.T) {
	report := Report{}

	assert.Error(t, report.UnmarshalBinary([]byte{}))
	assert.Error(t, report.UnmarshalBinary([]byte{0x00}))
	assert.Error(t, report.UnmarshalBinary([]byte{0x00, 0x00}))
	assert.Error(t, report.UnmarshalBinary(nil))
}

func TestMarshalReport(t *testing.T) {
	report := Report{
		Level:  struct{ SetpointType byte }{SetpointType: 0x01},
		Level2: struct{ Size, Scale, Precision byte }{0x01, 0x02, 0x03},
		Value:  []byte{0xAA},
	}

	data, err := report.MarshalBinary()

	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x01, 0x71, 0xAA}, data[2:])
}
