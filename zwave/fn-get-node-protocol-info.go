package zwave

type GetNodeProtocolInfoResponse struct {
	Capability          byte
	Security            byte
	BasicDeviceClass    byte
	GenericDeviceClass  byte
	SpecificDeviceClass byte
}

func ParseGetNodeProtocolInfoResponse(payload []byte) *GetNodeProtocolInfoResponse {
	val := &GetNodeProtocolInfoResponse{
		Capability:          payload[0],
		Security:            payload[1],
		BasicDeviceClass:    payload[3],
		GenericDeviceClass:  payload[4],
		SpecificDeviceClass: payload[5],
	}

	return val
}

func (n *GetNodeProtocolInfoResponse) IsListening() bool {
	return n.Capability&0x80 == 0x80
}

func (n *GetNodeProtocolInfoResponse) GetBasicDeviceClassName() string {
	return GetBasicTypeName(n.BasicDeviceClass)
}

func (n *GetNodeProtocolInfoResponse) GetGenericDeviceClassName() string {
	return GetGenericTypeName(n.BasicDeviceClass)
}
