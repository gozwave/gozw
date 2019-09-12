package frame

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParsingDataFrame(t *testing.T) {
	t.Parallel()

	parserInput := make(chan byte)
	parserOutput := make(chan *ParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	NewParser(ctx, parserInput, parserOutput, acks, naks, cans, logger)

	parserInput <- 0x01
	parserInput <- 0x04
	parserInput <- 0x01
	parserInput <- 0x13
	parserInput <- 0x01
	parserInput <- 0xe8

	parserEvent := <-parserOutput
	frame := parserEvent.frame

	assert.Equal(t, ParseOk, parserEvent.status)
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
	t.Parallel()

	parserInput := make(chan byte)
	parserOutput := make(chan *ParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	NewParser(ctx, parserInput, parserOutput, acks, naks, cans, logger)

	parserInput <- 0x01
	parserInput <- 0x04
	parserInput <- 0x01
	parserInput <- 0x13
	parserInput <- 0x01
	parserInput <- 0xe0

	parserEvent := <-parserOutput
	frame := parserEvent.frame

	assert.Equal(t, ParseNotOk, parserEvent.status)
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
	t.Parallel()

	parserInput := make(chan byte)
	parserOutput := make(chan *ParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	NewParser(ctx, parserInput, parserOutput, acks, naks, cans, logger)

	parserInput <- 0x01
	parserInput <- 0x04
	parserInput <- 0x01
	parserInput <- 0x13

	time.Sleep(readTimeout + 10*time.Millisecond)

	parserEvent := <-parserOutput

	assert.Equal(t, ParseTimeout, parserEvent.status)
}

func TestAcksNaksCans(t *testing.T) {
	t.Parallel()

	parserInput := make(chan byte)
	parserOutput := make(chan *ParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	NewParser(ctx, parserInput, parserOutput, acks, naks, cans, logger)

	var event bool

	parserInput <- HeaderAck
	event = <-acks

	assert.True(t, event)

	parserInput <- HeaderNak
	event = <-naks

	assert.True(t, event)

	parserInput <- HeaderCan
	event = <-cans

	assert.True(t, event)
}

func TestRecoversAfterInvalidLength(t *testing.T) {
	t.Parallel()

	parserInput := make(chan byte)
	parserOutput := make(chan *ParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	NewParser(ctx, parserInput, parserOutput, acks, naks, cans, logger)

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

	assert.Equal(t, ParseNotOk, parserEvent.status)
	assert.True(t, frame.IsData())
	assert.True(t, frame.IsResponse())
}

func BenchmarkParsingShortDataFrame(b *testing.B) {
	parserInput := make(chan byte)
	parserOutput := make(chan *ParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	NewParser(ctx, parserInput, parserOutput, acks, naks, cans, logger)

	for n := 0; n < b.N; n++ {

		parserInput <- 0x01
		parserInput <- 0x04
		parserInput <- 0x01
		parserInput <- 0x13
		parserInput <- 0x01
		parserInput <- 0xe8

		parserEvent := <-parserOutput
		assert.Equal(b, ParseOk, parserEvent.status)
	}
}

func BenchmarkParsingAckFrame(b *testing.B) {
	parserInput := make(chan byte)
	parserOutput := make(chan *ParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	NewParser(ctx, parserInput, parserOutput, acks, naks, cans, logger)

	for n := 0; n < b.N; n++ {

		parserInput <- 0x06

		<-acks
	}
}
