package zwave

type NodeList struct {
	Version      uint8
	Capabilities byte
	Nodes        []byte
	ChipType     byte
	ChipVersion  byte
}

// @todo functions to parse capabilities flags

func isBitSet(mask byte, pos uint) bool {
	return mask&(1<<(7-pos)) > 0
}

func (n *NodeList) Unmarshal(frame *Frame) {
	n.Version = frame.Payload[1]
	n.Capabilities = frame.Payload[2]
	n.Nodes = frame.Payload[4:33]
	n.ChipType = frame.Payload[33]
	n.ChipVersion = frame.Payload[34]
}

func (n *NodeList) GetNodeIds() []int {
	nodes := []int{}
	nodeNum := 1

	for i := 0; i < 29; i++ {
		for j := uint(0); j < 8; j++ {
			if isBitSet(n.Nodes[i], j) {
				nodes = append(nodes, nodeNum)
			}

			nodeNum++
		}
	}

	return nodes
}
