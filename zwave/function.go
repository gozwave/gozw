package zwave

type GenericPayload struct {
	CommandId byte
	Payload   []byte
}

func (p GenericPayload) Marshal() []byte {
	return append([]byte{p.CommandId}, p.Payload...)
}

func ParseFunctionPayload(payload []byte) interface{} {

	switch payload[0] {
	case FnGetInitAppData:
		return ParseNodeListResponse(payload)
	case FnSerialApiCapabilities:
		return ParseSerialApiCapabilitiesResponse(payload)
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
