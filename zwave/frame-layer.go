package zwave

import (
	"bytes"
	"time"

	"github.com/looplab/fsm"
)

const (
	minFrameSize uint8 = 3
	maxFrameSize uint8 = 88
)

const (
	readTimeout time.Duration = 1500 * time.Millisecond
	ackTimeCan1 time.Duration = 500 * time.Millisecond
	ackTimeCan2 time.Duration = 500 * time.Millisecond
	ackTime1    time.Duration = 1000 * time.Millisecond
	ackTime2    time.Duration = 2000 * time.Millisecond
)

type FrameLayer struct {
	transport                          *TransportLayer
	parseState                         *fsm.FSM
	sof, length, checksum, readCounter uint8
	payloadReadBuffer                  *bytes.Buffer
	parseTimeout                       *time.Timer
}

func NewFrameLayer(transport *TransportLayer) *FrameLayer {
	frameLayer := &FrameLayer{
		transport: transport,
	}

	return frameLayer
}

func (f *FrameLayer) Write(bytes []byte) {
	f.transport.Write(bytes)
}

func (f *FrameLayer) Read() <-chan *Frame {

	parsedFrameQueue := make(chan *Frame)

	f.payloadReadBuffer = bytes.NewBuffer([]byte{})

	f.parseTimeout = time.NewTimer(readTimeout)
	f.parseTimeout.Stop()

	f.parseState = fsm.NewFSM(
		"idle",
		fsm.Events{
			{Name: "PARSE_TIMEOUT", Src: []string{"idle", "length", "data", "checksum"}, Dst: "idle"},
			{Name: "RX_SOF", Src: []string{"idle"}, Dst: "length"},
			{Name: "RX_ACK", Src: []string{"idle"}, Dst: "idle"},
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
				f.parseTimeout.Stop()
				f.payloadReadBuffer.Reset()
				// stop timeout, clear payload buffer
			},
			"RX_SOF": func(e *fsm.Event) {
				f.parseTimeout.Reset(readTimeout)
			},
			"RX_LENGTH": func(e *fsm.Event) {
				f.length = e.Args[0].(uint8)
				f.readCounter = f.length - 2
			},
			"RX_DATA": func(e *fsm.Event) {
				f.payloadReadBuffer.WriteByte(e.Args[0].(uint8))
				f.readCounter--
			},
			"checksum": func(e *fsm.Event) {
				e.Async()
			},
			"CRC_OK": func(e *fsm.Event) {
				f.sendAck()
				parsedFrameQueue <- e.Args[0].(*Frame)
			},
			"CRC_NOTOK": func(e *fsm.Event) {
				// @todo logging?
				f.sendNak()
			},
			// "before_event": func(e *fsm.Event) {
			// 	fmt.Printf("%s: %s -> %s\n", e.Event, e.Src, e.Dst)
			// },
		},
	)

	go f.readLoop()
	return parsedFrameQueue
}

func (f *FrameLayer) sendAck() error {
	_, err := f.transport.Write([]byte{FrameSOFAck})
	return err
}

func (f *FrameLayer) sendNak() error {
	_, err := f.transport.Write([]byte{FrameSOFNak})
	return err
}

func (f *FrameLayer) readLoop() {

	timeout := time.NewTimer(readTimeout)
	timeout.Stop()

	readQueue := f.transport.Read()

	for {
		select {
		case <-timeout.C:
			f.parseState.Event("PARSE_TIMEOUT")

		case currentByte := <-readQueue:
			switch {

			case f.parseState.Is("idle"):
				switch currentByte {
				case FrameSOFData:
					f.parseState.Event("RX_SOF", currentByte)

				case FrameSOFAck:
					// @todo make ACK channel
				case FrameSOFCan:
					// @todo make CAN channel
				case FrameSOFNak:
					// @todo make NAK channel
				}

			case f.parseState.Is("length"):
				if currentByte < minFrameSize || currentByte > maxFrameSize {
					f.parseState.Event("INVALID_LENGTH")
				} else {
					f.parseState.Event("RX_LENGTH", currentByte)
				}

			case f.parseState.Is("data"):
				if f.readCounter > 0 {
					f.parseState.Event("RX_DATA", currentByte)
				} else {
					f.parseState.Event("RX_DATA_COMPLETE")
				}

			case f.parseState.Is("data_complete"):
				f.parseState.Event("RX_CHECKSUM", currentByte)
				f.parseState.Transition()

				payload := f.payloadReadBuffer.Bytes()
				frame := &Frame{
					Header:   f.sof,
					Length:   f.length,
					Type:     payload[0],
					Payload:  payload[1:],
					Checksum: currentByte,
				}

				if frame.VerifyChecksum() == nil {
					f.parseState.Event("CRC_OK", frame)
				} else {
					f.parseState.Event("CRC_NOTOK")
				}

			}
		}
	}
}
