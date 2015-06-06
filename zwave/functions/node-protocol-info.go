package functions

type NodeProtocolInfo struct {
	FunctionId uint8
	NodeId     uint8
}

func NewNodeProtocolInfo(nodeId uint8) NodeProtocolInfo {
	return NodeProtocolInfo{
		FunctionId: ZwRequestNodeInfo,
		NodeId:     nodeId,
	}
}

func (f *NodeProtocolInfo) Marshal() []byte {
	return []byte{
		f.FunctionId,
		f.NodeId,
	}
}
