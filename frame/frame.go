package frame

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	// HeaderData is a byte for a data frame.
	HeaderData byte = 0x01
	// HeaderAck is a byte for a ack frame.
	HeaderAck = 0x06
	// HeaderNak is a byte for a nak frame.
	HeaderNak = 0x15
	// HeaderCan is a byte for a can frame.
	HeaderCan = 0x18

	// TypeRequest is for a request frame.
	TypeRequest byte = 0x00
	// TypeResponse is for a response frame.
	TypeResponse = 0x01
)

// Frame contains a frame.
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

// NewRequestFrame will build  a new request frame.
func NewRequestFrame(payload []byte) *Frame {
	return &Frame{
		Header:  HeaderData,
		Type:    TypeRequest,
		Payload: payload,
	}
}

// NewNakFrame returns a new  nak frame.
func NewNakFrame() *Frame {
	return &Frame{
		Header: HeaderNak,
	}
}

// NewAckFrame returns a new ack frame.
func NewAckFrame() *Frame {
	return &Frame{
		Header: HeaderAck,
	}
}

// NewCanFrame returns a new can frame.
func NewCanFrame() *Frame {
	return &Frame{
		Header: HeaderCan,
	}
}

// IsRequest checks if the frame is a request frame.
func (z *Frame) IsRequest() bool {
	return z.Type == TypeRequest
}

// IsResponse checks if the frame is a response frame.
func (z *Frame) IsResponse() bool {
	return z.Type == TypeResponse
}

// IsAck returns whether this is an ack frame.
func (z *Frame) IsAck() bool {
	return z.Header == HeaderAck
}

// IsNak returns whether this is an nak frame.
func (z *Frame) IsNak() bool {
	return z.Header == HeaderNak
}

// IsCan returns whether this is an can frame.
func (z *Frame) IsCan() bool {
	return z.Header == HeaderCan
}

// IsData returns whether this is a data frame.
func (z *Frame) IsData() bool {
	return z.Header == HeaderData
}

// CalcChecksum calculates the checksum for this frame, given the current data.
// The Z-Wave checksum is calculated by taking 0xFF XOR Length XOR Type XOR Payload[0:n].
func (z *Frame) CalcChecksum() byte {
	var csum byte = 0xff
	csum ^= z.Length
	csum ^= z.Type

	for i := 0; i < len(z.Payload); i++ {
		csum ^= z.Payload[i]
	}

	return csum
}

// SetChecksum calculates the frame checksum and saves it into the frame.
func (z *Frame) SetChecksum() {
	z.Checksum = z.CalcChecksum()
}

// VerifyChecksum calculates a checksum for the frame and compares it to the
// frame's checksum, returning an error if they do not agree.
func (z *Frame) VerifyChecksum() error {
	if z.Header != HeaderData {
		return nil
	}

	if z.Checksum != z.CalcChecksum() {
		return errors.New("Invalid checksum")
	}

	return nil
}

// MarshalBinary will marshal this frame into a byte slice.
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

// UnmarshalFrame turns a byte slice into a Frame.
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
