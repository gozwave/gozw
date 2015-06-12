package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type AddNode struct {
	layers.BaseLayer
	Options    uint8
	CallbackId uint8
}

func (an *AddNode) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	an.Options = data[0]
	an.CallbackId = data[1]
	return nil
}

func (an *AddNode) CanDecode() gopacket.LayerClass {
	return LayerTypeAddNode
}

func (an *AddNode) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypeZero
}

func (an *AddNode) LayerPayload() []byte {
	return []byte{an.Options, an.CallbackId}
}

func (an *AddNode) LayerType() gopacket.LayerType {
	return LayerTypeAddNode
}

func decodeAddNode(data []byte, p gopacket.PacketBuilder) error {
	addNode := &AddNode{}
	err := addNode.DecodeFromBytes(data, p)

	p.AddLayer(addNode)

	if err != nil {
		return err
	}

	return p.NextDecoder(nil)
}
