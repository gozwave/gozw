package zwave

type NodeProtocolInfoResponse struct {
	Capability          byte
	Security            byte
	BasicDeviceClass    byte
	GenericDeviceClass  byte
	SpecificDeviceClass byte
}

func ParseNodeProtocolInfoResponse(payload []byte) *NodeProtocolInfoResponse {
	val := &NodeProtocolInfoResponse{
		Capability:          payload[0],
		Security:            payload[1],
		BasicDeviceClass:    payload[3],
		GenericDeviceClass:  payload[4],
		SpecificDeviceClass: payload[5],
	}

	return val
}

func (n *NodeProtocolInfoResponse) IsListening() bool {
	return n.Capability&0x80 == 0x80
}

func (n *NodeProtocolInfoResponse) GetBasicDeviceClassName() string {
	return GetBasicTypeName(n.BasicDeviceClass)
}

func (n *NodeProtocolInfoResponse) GetGenericDeviceClassName() string {
	return GetGenericTypeName(n.BasicDeviceClass)
}
