package functions

type SendData struct {
	FunctionId uint8
	NodeId     uint8
	// @todo make command class struct/interface
	Payload []byte
}

func NewSendData(nodeId uint8) SendData {
	return SendData{
		FunctionId: ZwSendData,
		NodeId:     nodeId,
	}
}

func (f *SendData) Marshal() []byte {
	buf := []byte{
		f.FunctionId,
		byte(f.NodeId),
		byte(len(f.Payload)),
	}

	buf = append(buf, f.Payload...)
	// @todo transport options
	buf = append(buf, 0x1, 0x1)

	return buf
}
