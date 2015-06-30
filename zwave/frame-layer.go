package zwave

import (
	"fmt"

	"github.com/looplab/fsm"
)

type FrameLayer struct {
	transport   *TransportLayer
	frameParser *FrameParser

	parserInput  chan<- byte
	parserOutput <-chan *FrameParseEvent

	pendingWrites chan *Frame
	frameOutput   chan *Frame

	state *fsm.FSM
}

func NewFrameLayer(transport *TransportLayer, debug bool) *FrameLayer {
	parserInput := make(chan byte)
	parserOutput := make(chan *FrameParseEvent)

	frameLayer := &FrameLayer{
		transport:   transport,
		frameParser: NewFrameParser(parserInput, parserOutput),

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
			"RX_COMPLETE": func(e *fsm.Event) {
				event := e.Args[0].(*FrameParseEvent)

				if event.status == FrameParseOk {
					frameLayer.sendAck()
					frameLayer.frameOutput <- event.frame
				} else if event.status == FrameParseNotOk {
					fmt.Println("sent a nak for some reason")
					frameLayer.sendNak()
				} else {
					fmt.Println("frame parse timeout or something")
				}
			},
		},
	)

	go frameLayer.handle()

	return frameLayer
}

func (f *FrameLayer) Write(frame *Frame) {
	f.pendingWrites <- frame
}

func (f *FrameLayer) Read() <-chan *Frame {
	return f.frameOutput
}

func (layer *FrameLayer) handleReceive(transportBytesIn <-chan byte, firstByte byte) {
	if firstByte == FrameSOFData {
		layer.state.Event("RX_SOF")
	} else {
		if firstByte == FrameSOFAck {
			layer.state.Event("RX_ACK")
		} else if firstByte == FrameSOFNak {
			layer.state.Event("RX_NAK")
		} else if firstByte == FrameSOFCan {
			layer.state.Event("RX_CAN")
		}

		return
	}

	layer.parserInput <- firstByte

	for {
		select {
		case event := <-layer.parserOutput:
			layer.state.Event("RX_COMPLETE", event)
			break

		case nextByte := <-transportBytesIn:
			layer.parserInput <- nextByte
		}
	}
}

func (layer *FrameLayer) handle() {

	transportBytesIn := layer.transport.Read()

	for {
		select {
		case firstByte := <-transportBytesIn:
			layer.handleReceive(transportBytesIn, firstByte)

		case writeFrame := <-layer.pendingWrites:
			layer.state.Event("TX_DATA")
			layer.writeToTransport(writeFrame.Marshal())
		}
	}
}

func (f *FrameLayer) writeToTransport(buf []byte) (int, error) {
	return f.transport.Write(buf)
}

func (f *FrameLayer) sendAck() error {
	_, err := f.transport.Write([]byte{FrameSOFAck})
	return err
}

func (f *FrameLayer) sendNak() error {
	_, err := f.transport.Write([]byte{FrameSOFNak})
	return err
}
