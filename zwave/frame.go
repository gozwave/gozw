package zwave

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type FrameHeader uint8
type FrameType uint8

const (
	FrameHeaderData uint8 = 0x01
	FrameHeaderAck  uint8 = 0x06
	FrameHeaderNak  uint8 = 0x15
	FrameHeaderCan  uint8 = 0x18
)

type Frame struct {

	// Header is one of FrameHeader*
	Header uint8

	// Length = byte length of all fields, excluding Header and Checksum
	Length uint8

	// Type is one of FrameType*
	Type uint8

	// Payload is the command id and command payload
	Payload []byte

	// Checksum = 0xff XOR Type XOR Length XOR payload[0] XOR [...payload[n]]
	Checksum uint8
}

func NewRequestFrame() *Frame {
	return &Frame{
		Header: FrameHeaderData,
		Type:   FrameTypeReq,
	}
}

func NewNakFrame() *Frame {
	return &Frame{
		Header: FrameHeaderNak,
	}
}

func NewAckFrame() *Frame {
	return &Frame{
		Header: FrameHeaderAck,
	}
}

func NewCanFrame() *Frame {
	return &Frame{
		Header: FrameHeaderCan,
	}
}

func (z *Frame) IsRequest() bool {
	return z.Type == FrameTypeReq
}

func (z *Frame) IsResponse() bool {
	return z.Type == FrameTypeRes
}

func (z *Frame) IsAck() bool {
	return z.Header == FrameHeaderAck
}

func (z *Frame) IsNak() bool {
	return z.Header == FrameHeaderNak
}

func (z *Frame) IsCan() bool {
	return z.Header == FrameHeaderCan
}

func (z *Frame) IsData() bool {
	return z.Header == FrameHeaderData
}

// CalcChecksum calculates the checksum for this frame, given the current data.
// The Z-Wave checksum is calculated by taking 0xFF XOR Length XOR Type XOR Payload[0:n]
func (z *Frame) CalcChecksum() uint8 {
	var csum uint8 = 0xff
	csum ^= z.Length
	csum ^= z.Type

	for i := 0; i < len(z.Payload); i++ {
		csum ^= z.Payload[i]
	}

	return csum
}

// SetChecksum calculates the frame checksum and saves it into the frame
func (z *Frame) SetChecksum() {
	z.Checksum = z.CalcChecksum()
}

// VerifyChecksum calculates a checksum for the frame and compares it to the
// frame's checksum, returning an error if they do not agree
func (z *Frame) VerifyChecksum() error {
	if z.Header != FrameHeaderData {
		return nil
	}

	if z.Checksum != z.CalcChecksum() {
		return errors.New("Invalid checksum")
	}

	return nil
}

// Marshal this frame into a byte slice
func (z *Frame) Marshal() []byte {
	buf := new(bytes.Buffer)

	z.Length = uint8(len(z.Payload) + 2)
	z.SetChecksum()

	switch z.Header {
	case FrameHeaderData:
		// Data frames have the whole kit and caboodle
		binary.Write(buf, binary.BigEndian, uint8(z.Header))
		binary.Write(buf, binary.BigEndian, uint8(z.Length))
		binary.Write(buf, binary.BigEndian, uint8(z.Type))
		buf.Write(z.Payload)
		binary.Write(buf, binary.BigEndian, uint8(z.Checksum))
	default:
		// Non-data frames are just a single byte
		binary.Write(buf, binary.BigEndian, uint8(z.Header))
	}

	return buf.Bytes()
}

// UnmarshalFrame turns a byte slice into a Frame
func UnmarshalFrame(frame []byte) *Frame {
	if frame[0] != FrameHeaderData {
		return &Frame{
			Header: frame[0],
		}
	}

	return &Frame{
		Header:   frame[0],
		Length:   frame[1],
		Type:     frame[2],
		Payload:  frame[3 : len(frame)-1],
		Checksum: frame[len(frame)-1],
	}
}
