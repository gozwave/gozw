package zwave

import "fmt"

type ApplicationControllerUpdate struct {
	CommandId      byte
	Status         byte
	NodeId         byte
	Length         byte
	Basic          uint8
	Generic        uint8
	Specific       uint8
	CommandClasses []byte
}

func ParseApplicationControllerUpdate(payload []byte) *AddRemoveNodeCallback {
	val := &AddRemoveNodeCallback{
		CommandId: payload[0],
		Status:    payload[1],
		Source:    payload[2],
		Length:    payload[3],
	}

	fmt.Println(payload)

	if val.Length == 0 {
		return val
	}

	if val.Length >= 1 {
		val.Basic = payload[4]
	}

	if val.Length >= 2 {
		val.Generic = payload[5]
	}

	if val.Length >= 3 {
		val.Specific = payload[6]
	}

	if val.Length >= 4 {
		val.CommandClasses = payload[7:]
	}

	return val
}

func (a *ApplicationControllerUpdate) GetStatusString() string {
	switch a.Status {
	case UpdateStateNodeInfoReceived:
		return "Node Info Received"
	case UpdateStateNodeInfoReqDone:
		return "Node Info Req Done"
	case UpdateStateNodeInfoReqFailed:
		return "Node Info Req Failed"
	case UpdateStateRoutingPending:
		return "Routing Pending"
	case UpdateStateNewIdAssigned:
		return "New IDA ssigned"
	case UpdateStateDeleteDone:
		return "Delete Done"
	case UpdateStateSucId:
		return "Update SUC ID"
	default:
		return "Unknown"
	}
}
