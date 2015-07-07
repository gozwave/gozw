package zwave

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParsingDataFrame(t *testing.T) {

	parserInput := make(chan byte)
	parserOutput := make(chan *FrameParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	NewFrameParser(parserInput, parserOutput, acks, naks, cans)

	parserInput <- 0x01
	parserInput <- 0x04
	parserInput <- 0x01
	parserInput <- 0x13
	parserInput <- 0x01
	parserInput <- 0xe8

	parserEvent := <-parserOutput
	frame := parserEvent.frame

	assert.Equal(t, FrameParseOk, parserEvent.status)
	assert.True(t, frame.IsData())
	assert.True(t, frame.IsResponse())
	assert.False(t, frame.IsAck())
	assert.False(t, frame.IsCan())
	assert.False(t, frame.IsRequest())
	assert.NoError(t, frame.VerifyChecksum())

	assert.EqualValues(t, 0x01, frame.Header)
	assert.EqualValues(t, 0x04, frame.Length)
	assert.EqualValues(t, 0x13, frame.Payload[0])
	assert.EqualValues(t, 0xe8, frame.CalcChecksum())
	assert.Len(t, frame.Payload, 2)

}

func TestInvalidChecksum(t *testing.T) {

	parserInput := make(chan byte)
	parserOutput := make(chan *FrameParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	NewFrameParser(parserInput, parserOutput, acks, naks, cans)

	parserInput <- 0x01
	parserInput <- 0x04
	parserInput <- 0x01
	parserInput <- 0x13
	parserInput <- 0x01
	parserInput <- 0xe0

	parserEvent := <-parserOutput
	frame := parserEvent.frame

	assert.Equal(t, FrameParseNotOk, parserEvent.status)
	assert.True(t, frame.IsData())
	assert.True(t, frame.IsResponse())
	assert.False(t, frame.IsAck())
	assert.False(t, frame.IsNak())
	assert.False(t, frame.IsCan())
	assert.False(t, frame.IsRequest())
	assert.Error(t, frame.VerifyChecksum())

	assert.EqualValues(t, 0x01, frame.Header)
	assert.EqualValues(t, 0x04, frame.Length)
	assert.EqualValues(t, 0x13, frame.Payload[0])
	assert.EqualValues(t, 0xe8, frame.CalcChecksum())
	assert.Len(t, frame.Payload, 2)

}

func TestParseTimeout(t *testing.T) {
	parserInput := make(chan byte)
	parserOutput := make(chan *FrameParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	NewFrameParser(parserInput, parserOutput, acks, naks, cans)

	parserInput <- 0x01
	parserInput <- 0x04
	parserInput <- 0x01
	parserInput <- 0x13

	time.Sleep(readTimeout + 10*time.Millisecond)

	parserEvent := <-parserOutput

	assert.Equal(t, FrameParseTimeout, parserEvent.status)
}

func TestAcksNaksCans(t *testing.T) {
	parserInput := make(chan byte)
	parserOutput := make(chan *FrameParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	NewFrameParser(parserInput, parserOutput, acks, naks, cans)

	var event bool

	parserInput <- FrameHeaderAck
	event = <-acks

	assert.True(t, event)

	parserInput <- FrameHeaderNak
	event = <-naks

	assert.True(t, event)

	parserInput <- FrameHeaderCan
	event = <-cans

	assert.True(t, event)
}

func TestRecoversAfterInvalidLength(t *testing.T) {
	parserInput := make(chan byte)
	parserOutput := make(chan *FrameParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	NewFrameParser(parserInput, parserOutput, acks, naks, cans)

	parserInput <- 0x01
	parserInput <- 0xFF
	parserInput <- 0x01
	parserInput <- 0x04
	parserInput <- 0x01
	parserInput <- 0x13
	parserInput <- 0x01
	parserInput <- 0xe0

	parserEvent := <-parserOutput
	frame := parserEvent.frame

	assert.Equal(t, FrameParseNotOk, parserEvent.status)
	assert.True(t, frame.IsData())
	assert.True(t, frame.IsResponse())

}
