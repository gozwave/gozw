package serialapi

import (
	"errors"

	"github.com/helioslabs/gozw/zwave/frame"
	"github.com/helioslabs/gozw/zwave/protocol"
	"github.com/helioslabs/gozw/zwave/session"
)

func (s *SerialAPILayer) GetSerialApiCapabilities() (*SerialApiCapabilities, error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionId: protocol.FnSerialApiGetCapabilities,
		HasReturn:  true,
		ReturnCallback: func(err error, ret *frame.Frame) bool {
			done <- ret
			return false
		},
	}

	s.sessionLayer.MakeRequest(request)
	ret := <-done

	if ret == nil {
		return nil, errors.New("Error getting home/node id")
	}

	val := &SerialApiCapabilities{
		ApplicationVersion:  ret.Payload[1],
		ApplicationRevision: ret.Payload[2],
		Manufacturer1:       ret.Payload[3],
		Manufacturer2:       ret.Payload[4],
		ProductType1:        ret.Payload[5],
		ProductType2:        ret.Payload[6],
		ProductId1:          ret.Payload[7],
		ProductId2:          ret.Payload[8],
		SupportedFunctions:  ret.Payload[9:],
	}

	return val, nil
}

type SerialApiCapabilities struct {
	ApplicationVersion  byte
	ApplicationRevision byte
	Manufacturer1       byte
	Manufacturer2       byte
	ProductType1        byte
	ProductType2        byte
	ProductId1          byte
	ProductId2          byte
	SupportedFunctions  []byte
}

func (n *SerialApiCapabilities) GetSupportedFunctions() []byte {
	supportedFunctions := []byte{}

	var i byte
	for i = 1; i < 255; i++ {
		if isBitSet(n.SupportedFunctions, i) {
			supportedFunctions = append(supportedFunctions, i)
		}
	}

	return supportedFunctions
}
