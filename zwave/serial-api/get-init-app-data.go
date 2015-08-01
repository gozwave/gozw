package serialapi

import (
	"errors"

	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/session"
)

type InitAppData struct {
	CommandId    byte
	Version      byte
	Capabilities byte
	Nodes        []byte
	ChipType     byte
	ChipVersion  byte
}

func (s *SerialAPILayer) GetInitAppData() (*InitAppData, error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionId: protocol.FnSerialApiGetInitAppData,
		HasReturn:  true,
		ReturnCallback: func(err error, ret *frame.Frame) bool {
			done <- ret
			return false
		},
	}

	s.sessionLayer.MakeRequest(request)
	ret := <-done

	if ret == nil {
		return nil, errors.New("Error getting node information")
	}

	return &InitAppData{
		CommandId:    ret.Payload[0],
		Version:      ret.Payload[1],
		Capabilities: ret.Payload[2],
		Nodes:        ret.Payload[4:33],
		ChipType:     ret.Payload[33],
		ChipVersion:  ret.Payload[34],
	}, nil

}

func isBitSet(mask []byte, nodeId byte) bool {
	if (nodeId > 0) && (nodeId <= 232) {
		return ((mask[(nodeId-1)>>3] & (1 << ((nodeId - 1) & 0x07))) != 0)
	}

	return false
}

func (n *InitAppData) GetApiType() string {
	if n.CommandId&0x80 == 0x80 {
		return "Slave"
	} else {
		return "Controller"
	}
}

func (n *InitAppData) TimerFunctionsSupported() bool {
	if n.CommandId&0x40 == 0x40 {
		return true
	} else {
		return false
	}
}

func (n *InitAppData) IsPrimaryController() bool {
	if n.CommandId&0x20 == 0x20 {
		return false
	} else {
		return true
	}
}

func (n *InitAppData) GetNodeIds() []byte {
	nodes := []byte{}

	var i byte
	for i = 1; i <= 232; i++ {
		if isBitSet(n.Nodes, i) {
			nodes = append(nodes, i)
		}
	}

	return nodes
}
