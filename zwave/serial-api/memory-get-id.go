package serialapi

import (
	"encoding/binary"
	"errors"

	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/session"
)

func (s *SerialAPILayer) MemoryGetId() (homeId uint32, nodeId byte, err error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionId: protocol.FnMemoryGetId,
		HasReturn:  true,
		ReturnCallback: func(err error, ret *frame.Frame) bool {
			done <- ret
			return false
		},
	}

	s.sessionLayer.MakeRequest(request)
	ret := <-done

	if ret == nil {
		return 0, 0, errors.New("Error getting home/node id")
	}

	homeId = binary.BigEndian.Uint32(ret.Payload[1:5])
	nodeId = ret.Payload[5]

	return
}
