package frame

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	HeaderData byte = 0x01
	HeaderAck       = 0x06
	HeaderNak       = 0x15
	HeaderCan       = 0x18
)

const (
	TypeRequest  byte = 0x00
	TypeResponse      = 0x01
)

type Frame struct {

	// Header is one of FrameHeader*
	Header byte

	// Length = byte length of all fields, excluding Header and Checksum
	Length byte

	// Type is one of FrameType*
	Type byte

	// Payload is the command id and command payload
	Payload []byte

	// Checksum = 0xff XOR Type XOR Length XOR payload[0] XOR [...payload[n]]
	Checksum byte
}

func NewRequestFrame(payload []byte) *Frame {
	return &Frame{
		Header:  HeaderData,
		Type:    TypeRequest,
		Payload: payload,
	}
}

func NewNakFrame() *Frame {
	return &Frame{
		Header: HeaderNak,
	}
}

func NewAckFrame() *Frame {
	return &Frame{
		Header: HeaderAck,
	}
}

func NewCanFrame() *Frame {
	return &Frame{
		Header: HeaderCan,
	}
}

func (z *Frame) IsRequest() bool {
	return z.Type == TypeRequest
}

func (z *Frame) IsResponse() bool {
	return z.Type == TypeResponse
}

func (z *Frame) IsAck() bool {
	return z.Header == HeaderAck
}

func (z *Frame) IsNak() bool {
	return z.Header == HeaderNak
}

func (z *Frame) IsCan() bool {
	return z.Header == HeaderCan
}

func (z *Frame) IsData() bool {
	return z.Header == HeaderData
}

// CalcChecksum calculates the checksum for this frame, given the current data.
// The Z-Wave checksum is calculated by taking 0xFF XOR Length XOR Type XOR Payload[0:n]
func (z *Frame) CalcChecksum() byte {
	var csum byte = 0xff
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
	if z.Header != HeaderData {
		return nil
	}

	if z.Checksum != z.CalcChecksum() {
		return errors.New("Invalid checksum")
	}

	return nil
}

// Marshal this frame into a byte slice
func (z *Frame) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	z.Length = byte(len(z.Payload) + 2)
	z.SetChecksum()

	switch z.Header {
	case HeaderData:
		// Data frames have the whole kit and caboodle
		binary.Write(buf, binary.BigEndian, byte(z.Header))
		binary.Write(buf, binary.BigEndian, byte(z.Length))
		binary.Write(buf, binary.BigEndian, byte(z.Type))
		buf.Write(z.Payload)
		binary.Write(buf, binary.BigEndian, byte(z.Checksum))
	default:
		// Non-data frames are just a single byte
		binary.Write(buf, binary.BigEndian, byte(z.Header))
	}

	return buf.Bytes(), nil
}

// UnmarshalFrame turns a byte slice into a Frame
func UnmarshalFrame(frame []byte) *Frame {
	if frame[0] != HeaderData {
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
