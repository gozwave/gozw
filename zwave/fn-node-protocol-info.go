package zwave

import "github.com/bjyoungblood/gozw/zwave/protocol"

type NodeProtocolInfoResponse struct {
	CommandId           byte
	Capability          byte
	Security            byte
	BasicDeviceClass    byte
	GenericDeviceClass  byte
	SpecificDeviceClass byte
}

func ParseNodeProtocolInfoResponse(payload []byte) *NodeProtocolInfoResponse {
	val := &NodeProtocolInfoResponse{
		CommandId:           payload[0],
		Capability:          payload[1],
		Security:            payload[2],
		BasicDeviceClass:    payload[4],
		GenericDeviceClass:  payload[5],
		SpecificDeviceClass: payload[6],
	}

	return val
}

func (n *NodeProtocolInfoResponse) IsListening() bool {
	return n.Capability&0x80 == 0x80
}

func (n *NodeProtocolInfoResponse) GetBasicDeviceClassName() string {
	return protocol.GetBasicDeviceTypeName(n.BasicDeviceClass)
}

func (n *NodeProtocolInfoResponse) GetGenericDeviceClassName() string {
	return protocol.GetGenericDeviceTypeName(n.GenericDeviceClass)
}

func (n *NodeProtocolInfoResponse) GetSpecificDeviceClassName() string {
	return protocol.GetSpecificDeviceTypeName(n.GenericDeviceClass, n.SpecificDeviceClass)
}
