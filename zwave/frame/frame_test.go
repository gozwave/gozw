package frame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalFrame(t *testing.T) {
	t.Parallel()

	frame := NewRequestFrame([]byte{
		0x13,
		0x01,
	})

	marshalled, err := frame.MarshalBinary()

	assert.NoError(t, err)
	assert.Len(t, marshalled, 6)
	assert.EqualValues(t, []byte{0x01, 0x04, 0x00, 0x13, 0x01, 0xe9}, marshalled)

	marshalled, err = NewAckFrame().MarshalBinary()
	assert.NoError(t, err)
	assert.Len(t, marshalled, 1)
	assert.EqualValues(t, []byte{HeaderAck}, marshalled)

	marshalled, err = NewNakFrame().MarshalBinary()
	assert.NoError(t, err)
	assert.Len(t, marshalled, 1)
	assert.EqualValues(t, []byte{HeaderNak}, marshalled)

	marshalled, err = NewCanFrame().MarshalBinary()
	assert.NoError(t, err)
	assert.Len(t, marshalled, 1)
	assert.EqualValues(t, []byte{HeaderCan}, marshalled)

}

func TestChecksum(t *testing.T) {
	t.Parallel()

	frame := NewRequestFrame([]byte{
		0x13,
		0x01,
	})

	marshalled, err := frame.MarshalBinary()

	assert.NoError(t, err)
	assert.Len(t, marshalled, 6)
	assert.EqualValues(t, 0xe9, frame.CalcChecksum())

	frame = NewAckFrame()
	assert.NoError(t, frame.VerifyChecksum())
}

func TestUnmarshalFrame(t *testing.T) {
	t.Parallel()

	frame := NewRequestFrame([]byte{
		0x13,
		0x01,
	})

	marshalled, err := frame.MarshalBinary()
	assert.NoError(t, err)

	frame = UnmarshalFrame(marshalled)
	assert.EqualValues(t, []byte{0x13, 0x01}, frame.Payload)

	frame = UnmarshalFrame([]byte{0x06})
	assert.True(t, frame.IsAck())
}
