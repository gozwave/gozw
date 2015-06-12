package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const (
	FrameSOFData uint8 = 0x01
	FrameSOFAck  uint8 = 0x06
	FrameSOFNak  uint8 = 0x15
	FrameSOFCan  uint8 = 0x18
)

const (
	FrameTypeReq uint8 = 0x00
	FrameTypeRes uint8 = 0x01
)

type Frame struct {
	layers.BaseLayer

	// SOF is one of FrameSOF*
	SOF uint8

	// Length = byte length of all fields, excluding Header and Checksum
	Length uint8

	// Type is one of FrameType*
	Type uint8

	// FunctionId is a Z-Wave API Function ID
	FunctionId uint8

	// Payload is the command id and command payload
	Payload []byte

	// Checksum = 0xff XOR Type XOR Length XOR payload[0] XOR [...payload[n]]
	Checksum uint8
}

func (f *Frame) IsAck() bool {
	return f.SOF == FrameSOFAck
}

func (f *Frame) IsNak() bool {
	return f.SOF == FrameSOFNak
}

func (f *Frame) IsCan() bool {
	return f.SOF == FrameSOFCan
}

func (f *Frame) IsData() bool {
	return f.SOF == FrameSOFData
}

func (f *Frame) IsRequest() bool {
	return f.Type == FrameTypeReq
}

func (f *Frame) IsResponse() bool {
	return f.Type == FrameTypeRes
}

// CalcChecksum calculates the checksum for this frame, given the current data.
// The Z-Wave checksum is calculated by taking 0xFF XOR Length XOR Type XOR Payload[0:n]
func (f *Frame) CalcChecksum() uint8 {
	var csum uint8 = 0xff
	csum ^= f.Length
	csum ^= f.Type

	for i := 0; i < len(f.Payload); i++ {
		csum ^= f.Payload[i]
	}

	return csum
}

// Layer interface
func (f *Frame) LayerType() gopacket.LayerType {
	return LayerTypeFrame
}

// DecodingLayer interface
func (f *Frame) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	f.SOF = data[0]

	if !f.IsData() {
		return nil
	}

	f.Length = data[1]
	f.Type = data[2]
	f.FunctionId = data[3]

	f.Payload = data[3:f.Length]
	f.Checksum = data[f.Length+1]

	return nil
}

func (f *Frame) CanDecode() gopacket.LayerClass {
	return LayerClassFrame
}

func (f *Frame) NextLayerType() gopacket.LayerType {
	switch f.FunctionId {
	case ZwGetInitData:
		return LayerTypeInitData
	case ZwAddNodeToNetwork:
		return LayerTypeAddNode
	default:
		return gopacket.LayerTypePayload
	}
}

func (f *Frame) LayerPayload() []byte {
	return f.Payload
}

func (f *Frame) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	bytes, err := b.PrependBytes(5 + len(f.Payload))
	if err != nil {
		return err
	}

	bytes[0] = f.SOF
	bytes[1] = byte(len(f.Payload) + 2)
	bytes[2] = f.Type
	bytes[3] = f.FunctionId

	for i := 0; i < len(f.Payload); i++ {
		bytes[i+4] = f.Payload[i]
	}

	if opts.ComputeChecksums {
		bytes[len(f.Payload)+4] = f.CalcChecksum()
	} else {
		bytes[len(f.Payload)+4] = 0
	}

	return nil
}

func decodeFrame(data []byte, p gopacket.PacketBuilder) error {
	frame := &Frame{}
	err := frame.DecodeFromBytes(data, p)

	p.AddLayer(frame)

	if err != nil {
		return err
	}

	return p.NextDecoder(frame.NextLayerType())
}
