package zwave

import "fmt"

const (
	FnGetInitAppData              = 0x02
	FnApplicationNodeInformation  = 0x03
	FnSerialApiCapabilities       = 0x07
	FnSendData                    = 0x13
	FnGetVersion                  = 0x15
	FnMemoryGetId                 = 0x20
	FnGetNodeProtocolInfo         = 0x41
	FnSetDefault                  = 0x42
	FnApplicationControllerUpdate = 0x49
	FnAddNodeToNetwork            = 0x4a
	FnRemoveNodeFromNetwork       = 0x4b
	FnRequestNetworkUpdate        = 0x53
	FnRequestNodeInfo             = 0x60
	FnRemoveFailingNode           = 0x61
	FnIsNodeFailed                = 0x62
	FnSerialAPIReady              = 0xEF
)

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
	case FnGetNodeProtocolInfo:
		return ParseNodeProtocolInfoResponse(payload)
	case FnGetVersion:
		return ParseVersionResponse(payload)
	case FnMemoryGetId:
		return ParseMemoryGetIdResponse(payload)
	case FnAddNodeToNetwork, FnRemoveNodeFromNetwork:
		return ParseAddNodeCallback(payload)
	case FnApplicationControllerUpdate:
		return ParseApplicationControllerUpdate(payload)
	default:
		fmt.Println("UNKNOWN:", payload[0])
		val := &GenericPayload{
			CommandId: payload[0],
		}

		if len(payload) > 1 {
			val.Payload = payload[1:]
		}

		return val
	}

}
