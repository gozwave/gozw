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

const (
	FrameTypeReq uint8 = 0x00
	FrameTypeRes uint8 = 0x01
)

type ZFrame struct {

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

func NewRequestFrame(payload []byte) *ZFrame {
	frame := ZFrame{
		Header:  FrameHeaderData,
		Type:    FrameTypeReq,
		Length:  uint8(len(payload) + 2), // payload length plus Type and Length
		Payload: payload,
	}

	frame.SetChecksum()

	return &frame
}

func NewNakFrame() *ZFrame {
	return &ZFrame{
		Header: FrameHeaderNak,
	}
}

func NewAckFrame() *ZFrame {
	return &ZFrame{
		Header: FrameHeaderAck,
	}
}

func NewCanFrame() *ZFrame {
	return &ZFrame{
		Header: FrameHeaderCan,
	}
}

func (z *ZFrame) IsRequest() bool {
	return z.Type == FrameTypeReq
}

func (z *ZFrame) IsResponse() bool {
	return z.Type == FrameTypeRes
}

func (z *ZFrame) IsAck() bool {
	return z.Header == FrameHeaderAck
}

// CalcChecksum calculates the checksum for this frame, given the current data.
// The Z-Wave checksum is calculated by taking 0xFF XOR Length XOR Type XOR Payload[0:n]
func (z *ZFrame) CalcChecksum() uint8 {
	var csum uint8 = 0xff
	csum ^= z.Length
	csum ^= z.Type

	for i := 0; i < len(z.Payload); i++ {
		csum ^= z.Payload[i]
	}

	return csum
}

// SetChecksum calculates the frame checksum and saves it into the frame
func (z *ZFrame) SetChecksum() {
	z.Checksum = z.CalcChecksum()
}

// VerifyChecksum calculates a checksum for the frame and compares it to the
// frame's checksum, returning an error if they do not agree
func (z *ZFrame) VerifyChecksum() error {
	if z.Header != FrameHeaderData {
		return nil
	}

	if z.Checksum != z.CalcChecksum() {
		return errors.New("Invalid checksum")
	}

	return nil
}

// Marshal this frame into a byte slice
func (z *ZFrame) Marshal() []byte {
	buf := new(bytes.Buffer)

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

// UnmarshalFrame turns a byte slice into a ZFrame
func UnmarshalFrame(frame []byte) *ZFrame {
	if frame[0] != FrameHeaderData {
		return &ZFrame{
			Header: frame[0],
		}
	}

	return &ZFrame{
		Header:   frame[0],
		Length:   frame[1],
		Type:     frame[2],
		Payload:  frame[3 : len(frame)-2],
		Checksum: frame[len(frame)-1],
	}
}
