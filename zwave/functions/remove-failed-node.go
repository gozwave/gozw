package functions

type RemoveFailedNode struct {
	FunctionId uint8
	NodeId     uint8
	CallbackId uint8
}

func NewRemoveFailedNode(nodeId uint8) RemoveFailedNode {
	return RemoveFailedNode{
		FunctionId: ZwRemoveFailingNode,
		NodeId:     nodeId,
		CallbackId: 0x01,
	}
}

func (f *RemoveFailedNode) Marshal() []byte {
	return []byte{
		f.FunctionId,
		f.NodeId,
		f.CallbackId,
	}
}
