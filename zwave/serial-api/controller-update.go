package serialapi

type ControllerUpdate struct {
	Status         byte
	NodeId         byte
	Length         byte
	Basic          uint8
	Generic        uint8
	Specific       uint8
	CommandClasses []byte
}

const (
	UpdateStateNodeInfoReceived  = 0x84
	UpdateStateNodeInfoReqDone   = 0x82
	UpdateStateNodeInfoReqFailed = 0x81
	UpdateStateRoutingPending    = 0x80
	UpdateStateNewIdAssigned     = 0x40
	UpdateStateDeleteDone        = 0x20
	UpdateStateSucId             = 0x10
)

func parseControllerUpdate(payload []byte) ControllerUpdate {
	val := ControllerUpdate{
		Status: payload[1],
		NodeId: payload[2],
		Length: payload[3],
	}

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

func (a *ControllerUpdate) GetStatusString() string {
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
		return "New ID Assigned"
	case UpdateStateDeleteDone:
		return "Delete Done"
	case UpdateStateSucId:
		return "Update SUC ID"
	default:
		return "Unknown"
	}
}
