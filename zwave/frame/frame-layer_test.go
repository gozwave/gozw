package frame

import (
	"testing"
	"time"

	"github.com/bjyoungblood/gozw/zwave/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGoodIncomingFrameResultsInAck(t *testing.T) {
	t.Parallel()

	bytes := make(chan byte, 100)
	var bytesFromTransport <-chan byte = bytes

	transport := new(mocks.TransportLayer)
	transport.On("Read").Return(bytesFromTransport).Once()
	transport.On("Write", []byte{FrameHeaderAck}).Return(1, nil)

	frameLayer := NewFrameLayer(transport)

	bytes <- 0x01
	bytes <- 0x04
	bytes <- 0x01
	bytes <- 0x13
	bytes <- 0x01
	bytes <- 0xe8

	frame := <-frameLayer.GetOutputChannel()

	// Ensure the other goroutines have time to do their thing
	time.Sleep(200 * time.Millisecond)

	// Ensure ack was written back to the transport
	transport.AssertCalled(t, "Write", []byte{FrameHeaderAck})
	transport.AssertExpectations(t)

	// Ensure the frame read from the transport is correct
	assert.True(t, frame.IsResponse())
	assert.True(t, frame.IsData())
	assert.EqualValues(t, 0x13, frame.Payload[0])
	assert.EqualValues(t, 0x01, frame.Payload[1])
	assert.NoError(t, frame.VerifyChecksum())
}

func TestBadIncomingFrameResultsInNak(t *testing.T) {
	t.Parallel()

	bytes := make(chan byte, 100)
	var bytesFromTransport <-chan byte = bytes

	transport := new(mocks.TransportLayer)
	transport.On("Read").Return(bytesFromTransport)
	transport.On("Write", []byte{FrameHeaderNak}).Return(1, nil)

	_ = NewFrameLayer(transport)

	bytes <- 0x01
	bytes <- 0x04
	bytes <- 0x01
	bytes <- 0x13
	bytes <- 0x01
	bytes <- 0x99

	// Ensure the other goroutines have time to do their thing
	time.Sleep(200 * time.Millisecond)

	// Ensure nak was written back to the transport
	transport.AssertCalled(t, "Write", []byte{FrameHeaderNak})
	transport.AssertExpectations(t)
}

func TestOutgoingFrameWrittenCorrectly(t *testing.T) {
	// @todo write me
}
