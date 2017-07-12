package frame

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/gozwave/gozw/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGoodIncomingFrameResultsInAck(t *testing.T) {
	t.Parallel()

	pr, pw := io.Pipe()
	buf := &testutil.TestBuffer{
		ReadableBytes: pr,
		BytesWritten:  bytes.NewBuffer([]byte{}),
	}

	frameLayer := NewFrameLayer(io.ReadWriter(buf))

	pw.Write([]byte{
		0x01,
		0x04,
		0x01,
		0x13,
		0x01,
		0xe8,
	})

	select {
	case frame := <-frameLayer.GetOutputChannel():

		// Ensure ack was written back to the transport
		assert.EqualValues(t, []byte{HeaderAck}, buf.BytesWritten.Bytes())

		// Ensure the frame read from the transport is correct
		assert.True(t, frame.IsResponse())
		assert.True(t, frame.IsData())
		assert.EqualValues(t, 0x13, frame.Payload[0])
		assert.EqualValues(t, 0x01, frame.Payload[1])
		assert.NoError(t, frame.VerifyChecksum())
	case <-time.After(time.Second):
		assert.Fail(t, "Timeout >1s")
	}
}

func TestBadIncomingFrameResultsInNak(t *testing.T) {
	t.Parallel()

	pr, pw := io.Pipe()
	buf := &testutil.TestBuffer{
		ReadableBytes: pr,
		BytesWritten:  bytes.NewBuffer([]byte{}),
	}

	_ = NewFrameLayer(io.ReadWriter(buf))

	pw.Write([]byte{
		0x01,
		0x04,
		0x01,
		0x13,
		0x01,
		0x99,
	})

	// Ensure the other goroutines have time to do their thing
	time.Sleep(200 * time.Millisecond)

	// Ensure nak was written back to the transport
	assert.EqualValues(t, []byte{HeaderNak}, buf.BytesWritten.Bytes())
}

func TestOutgoingFrameWrittenCorrectly(t *testing.T) {
	// @todo write me
}
