package frame

import (
	"bytes"
	"context"
	"time"

	"github.com/looplab/fsm"
	"go.uber.org/zap"
)

// ParseStatus respresents an enum for the status parsing a frame.
type ParseStatus int

const (
	// ParseOk means a frame was successfully parsed.
	ParseOk ParseStatus = iota
	// ParseNotOk means a frame was not successfully parsed.
	ParseNotOk
	// ParseTimeout means a frame timed out.
	ParseTimeout

	minFrameSize uint8 = 3
	maxFrameSize uint8 = 88
	readTimeout        = 1500 * time.Millisecond
)

// ParseEvent represents a parsed frame and whether it was successful.
type ParseEvent struct {
	status ParseStatus
	frame  Frame
}

// Parser contains  a frame parser.
type Parser struct {
	state                              *fsm.FSM
	input                              <-chan byte
	framesReceived                     chan<- *ParseEvent
	acks, naks, cans                   chan<- bool
	sof, length, checksum, readCounter byte
	payloadReadBuffer                  *bytes.Buffer
	parseTimeout                       *time.Timer
	l                                  *zap.Logger
	ctx                                context.Context
}

// NewParser will return  a new  parser
func NewParser(ctx context.Context, input <-chan byte, output chan<- *ParseEvent, acks, naks, cans chan bool, logger *zap.Logger) *Parser {
	parser := &Parser{
		input:             input,
		framesReceived:    output,
		acks:              acks,
		naks:              naks,
		cans:              cans,
		payloadReadBuffer: bytes.NewBuffer([]byte{}),
		parseTimeout:      time.NewTimer(readTimeout),
		l:                 logger,
		ctx:               ctx,
	}

	parser.parseTimeout.Stop()

	parser.state = fsm.NewFSM(
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
				parser.parseTimeout.Stop()
				parser.payloadReadBuffer.Reset()
			},
			"PARSE_TIMEOUT": func(e *fsm.Event) {
				event := &ParseEvent{
					status: ParseTimeout,
					frame:  Frame{},
				}

				go func() {
					parser.framesReceived <- event
				}()
			},
			"RX_ACK": func(e *fsm.Event) {

				parser.acks <- true
			},
			"RX_NAK": func(e *fsm.Event) {
				parser.naks <- true
			},
			"RX_CAN": func(e *fsm.Event) {
				parser.cans <- true
			},
			"RX_SOF": func(e *fsm.Event) {
				parser.sof = e.Args[0].(byte)
				parser.parseTimeout.Reset(readTimeout)
			},
			"RX_LENGTH": func(e *fsm.Event) {
				parser.length = e.Args[0].(byte)
				parser.readCounter = parser.length - 2
			},
			"RX_DATA": func(e *fsm.Event) {
				parser.payloadReadBuffer.WriteByte(e.Args[0].(byte))
				parser.readCounter--
			},
			"checksum": func(e *fsm.Event) {
				e.Async()
			},
			"CRC_OK": func(e *fsm.Event) {
				parser.l.Debug("event received, attempting to parse")

				event := &ParseEvent{
					status: ParseOk,
					frame:  e.Args[0].(Frame),
				}

				parser.l.Debug("event parsed", zap.Int("status", int(event.status)))

				go func() {
					parser.framesReceived <- event
				}()
			},
			"CRC_NOTOK": func(e *fsm.Event) {
				parser.l.Debug("event received, attempting to parse")

				event := &ParseEvent{
					status: ParseNotOk,
					frame:  e.Args[0].(Frame),
				}

				parser.l.Debug("event parsed", zap.Int("status", int(event.status)))

				go func() {
					parser.framesReceived <- event
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

	go parser.parse()

	return parser
}

func (p *Parser) parse() {
	for {
		select {
		case <-p.parseTimeout.C:
			p.state.Event("PARSE_TIMEOUT")

		case currentByte := <-p.input:
			p.processByte(currentByte)
		case <-p.ctx.Done():
			p.l.Info("closing frame parser")
			return
		}
	}
}

func (p *Parser) processByte(currentByte byte) {
	// uncomment to really dig deep. this line is a little verbose
	// p.l.Debug("parser processing byte", zap.String("byte", fmt.Sprint(currentByte)), zap.String("current_state", p.state.Current()))

	switch {

	case p.state.Is("idle"):
		switch currentByte {
		case HeaderData:
			p.state.Event("RX_SOF", currentByte)
		case HeaderAck:
			p.state.Event("RX_ACK", currentByte)
		case HeaderCan:
			p.state.Event("RX_CAN", currentByte)
		case HeaderNak:
			p.state.Event("RX_NAK", currentByte)
		}

	case p.state.Is("length"):
		if currentByte < minFrameSize || currentByte > maxFrameSize {
			p.state.Event("INVALID_LENGTH")
		} else {
			p.state.Event("RX_LENGTH", currentByte)
		}

	case p.state.Is("data"):
		if p.readCounter > 0 {
			p.state.Event("RX_DATA", currentByte)
		} else {
			p.state.Event("RX_DATA", currentByte)
			p.state.Event("RX_DATA_COMPLETE")
		}

	case p.state.Is("data_complete"):
		p.state.Event("RX_CHECKSUM", currentByte)
		p.state.Transition()

		payload := p.payloadReadBuffer.Bytes()
		frame := Frame{
			Header:   p.sof,
			Length:   p.length,
			Type:     payload[0],
			Payload:  payload[1:],
			Checksum: currentByte,
		}

		p.l.Debug("data complete")

		if frame.VerifyChecksum() == nil {
			p.l.Debug("checksum verified")
			p.state.Event("CRC_OK", frame)
		} else {
			p.l.Debug("checksum invalid")
			p.state.Event("CRC_NOTOK", frame)
		}

	}
}
