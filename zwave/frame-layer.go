package zwave

import (
	"bytes"
	"fmt"
	"time"
)

type frameParseState int

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

const (
	FRS_SOF_HUNT frameParseState = iota
	FRS_LENGTH
	FRS_TYPE
	FRS_DATA
	FRS_CHECKSUM
)

type FrameLayer struct {
	transport  *TransportLayer
	parseState frameParseState
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
	frames := make(chan *Frame)
	go f.readFromTransport(frames)
	return frames
}

func (f *FrameLayer) sendAck() error {
	_, err := f.transport.Write([]byte{FrameSOFAck})
	return err
}

func (f *FrameLayer) sendNak() error {
	_, err := f.transport.Write([]byte{FrameSOFNak})
	return err
}

func (f *FrameLayer) readFromTransport(frames chan<- *Frame) {

	f.parseState = FRS_SOF_HUNT

	timeout := time.NewTimer(readTimeout)
	timeout.Stop()

	var sof, length, frameType, counter uint8
	payload := bytes.NewBuffer([]byte{})

	inputBytes := f.transport.Read()

	for {
		select {
		case <-timeout.C:
			f.parseState = FRS_SOF_HUNT
			payload.Reset()

		case currentByte := <-inputBytes:
			switch f.parseState {

			case FRS_SOF_HUNT:
				switch currentByte {
				case FrameSOFData:
					sof = currentByte
					f.parseState = FRS_LENGTH
					timeout.Reset(readTimeout)

				case FrameSOFAck:
					f.parseState = FRS_SOF_HUNT
					// @todo make ACK channel
				case FrameSOFCan:
					f.parseState = FRS_SOF_HUNT
					// @todo make CAN channel
				case FrameSOFNak:
					f.parseState = FRS_SOF_HUNT
					// @todo make NAK channel
				}

			case FRS_LENGTH:
				length = currentByte
				if length < minFrameSize || length > maxFrameSize {
					f.parseState = FRS_SOF_HUNT
				} else {
					counter = length - 2
					f.parseState = FRS_TYPE
				}

			case FRS_TYPE:
				frameType = currentByte
				counter--
				f.parseState = FRS_DATA

			case FRS_DATA:
				if counter > 0 {
					payload.WriteByte(currentByte)
					counter--
				} else {
					f.parseState = FRS_CHECKSUM
				}

			case FRS_CHECKSUM:
				f.parseState = FRS_SOF_HUNT

				frame := &Frame{
					Header:   sof,
					Length:   length,
					Type:     frameType,
					Payload:  payload.Bytes(),
					Checksum: currentByte,
				}

				payload.Reset()

				if frame.VerifyChecksum() == nil {
					fmt.Println("yo")
					frames <- frame
				} else {
					fmt.Println("invalid frame:", frame)
					// @todo send NAK
				}

			}
		}
	}
}
