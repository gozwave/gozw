package functions

type IsNodeFailed struct {
	FunctionId uint8
	NodeId     uint8
}

func NewIsNodeFailed(nodeId uint8) IsNodeFailed {
	return IsNodeFailed{
		FunctionId: ZwIsNodeFailed,
		NodeId:     nodeId,
	}
}

func (f *IsNodeFailed) Marshal() []byte {
	return []byte{
		f.FunctionId,
		f.NodeId,
	}
}
