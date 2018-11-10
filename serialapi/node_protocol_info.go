package serialapi

import (
	"errors"

	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/protocol"
	"github.com/gozwave/gozw/session"
)

// GetNodeProtocolInfo will retrieve protocol info for a node.
func (s *Layer) GetNodeProtocolInfo(nodeID byte) (nodeInfo *NodeProtocolInfo, err error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionID: protocol.FnGetNodeProtocolInfo,
		Payload:    []byte{nodeID},
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

	nodeInfo = &NodeProtocolInfo{
		Capability:          ret.Payload[1],
		Security:            ret.Payload[2],
		BasicDeviceClass:    ret.Payload[4],
		GenericDeviceClass:  ret.Payload[5],
		SpecificDeviceClass: ret.Payload[6],
	}

	return
}

// NodeProtocolInfo contains protocol info for a node.
type NodeProtocolInfo struct {
	Capability          byte
	Security            byte
	BasicDeviceClass    byte
	GenericDeviceClass  byte
	SpecificDeviceClass byte
}

// IsListening returns whether a node is listening.
func (n *NodeProtocolInfo) IsListening() bool {
	return n.Capability&0x80 == 0x80
}

// GetBasicDeviceClassName will return the basic device class as a string
func (n *NodeProtocolInfo) GetBasicDeviceClassName() string {
	return protocol.GetBasicDeviceTypeName(n.BasicDeviceClass)
}

// GetGenericDeviceClassName will return the generic device class as a string
func (n *NodeProtocolInfo) GetGenericDeviceClassName() string {
	return protocol.GetGenericDeviceTypeName(n.GenericDeviceClass)
}

// GetSpecificDeviceClassName will return the specific device class as a string
func (n *NodeProtocolInfo) GetSpecificDeviceClassName() string {
	return protocol.GetSpecificDeviceTypeName(n.GenericDeviceClass, n.SpecificDeviceClass)
}
