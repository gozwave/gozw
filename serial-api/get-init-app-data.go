package serialapi

import (
	"errors"

	"github.com/helioslabs/gozw/frame"
	"github.com/helioslabs/gozw/protocol"
	"github.com/helioslabs/gozw/session"
)

type InitAppData struct {
	CommandID    byte
	Version      byte
	Capabilities byte
	Nodes        []byte
	ChipType     byte
	ChipVersion  byte
}

func (s *Layer) GetInitAppData() (*InitAppData, error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionID: protocol.FnSerialAPIGetInitAppData,
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
		CommandID:    ret.Payload[0],
		Version:      ret.Payload[1],
		Capabilities: ret.Payload[2],
		Nodes:        ret.Payload[4:33],
		ChipType:     ret.Payload[33],
		ChipVersion:  ret.Payload[34],
	}, nil

}

func isBitSet(mask []byte, nodeID byte) bool {
	if (nodeID > 0) && (nodeID <= 232) {
		return ((mask[(nodeID-1)>>3] & (1 << ((nodeID - 1) & 0x07))) != 0)
	}

	return false
}

func (n *InitAppData) GetAPIType() string {
	if n.CommandID&0x80 == 0x80 {
		return "Slave"
	}

	return "Controller"
}

func (n *InitAppData) TimerFunctionsSupported() bool {
	if n.CommandID&0x40 == 0x40 {
		return true
	}

	return false
}

func (n *InitAppData) IsPrimaryController() bool {
	if n.CommandID&0x20 == 0x20 {
		return false
	}

	return true
}

func (n *InitAppData) GetNodeIDs() []byte {
	nodes := []byte{}

	var i byte
	for i = 1; i <= 232; i++ {
		if isBitSet(n.Nodes, i) {
			nodes = append(nodes, i)
		}
	}

	return nodes
}
