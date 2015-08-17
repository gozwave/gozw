package thermostatsetpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSupportedReport(t *testing.T) {
	report := SupportedReport{}

	err := report.UnmarshalBinary([]byte{
		0xFF,
		0x71,
		0xAA,
	})

	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0xFF, 0x71, 0xAA}, report.BitMask)
}

func TestUnmarshalSupportedReportHandlesBadPayloads(t *testing.T) {
	report := SupportedReport{}

	assert.Error(t, report.UnmarshalBinary([]byte{}))
	assert.Error(t, report.UnmarshalBinary(nil))
}

func TestMarshalSupportedReport(t *testing.T) {
	report := SupportedReport{
		BitMask: []byte{0x01, 0x71, 0xAA},
	}

	data, err := report.MarshalBinary()

	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x01, 0x71, 0xAA}, data)
}
