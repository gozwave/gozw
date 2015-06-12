package layers

import "github.com/google/gopacket"

var (
	LayerTypeFrame = gopacket.RegisterLayerType(4000, gopacket.LayerTypeMetadata{
		Name:    "Frame",
		Decoder: gopacket.DecodeFunc(decodeFrame),
	})

	LayerTypeAddNode = gopacket.RegisterLayerType(4001, gopacket.LayerTypeMetadata{
		Name:    "FnAddNode",
		Decoder: gopacket.DecodeFunc(decodeAddNode),
	})

	LayerTypeInitData = gopacket.RegisterLayerType(4002, gopacket.LayerTypeMetadata{
		Name:    "FnInitData",
		Decoder: gopacket.DecodeFunc(decodeInitData),
	})
)

var LayerClassFrame = gopacket.NewLayerClass([]gopacket.LayerType{
	LayerTypeFrame,
})

var LayerClassFunction = gopacket.NewLayerClass([]gopacket.LayerType{
	LayerTypeAddNode,
})
