package zwave

import (
	"fmt"

	"github.com/looplab/fsm"
)

type FrameLayer struct {
	transportLayer *TransportLayer
	frameParser    *FrameParser

	parserInput  chan<- byte
	parserOutput <-chan *FrameParseEvent

	pendingWrites chan *Frame
	frameOutput   chan *Frame

	state *fsm.FSM
}

func NewFrameLayer(transportLayer *TransportLayer) *FrameLayer {
	parserInput := make(chan byte)
	parserOutput := make(chan *FrameParseEvent)

	frameLayer := &FrameLayer{
		transportLayer: transportLayer,
		frameParser:    NewFrameParser(parserInput, parserOutput),

		parserInput:  parserInput,
		parserOutput: parserOutput,

		pendingWrites: make(chan *Frame),
		frameOutput:   make(chan *Frame),
	}

	frameLayer.state = fsm.NewFSM(
		"idle",
		fsm.Events{
			{Name: "TX_DATA", Src: []string{"idle"}, Dst: "awaiting_ack"},
			{Name: "RX_ACK", Src: []string{"awaiting_ack"}, Dst: "idle"},
			{Name: "RX_NAK", Src: []string{"awaiting_ack"}, Dst: "idle"},
			{Name: "RX_CAN", Src: []string{"idle", "awaiting_ack"}, Dst: "idle"},
			{Name: "RX_SOF", Src: []string{"idle"}, Dst: "parse_frame"},
			{Name: "RX_COMPLETE", Src: []string{"parse_frame"}, Dst: "idle"},
		},
		fsm.Callbacks{
			"before_event": func(e *fsm.Event) {
				fmt.Printf("%s: %s -> %s\n", e.Event, e.Src, e.Dst)
			},
		},
	)

	go frameLayer.loop()

	return frameLayer
}

func (layer *FrameLayer) loop() {
	transportBytesIn := layer.transportLayer.Read()

start:
	for {
		select {

		// Read
		case firstByte := <-transportBytesIn:
			layer.parserInput <- firstByte

			for {
				select {
				case event := <-layer.parserOutput:

					if event.status == FrameParseOk {
						layer.sendAck()
						layer.frameOutput <- event.frame
					} else if event.status == FrameParseNotOk {
						layer.sendNak()
						layer.frameOutput <- event.frame
					} else {
						fmt.Println("frame parse timeout or something")
					}

					goto start

				case nextByte := <-transportBytesIn:
					layer.parserInput <- nextByte
				}
			}

		// Write
		case writeFrame := <-layer.pendingWrites:
			layer.writeToTransport(writeFrame.Marshal())

		}
	}
}

func (f *FrameLayer) Write(frame *Frame) {
	go func() {
		f.pendingWrites <- frame
	}()
}

func (f *FrameLayer) GetOutput() <-chan *Frame {
	return f.frameOutput
}

func (f *FrameLayer) writeToTransport(buf []byte) (int, error) {
	return f.transportLayer.Write(buf)
}

func (f *FrameLayer) sendAck() error {
	_, err := f.transportLayer.Write([]byte{FrameSOFAck})
	return err
}

func (f *FrameLayer) sendNak() error {
	_, err := f.transportLayer.Write([]byte{FrameSOFNak})
	return err
}
