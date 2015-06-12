package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type GetInitData struct {
	layers.BaseLayer
	Version      uint8
	Capabilities uint8
	Nodes        []byte
	ChipType     uint8
	ChipVersion  uint8
}

func (f *GetInitData) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	f.Version = data[0]
	f.Capabilities = data[1]
	f.Nodes = data[4:33]
	f.ChipType = data[33]
	f.ChipVersion = data[34]
	return nil
}

func (f *GetInitData) CanDecode() gopacket.LayerClass {
	return LayerTypeGetInitData
}

func (f *GetInitData) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypeZero
}

func (f *GetInitData) LayerPayload() []byte {
	buf := []byte{f.Version, f.Capabilities, 29}
	buf = append(buf, f.Nodes...)
	buf = append(buf, f.ChipType, f.ChipVersion)
	return buf
}

func (f *GetInitData) LayerType() gopacket.LayerType {
	return LayerTypeGetInitData
}

func decodeGetInitData(data []byte, p gopacket.PacketBuilder) error {
	packet := &GetInitData{}
	err := packet.DecodeFromBytes(data, p)

	p.AddLayer(packet)

	if err != nil {
		return err
	}

	return p.NextDecoder(nil)
}
