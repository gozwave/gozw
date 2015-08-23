package frame

import (
	"log"
	"os"

	"github.com/comail/colog"
	"github.com/helioslabs/gozw/zwave/transport"
)

type ILayer interface {
	Write(frame *Frame)
	GetOutputChannel() <-chan Frame
}

type Layer struct {
	transportLayer transport.Transport

	frameParser      *Parser
	parserInput      chan<- byte
	parserOutput     <-chan *ParseEvent
	acks, naks, cans <-chan bool

	logger *log.Logger

	pendingWrites chan *Frame
	frameOutput   chan Frame
}

func NewFrameLayer(transportLayer transport.Transport) *Layer {
	frameLogger := colog.NewCoLog(os.Stdout, "frame ", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	frameLogger.ParseFields(true)

	parserInput := make(chan byte)
	parserOutput := make(chan *ParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	frameLayer := &Layer{
		transportLayer: transportLayer,

		frameParser:  NewParser(parserInput, parserOutput, acks, naks, cans),
		parserInput:  parserInput,
		parserOutput: parserOutput,
		acks:         acks,
		naks:         naks,
		cans:         cans,

		logger: frameLogger.NewLogger(),

		pendingWrites: make(chan *Frame),
		frameOutput:   make(chan Frame, 5),
	}

	go frameLayer.bgWork()
	go frameLayer.bgRead()

	return frameLayer
}

func (l *Layer) bgWork() {

	for {
		select {
		case frameIn := <-l.parserOutput:
			if frameIn.status == ParseOk {
				l.sendAck()
				l.frameOutput <- frameIn.frame
			} else if frameIn.status == ParseNotOk {
				l.sendNak()
			} else {
				// @todo handle timeout(?)
			}

		case <-l.acks:
			l.logger.Print("warn: rx ack")
		case <-l.naks:
			l.logger.Print("warn: rx nak")
		case <-l.cans:
			l.logger.Print("warn: rx can")

		case frameToWrite := <-l.pendingWrites:
			l.writeToTransport(frameToWrite.Marshal())
			_ = <-l.acks

		}
	}
}

func (l *Layer) Write(frame *Frame) {
	go func() {
		l.pendingWrites <- frame
	}()
}

func (l *Layer) GetOutputChannel() <-chan Frame {
	return l.frameOutput
}

func (l *Layer) bgRead() {
	for eachByte := range l.transportLayer.Read() {
		l.parserInput <- eachByte
	}
}

func (l *Layer) writeToTransport(buf []byte) (int, error) {
	return l.transportLayer.Write(buf)
}

func (l *Layer) sendAck() error {
	_, err := l.transportLayer.Write([]byte{HeaderAck})
	return err
}

func (l *Layer) sendNak() error {
	_, err := l.transportLayer.Write([]byte{HeaderNak})
	return err
}
