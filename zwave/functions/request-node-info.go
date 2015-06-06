package functions

type RequestNodeInfo struct {
	FunctionId uint8
	NodeId     uint8
}

func NewRequestNodeInfo(nodeId uint8) RequestNodeInfo {
	return RequestNodeInfo{
		FunctionId: ZwRequestNodeInfo,
		NodeId:     nodeId,
	}
}

func (f *RequestNodeInfo) Marshal() []byte {
	return []byte{
		f.FunctionId,
		f.NodeId,
	}
}
