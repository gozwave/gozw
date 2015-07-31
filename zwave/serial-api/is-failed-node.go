package serialapi

import (
	"errors"

	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/session"
)

func (s *SerialAPILayer) IsFailedNode(nodeId byte) (failed bool, err error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionId: protocol.FnIsNodeFailed,
		Payload:    []byte{nodeId},
		HasReturn:  true,
		ReturnCallback: func(err error, ret *frame.Frame) bool {
			done <- ret
			return false
		},
	}

	s.sessionLayer.MakeRequest(request)
	ret := <-done

	if ret == nil {
		err = errors.New("Error checking failure status")
		return
	}

	failed = ret.Payload[1] == 1

	return
}
