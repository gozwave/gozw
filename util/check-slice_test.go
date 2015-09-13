package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSliceLength(t *testing.T) {
	testSlice := []byte{1, 2, 3, 4, 5}

	assert.NoError(t, CheckSliceLength(testSlice, 0))
	assert.NoError(t, CheckSliceLength(testSlice, 1))
	assert.NoError(t, CheckSliceLength(testSlice, 2))
	assert.NoError(t, CheckSliceLength(testSlice, 3))
	assert.NoError(t, CheckSliceLength(testSlice, 4))

	assert.Error(t, CheckSliceLength(testSlice, 5))

	var nilSlice []byte
	assert.Error(t, CheckSliceLength(nilSlice, 0))

	emptySlice := []byte{}
	assert.Error(t, CheckSliceLength(emptySlice, 0))
}
