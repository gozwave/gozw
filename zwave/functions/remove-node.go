package functions

type RemoveNode struct {
	FunctionId uint8
	Options    uint8
	CallbackId uint8
}

func NewRemoveNode() RemoveNode {
	return RemoveNode{
		FunctionId: ZwAddNodeToNetwork,
		Options:    RemoveNodeAny | RemoveNodeOptionNormalPower | RemoveNodeOptionNetworkWide,
		CallbackId: 0x01,
	}
}

func NewRemoveNodeEnd() RemoveNode {
	return RemoveNode{
		FunctionId: ZwRemoveNodeFromNetwork,
		Options:    RemoveNodeStop,
		CallbackId: 0x01,
	}
}

func (f *RemoveNode) Marshal() []byte {
	return []byte{
		f.FunctionId,
		f.Options,
		f.CallbackId,
	}
}
