package commandclass

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseThermostatSetpointReport(t *testing.T) {
	payload := []byte{
		CommandClassThermostatSetpoint,
		CommandThermostatSetpointReport,
		0x01,             // cooling
		0x29,             // Precision = 1, Scale = 1, Size = 1
		byte(int8(0x50)), // 80 in decimal
	}

	report := ParseThermostatSetpointReport(payload)

	assert.EqualValues(t, 0x01, report.Type)
	assert.EqualValues(t, 0x01, report.Precision)
	assert.EqualValues(t, 0x01, report.Scale)
	assert.EqualValues(t, 0x01, report.Size)
	assert.EqualValues(t, []byte{0x50}, report.Value)
}

func TestGetTemperatureOneByteValue(t *testing.T) {
	testcases := []struct {
		Expected  float64
		Precision byte
		Value     byte
	}{
		{-128, 0, 0x80},
		{-2, 0, 0xFE},
		{-1, 0, 0xFF},
		{0, 0, 0x00},
		{1, 0, 0x01},
		{2, 0, 0x02},
		{127, 0, 0x7F},

		{-12.8, 1, 0x80},
		{-0.2, 1, 0xFE},
		{-0.1, 1, 0xFF},
		{0, 1, 0x00},
		{0.1, 1, 0x01},
		{0.2, 1, 0x02},
		{12.7, 1, 0x7F},

		{-1.28, 2, 0x80},
		{-0.02, 2, 0xFE},
		{-0.01, 2, 0xFF},
		{0, 2, 0x00},
		{0.01, 2, 0x01},
		{0.02, 2, 0x02},
		{1.27, 2, 0x7F},
	}

	for _, testcase := range testcases {
		payload := []byte{
			CommandClassThermostatSetpoint,
			CommandThermostatSetpointReport,
			0x01, // cooling
			(testcase.Precision << 5) | 0x09, // Precision = X, Scale = 1, Size = 1
			testcase.Value,
		}

		report := ParseThermostatSetpointReport(payload)
		temp, err := report.GetTemperature()

		assert.NoError(t, err)

		assert.EqualValues(t, SetpointScaleFarenheit, temp.Scale)
		assert.EqualValues(t, testcase.Precision, report.Precision)
		assert.InDelta(t, testcase.Expected, temp.Value, 0.001)
	}
}

func TestGetTemperatureTwoByteValue(t *testing.T) {
	testcases := []struct {
		Expected  float64
		Precision byte
		Value     uint16
	}{
		{-32768, 0, 0x8000},
		{-2, 0, 0xFFFE},
		{-1, 0, 0xFFFF},
		{0, 0, 0x0000},
		{1, 0, 0x0001},
		{2, 0, 0x0002},
		{32767, 0, 0x7FFF},

		{-3276.8, 1, 0x8000},
		{-0.2, 1, 0xFFFE},
		{-0.1, 1, 0xFFFF},
		{0, 1, 0x0000},
		{0.1, 1, 0x0001},
		{0.2, 1, 0x0002},
		{3276.7, 1, 0x7FFF},

		{-327.68, 2, 0x8000},
		{-0.02, 2, 0xFFFE},
		{-0.01, 2, 0xFFFF},
		{0, 2, 0x0000},
		{0.01, 2, 0x0001},
		{0.02, 2, 0x0002},
		{327.67, 2, 0x7FFF},
	}

	for _, testcase := range testcases {
		payload := []byte{
			CommandClassThermostatSetpoint,
			CommandThermostatSetpointReport,
			0x01, // cooling
			(testcase.Precision << 5) | 0x0A, // Precision = X, Scale = 1, Size = 2
		}

		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, testcase.Value)
		payload = append(payload, buf...)

		report := ParseThermostatSetpointReport(payload)
		temp, err := report.GetTemperature()

		assert.NoError(t, err)

		assert.EqualValues(t, SetpointScaleFarenheit, temp.Scale)
		assert.EqualValues(t, testcase.Precision, report.Precision)
		assert.InDelta(t, testcase.Expected, temp.Value, 0.001)
	}
}

func TestGetTemperatureFourByteValue(t *testing.T) {
	testcases := []struct {
		Expected  float64
		Precision byte
		Value     uint32
	}{
		{-2147483648, 0, 0x80000000},
		{-2, 0, 0xFFFFFFFE},
		{-1, 0, 0xFFFFFFFF},
		{0, 0, 0x00000000},
		{1, 0, 0x00000001},
		{2, 0, 0x00000002},
		{2147483647, 0, 0x7FFFFFFF},

		{-214748364.8, 1, 0x80000000},
		{-0.2, 1, 0xFFFFFFFE},
		{-0.1, 1, 0xFFFFFFFF},
		{0, 1, 0x00000000},
		{0.1, 1, 0x00000001},
		{0.2, 1, 0x00000002},
		{214748364.7, 1, 0x7FFFFFFF},

		{-21474836.48, 2, 0x80000000},
		{-0.02, 2, 0xFFFFFFFE},
		{-0.01, 2, 0xFFFFFFFF},
		{0, 2, 0x00000000},
		{0.01, 2, 0x00000001},
		{0.02, 2, 0x00000002},
		{21474836.47, 2, 0x7FFFFFFF},
	}

	for _, testcase := range testcases {
		payload := []byte{
			CommandClassThermostatSetpoint,
			CommandThermostatSetpointReport,
			0x01, // cooling
			(testcase.Precision << 5) | 0x0C, // Precision = X, Scale = 1, Size = 4
		}

		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, testcase.Value)
		payload = append(payload, buf...)

		report := ParseThermostatSetpointReport(payload)
		temp, err := report.GetTemperature()

		assert.NoError(t, err)

		assert.EqualValues(t, SetpointScaleFarenheit, temp.Scale)
		assert.EqualValues(t, testcase.Precision, report.Precision)
		assert.InDelta(t, testcase.Expected, temp.Value, 0.001)
	}
}

func TestGetTemperatureInvalidSize(t *testing.T) {
	payload := []byte{
		CommandClassThermostatSetpoint,
		CommandThermostatSetpointReport,
		0x01, // cooling
		0x03, // Precision = 0, Scale = 0, Size = 3
		0x00, 0x00, 0x00,
	}

	report := ParseThermostatSetpointReport(payload)
	temp, err := report.GetTemperature()

	assert.Error(t, err)
	assert.Nil(t, temp)
}
