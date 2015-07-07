package zwave

import "fmt"

type FrameLayer interface {
	Write(frame *Frame)
	GetOutputChannel() <-chan Frame
}

type SerialFrameLayer struct {
	transportLayer TransportLayer

	frameParser      *FrameParser
	parserInput      chan<- byte
	parserOutput     <-chan *FrameParseEvent
	acks, naks, cans <-chan bool

	pendingWrites chan *Frame
	frameOutput   chan Frame
}

func NewFrameLayer(transportLayer TransportLayer) *SerialFrameLayer {
	parserInput := make(chan byte)
	parserOutput := make(chan *FrameParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	frameLayer := &SerialFrameLayer{
		transportLayer: transportLayer,

		frameParser:  NewFrameParser(parserInput, parserOutput, acks, naks, cans),
		parserInput:  parserInput,
		parserOutput: parserOutput,
		acks:         acks,
		naks:         naks,
		cans:         cans,

		pendingWrites: make(chan *Frame),
		frameOutput:   make(chan Frame),
	}

	go frameLayer.bgWork()
	go frameLayer.bgRead()

	return frameLayer
}

func (layer *SerialFrameLayer) bgWork() {

	for {
		select {
		case frameIn := <-layer.parserOutput:
			if frameIn.status == FrameParseOk {
				layer.sendAck()
				layer.frameOutput <- frameIn.frame
			} else if frameIn.status == FrameParseNotOk {
				layer.sendNak()
			} else {
				// @todo handle timeout(?)
			}

		case <-layer.acks:
			fmt.Println("frame layer: rx ack")
		case <-layer.naks:
			fmt.Println("frame layer: rx nak")
		case <-layer.cans:
			fmt.Println("frame layer: rx can")

		case frameToWrite := <-layer.pendingWrites:
			layer.writeToTransport(frameToWrite.Marshal())
			_ = <-layer.acks

		}
	}
}

func (f *SerialFrameLayer) Write(frame *Frame) {
	go func() {
		f.pendingWrites <- frame
	}()
}

func (f *SerialFrameLayer) GetOutputChannel() <-chan Frame {
	return f.frameOutput
}

func (f *SerialFrameLayer) bgRead() {
	for eachByte := range f.transportLayer.Read() {
		f.parserInput <- eachByte
	}
}

func (f *SerialFrameLayer) writeToTransport(buf []byte) (int, error) {
	return f.transportLayer.Write(buf)
}

func (f *SerialFrameLayer) sendAck() error {
	_, err := f.transportLayer.Write([]byte{FrameHeaderAck})
	return err
}

func (f *SerialFrameLayer) sendNak() error {
	_, err := f.transportLayer.Write([]byte{FrameHeaderNak})
	return err
}
