package zwave

import "encoding/binary"

type MemoryGetIdResponse struct {
	CommandId byte
	HomeId    uint32
	NodeId    byte
}

func ParseMemoryGetIdResponse(payload []byte) *MemoryGetIdResponse {
	val := &MemoryGetIdResponse{
		CommandId: payload[0],
		HomeId:    binary.BigEndian.Uint32(payload[1:5]),
		NodeId:    payload[5],
	}

	return val
}
