package zwave

type NodeListResponse struct {
	CommandId    byte
	Version      byte
	Capabilities byte
	Nodes        []byte
	ChipType     byte
	ChipVersion  byte
}

func isBitSet(mask []byte, nodeId uint8) bool {
	if (nodeId > 0) && (nodeId <= 232) {
		return ((mask[(nodeId-1)>>3] & (1 << ((nodeId - 1) & 0x07))) != 0)
	}

	return false
}

func ParseNodeListResponse(payload []byte) *NodeListResponse {
	val := &NodeListResponse{
		CommandId:    payload[0],
		Version:      payload[1],
		Capabilities: payload[2],
		Nodes:        payload[4:33],
		ChipType:     payload[33],
		ChipVersion:  payload[34],
	}

	return val
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
