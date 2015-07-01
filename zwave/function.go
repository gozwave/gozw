package zwave

type FunctionPayload interface {
	Marshal() []byte
}

type GenericPayload struct {
	CommandId byte
	Payload   []byte
}

func (p GenericPayload) Marshal() []byte {
	return append([]byte{p.CommandId}, p.Payload...)
}

func MarshalPayload(payload FunctionPayload) []byte {
	return payload.Marshal()
}

func ParseFunctionPayload(payload []byte) FunctionPayload {

	switch payload[0] {
	case 0x02:
		return ParseNodeListResponse(payload)
	default:
		val := &GenericPayload{
			CommandId: payload[0],
		}

		if len(payload) > 1 {
			val.Payload = payload[1:]
		}

		return val
	}

}
