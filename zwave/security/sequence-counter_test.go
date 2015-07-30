package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSequenceCounter(t *testing.T) {
	var seq uint8

	counter := NewSequenceCounter()

	for i := 1; i <= 15; i++ {
		seq = counter.Get(1)
		assert.EqualValues(t, i, seq)
	}

	seq = counter.Get(1)
	assert.EqualValues(t, 1, seq)

	seq = counter.Get(1)
	assert.EqualValues(t, 2, seq)

	seq = counter.Get(2)
	assert.EqualValues(t, 1, seq)

	seq = counter.Get(2)
	assert.EqualValues(t, 2, seq)
}
