package functions

type AddNode struct {
	FunctionId uint8
	Options    uint8
	CallbackId uint8
}

func NewAddNode() AddNode {
	return AddNode{
		FunctionId: ZwAddNodeToNetwork,
		Options:    AddNodeAny | AddNodeOptionNormalPower | AddNodeOptionNetworkWide,
		CallbackId: 0x01,
	}
}

func NewAddNodeEnd() AddNode {
	return AddNode{
		FunctionId: ZwAddNodeToNetwork,
		Options:    AddNodeStop,
		CallbackId: 0x01,
	}
}

func (f *AddNode) Marshal() []byte {
	return []byte{
		f.FunctionId,
		f.Options,
		f.CallbackId,
	}
}
