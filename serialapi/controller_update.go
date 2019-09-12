package serialapi

import "github.com/gozwave/gozw/protocol"

// ControllerUpdate .
type ControllerUpdate struct {
	Status         byte
	NodeID         byte
	Length         byte
	Basic          byte
	Generic        byte
	Specific       byte
	CommandClasses []byte
}

func parseControllerUpdate(payload []byte) ControllerUpdate {
	val := ControllerUpdate{
		Status: payload[1],
		NodeID: payload[2],
		Length: payload[3],
	}

	if val.Length == 0 || val.Status == protocol.UpdateStateNodeInfoReqFailed {
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

// GetStatusString will get the  controller status as a string
func (a *ControllerUpdate) GetStatusString() string {
	switch a.Status {
	case protocol.UpdateStateNodeInfoReceived:
		return "Node Info Received"
	case protocol.UpdateStateNodeInfoReqDone:
		return "Node Info Req Done"
	case protocol.UpdateStateNodeInfoReqFailed:
		return "Node Info Req Failed"
	case protocol.UpdateStateRoutingPending:
		return "Routing Pending"
	case protocol.UpdateStateNewIDAssigned:
		return "New ID Assigned"
	case protocol.UpdateStateDeleteDone:
		return "Delete Done"
	case protocol.UpdateStateSucID:
		return "Update SUC ID"
	default:
		return "Unknown"
	}
}
