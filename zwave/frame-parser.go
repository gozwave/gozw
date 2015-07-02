package zwave

import (
	"bytes"
	"time"

	"github.com/looplab/fsm"
)

type FrameParseStatus int

const (
	FrameParseOk FrameParseStatus = iota
	FrameParseNotOk
	FrameParseTimeout
)

const (
	minFrameSize uint8 = 3
	maxFrameSize uint8 = 88
)

const readTimeout time.Duration = 1500 * time.Millisecond

type FrameParseEvent struct {
	status FrameParseStatus
	frame  Frame
}

type FrameParser struct {
	state                              *fsm.FSM
	input                              <-chan byte
	framesReceived                     chan<- *FrameParseEvent
	acks, naks, cans                   chan<- bool
	sof, length, checksum, readCounter uint8
	payloadReadBuffer                  *bytes.Buffer
	parseTimeout                       *time.Timer
}

func NewFrameParser(input <-chan byte, output chan<- *FrameParseEvent, acks, naks, cans chan bool) *FrameParser {
	frameParser := &FrameParser{
		input:             input,
		framesReceived:    output,
		acks:              acks,
		naks:              naks,
		cans:              cans,
		payloadReadBuffer: bytes.NewBuffer([]byte{}),
		parseTimeout:      time.NewTimer(readTimeout),
	}

	frameParser.parseTimeout.Stop()

	frameParser.state = fsm.NewFSM(
		"idle",
		fsm.Events{
			{Name: "PARSE_TIMEOUT", Src: []string{"idle", "length", "data", "checksum"}, Dst: "idle"},
			{Name: "RX_SOF", Src: []string{"idle"}, Dst: "length"},
			{Name: "RX_ACK", Src: []string{"idle"}, Dst: "idle"},
			{Name: "RX_CAN", Src: []string{"idle"}, Dst: "idle"},
			{Name: "RX_NAK", Src: []string{"idle"}, Dst: "idle"},
			{Name: "RX_LENGTH", Src: []string{"length"}, Dst: "data"},
			{Name: "INVALID_LENGTH", Src: []string{"length"}, Dst: "idle"},
			{Name: "RX_DATA", Src: []string{"length", "data"}, Dst: "data"},
			{Name: "RX_DATA_COMPLETE", Src: []string{"data"}, Dst: "data_complete"},
			{Name: "RX_CHECKSUM", Src: []string{"data_complete"}, Dst: "checksum"},
			{Name: "CRC_OK", Src: []string{"checksum"}, Dst: "idle"},
			{Name: "CRC_NOTOK", Src: []string{"checksum"}, Dst: "idle"},
		},
		fsm.Callbacks{
			"enter_idle": func(e *fsm.Event) {
				frameParser.parseTimeout.Stop()
				frameParser.payloadReadBuffer.Reset()
			},
			"PARSE_TIMEOUT": func(e *fsm.Event) {
				event := &FrameParseEvent{
					status: FrameParseTimeout,
					frame:  Frame{},
				}

				go func() {
					frameParser.framesReceived <- event
				}()
			},
			"RX_ACK": func(e *fsm.Event) {
				frameParser.acks <- true
			},
			"RX_NAK": func(e *fsm.Event) {
				frameParser.naks <- true
			},
			"RX_CAN": func(e *fsm.Event) {
				frameParser.cans <- true
			},
			"RX_SOF": func(e *fsm.Event) {
				frameParser.sof = e.Args[0].(uint8)
				frameParser.parseTimeout.Reset(readTimeout)
			},
			"RX_LENGTH": func(e *fsm.Event) {
				frameParser.length = e.Args[0].(uint8)
				frameParser.readCounter = frameParser.length - 2
			},
			"RX_DATA": func(e *fsm.Event) {
				frameParser.payloadReadBuffer.WriteByte(e.Args[0].(uint8))
				frameParser.readCounter--
			},
			"checksum": func(e *fsm.Event) {
				e.Async()
			},
			"CRC_OK": func(e *fsm.Event) {
				event := &FrameParseEvent{
					status: FrameParseOk,
					frame:  e.Args[0].(Frame),
				}

				go func() {
					frameParser.framesReceived <- event
				}()
			},
			"CRC_NOTOK": func(e *fsm.Event) {
				event := &FrameParseEvent{
					status: FrameParseNotOk,
					frame:  e.Args[0].(Frame),
				}

				go func() {
					frameParser.framesReceived <- event
				}()
			},
			// "before_event": func(e *fsm.Event) {
			// 	if e.Src == "data" && e.Dst == "data" {
			// 		return
			// 	}
			// 	fmt.Printf("%s: %s -> %s\n", e.Event, e.Src, e.Dst)
			// },
		},
	)

	go frameParser.parse()

	return frameParser
}

func (parser *FrameParser) parse() {
	for {
		select {
		case <-parser.parseTimeout.C:
			parser.state.Event("PARSE_TIMEOUT")

		case currentByte := <-parser.input:
			parser.processByte(currentByte)
		}
	}
}

func (parser *FrameParser) processByte(currentByte byte) {
	switch {

	case parser.state.Is("idle"):
		switch currentByte {
		case FrameHeaderData:
			parser.state.Event("RX_SOF", currentByte)
		case FrameHeaderAck:
			parser.state.Event("RX_ACK", currentByte)
		case FrameHeaderCan:
			parser.state.Event("RX_CAN", currentByte)
		case FrameHeaderNak:
			parser.state.Event("RX_NAK", currentByte)
		}

	case parser.state.Is("length"):
		if currentByte < minFrameSize || currentByte > maxFrameSize {
			parser.state.Event("INVALID_LENGTH")
		} else {
			parser.state.Event("RX_LENGTH", currentByte)
		}

	case parser.state.Is("data"):
		if parser.readCounter > 0 {
			parser.state.Event("RX_DATA", currentByte)
		} else {
			parser.state.Event("RX_DATA", currentByte)
			parser.state.Event("RX_DATA_COMPLETE")
		}

	case parser.state.Is("data_complete"):
		parser.state.Event("RX_CHECKSUM", currentByte)
		parser.state.Transition()

		payload := parser.payloadReadBuffer.Bytes()
		frame := Frame{
			Header:   parser.sof,
			Length:   parser.length,
			Type:     payload[0],
			Payload:  payload[1:],
			Checksum: currentByte,
		}

		if frame.VerifyChecksum() == nil {
			parser.state.Event("CRC_OK", frame)
		} else {
			parser.state.Event("CRC_NOTOK", frame)
		}

	}
}
