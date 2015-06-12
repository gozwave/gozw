package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type InitData struct {
	layers.BaseLayer
	Version      uint8
	Capabilities uint8
	Nodes        []byte
	ChipType     uint8
	ChipVersion  uint8
}

func (f *InitData) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	f.Version = data[0]
	f.Capabilities = data[1]
	f.Nodes = data[4:33]
	f.ChipType = data[33]
	f.ChipVersion = data[34]
	return nil
}

func (f *InitData) CanDecode() gopacket.LayerClass {
	return LayerTypeAddNode
}

func (f *InitData) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypeZero
}

func (f *InitData) LayerPayload() []byte {
	return []byte{f.Options, f.CallbackId}
}

func (f *InitData) LayerType() gopacket.LayerType {
	return LayerTypeAddNode
}

func decodeAddNode(data []byte, p gopacket.PacketBuilder) error {
	addNode := &InitData{}
	err := addNode.DecodeFromBytes(data, p)

	p.AddLayer(addNode)

	if err != nil {
		return err
	}

	return p.NextDecoder(nil)
}
