package zwave

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

func (n *NodeProtocolInfoResponse) GetSpecificDeviceClassName() string {
	return GetSpecificTypeName(n.BasicDeviceClass, n.SpecificDeviceClass)
}
