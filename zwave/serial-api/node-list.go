package serialapi

import (
	"errors"

	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/session"
)

type NodeListResponse struct {
	CommandId    byte
	Version      byte
	Capabilities byte
	Nodes        []byte
	ChipType     byte
	ChipVersion  byte
}

func (s *SerialAPILayer) GetNodeList() (*NodeListResponse, error) {

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

	return &NodeListResponse{
		CommandId:    ret.Payload[0],
		Version:      ret.Payload[1],
		Capabilities: ret.Payload[2],
		Nodes:        ret.Payload[4:33],
		ChipType:     ret.Payload[33],
		ChipVersion:  ret.Payload[34],
	}, nil

}

func isBitSet(mask []byte, nodeId uint8) bool {
	if (nodeId > 0) && (nodeId <= 232) {
		return ((mask[(nodeId-1)>>3] & (1 << ((nodeId - 1) & 0x07))) != 0)
	}

	return false
}

func (n *NodeListResponse) GetApiType() string {
	if n.CommandId&0x80 == 0x80 {
		return "Slave"
	} else {
		return "Controller"
	}
}

func (n *NodeListResponse) TimerFunctionsSupported() bool {
	if n.CommandId&0x40 == 0x40 {
		return true
	} else {
		return false
	}
}

func (n *NodeListResponse) IsPrimaryController() bool {
	if n.CommandId&0x20 == 0x20 {
		return false
	} else {
		return true
	}
}

func (n *NodeListResponse) GetNodeIds() []uint8 {
	nodes := []uint8{}

	var i uint8
	for i = 1; i <= 232; i++ {
		if isBitSet(n.Nodes, i) {
			nodes = append(nodes, i)
		}
	}

	return nodes
}
