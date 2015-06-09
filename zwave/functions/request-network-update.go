package functions

type RequestNetworkUpdate struct {
	FunctionId uint8
	CallbackId uint8
}

func NewRequestNetworkUpdate() RequestNetworkUpdate {
	return RequestNetworkUpdate{
		FunctionId: ZwRequestNetworkUpdate,
		CallbackId: 0x0,
	}
}

func (f *RequestNetworkUpdate) Marshal() []byte {
	return []byte{
		f.FunctionId,
		f.CallbackId,
	}
}
