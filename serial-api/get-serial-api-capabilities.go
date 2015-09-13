package serialapi

import (
	"errors"

	"github.com/helioslabs/gozw/frame"
	"github.com/helioslabs/gozw/protocol"
	"github.com/helioslabs/gozw/session"
)

func (s *Layer) GetSerialAPICapabilities() (*SerialAPICapabilities, error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionID: protocol.FnSerialAPIGetCapabilities,
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

	val := &SerialAPICapabilities{
		ApplicationVersion:  ret.Payload[1],
		ApplicationRevision: ret.Payload[2],
		Manufacturer1:       ret.Payload[3],
		Manufacturer2:       ret.Payload[4],
		ProductType1:        ret.Payload[5],
		ProductType2:        ret.Payload[6],
		ProductID1:          ret.Payload[7],
		ProductID2:          ret.Payload[8],
		SupportedFunctions:  ret.Payload[9:],
	}

	return val, nil
}

type SerialAPICapabilities struct {
	ApplicationVersion  byte
	ApplicationRevision byte
	Manufacturer1       byte
	Manufacturer2       byte
	ProductType1        byte
	ProductType2        byte
	ProductID1          byte
	ProductID2          byte
	SupportedFunctions  []byte
}

func (s *SerialAPICapabilities) GetSupportedFunctions() []byte {
	supportedFunctions := []byte{}

	var i byte
	for i = 1; i < 255; i++ {
		if isBitSet(s.SupportedFunctions, i) {
			supportedFunctions = append(supportedFunctions, i)
		}
	}

	return supportedFunctions
}
