package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseZWFloat(t *testing.T) {
	testcases := []struct {
		Expected  float64
		Precision byte
		Size      byte
		Value     []byte
	}{
		{-128, 0, 1, []byte{0x80}},
		{-2, 0, 1, []byte{0xFE}},
		{-1, 0, 1, []byte{0xFF}},
		{0, 0, 1, []byte{0x00}},
		{1, 0, 1, []byte{0x01}},
		{2, 0, 1, []byte{0x02}},
		{127, 0, 1, []byte{0x7F}},

		{-12.8, 1, 1, []byte{0x80}},
		{-0.2, 1, 1, []byte{0xFE}},
		{-0.1, 1, 1, []byte{0xFF}},
		{0, 1, 1, []byte{0x00}},
		{0.1, 1, 1, []byte{0x01}},
		{0.2, 1, 1, []byte{0x02}},
		{12.7, 1, 1, []byte{0x7F}},

		{-1.28, 2, 1, []byte{0x80}},
		{-0.02, 2, 1, []byte{0xFE}},
		{-0.01, 2, 1, []byte{0xFF}},
		{0, 2, 1, []byte{0x00}},
		{0.01, 2, 1, []byte{0x01}},
		{0.02, 2, 1, []byte{0x02}},
		{1.27, 2, 1, []byte{0x7F}},

		{-32768, 0, 2, []byte{0x80, 0x00}},
		{-2, 0, 2, []byte{0xFF, 0xFE}},
		{-1, 0, 2, []byte{0xFF, 0xFF}},
		{0, 0, 2, []byte{0x00, 0x00}},
		{1, 0, 2, []byte{0x00, 0x01}},
		{2, 0, 2, []byte{0x00, 0x02}},
		{32767, 0, 2, []byte{0x7F, 0xFF}},

		{-3276.8, 1, 2, []byte{0x80, 0x00}},
		{-0.2, 1, 2, []byte{0xFF, 0xFE}},
		{-0.1, 1, 2, []byte{0xFF, 0xFF}},
		{0, 1, 2, []byte{0x00, 0x00}},
		{0.1, 1, 2, []byte{0x00, 0x01}},
		{0.2, 1, 2, []byte{0x00, 0x02}},
		{3276.7, 1, 2, []byte{0x7F, 0xFF}},

		{-327.68, 2, 2, []byte{0x80, 0x00}},
		{-0.02, 2, 2, []byte{0xFF, 0xFE}},
		{-0.01, 2, 2, []byte{0xFF, 0xFF}},
		{0, 2, 2, []byte{0x00, 0x00}},
		{0.01, 2, 2, []byte{0x00, 0x01}},
		{0.02, 2, 2, []byte{0x00, 0x02}},
		{327.67, 2, 2, []byte{0x7F, 0xFF}},

		{-2147483648, 0, 4, []byte{0x80, 0x00, 0x00, 0x00}},
		{-2, 0, 4, []byte{0xFF, 0xFF, 0xFF, 0xFE}},
		{-1, 0, 4, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
		{0, 0, 4, []byte{0x00, 0x00, 0x00, 0x00}},
		{1, 0, 4, []byte{0x00, 0x00, 0x00, 0x01}},
		{2, 0, 4, []byte{0x00, 0x00, 0x00, 0x02}},
		{2147483647, 0, 4, []byte{0x7F, 0xFF, 0xFF, 0xFF}},

		{-214748364.8, 1, 4, []byte{0x80, 0x00, 0x00, 0x00}},
		{-0.2, 1, 4, []byte{0xFF, 0xFF, 0xFF, 0xFE}},
		{-0.1, 1, 4, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
		{0, 1, 4, []byte{0x00, 0x00, 0x00, 0x00}},
		{0.1, 1, 4, []byte{0x00, 0x00, 0x00, 0x01}},
		{0.2, 1, 4, []byte{0x00, 0x00, 0x00, 0x02}},
		{214748364.7, 1, 4, []byte{0x7F, 0xFF, 0xFF, 0xFF}},

		{-21474836.48, 2, 4, []byte{0x80, 0x00, 0x00, 0x00}},
		{-0.02, 2, 4, []byte{0xFF, 0xFF, 0xFF, 0xFE}},
		{-0.01, 2, 4, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
		{0, 2, 4, []byte{0x00, 0x00, 0x00, 0x00}},
		{0.01, 2, 4, []byte{0x00, 0x00, 0x00, 0x01}},
		{0.02, 2, 4, []byte{0x00, 0x00, 0x00, 0x02}},
		{21474836.47, 2, 4, []byte{0x7F, 0xFF, 0xFF, 0xFF}},
	}

	for _, testcase := range testcases {
		val, err := ParseZWFloat(testcase.Size, 1, testcase.Precision, testcase.Value)

		assert.NoError(t, err)

		assert.EqualValues(t, 1, val.Scale)
		assert.InDelta(t, testcase.Expected, val.Value, 0.001)
	}
}

func TestGetTemperatureInvalidSize(t *testing.T) {
	val, err := ParseZWFloat(3, 1, 1, []byte{0})

	assert.Error(t, err)
	assert.Nil(t, val)
}
