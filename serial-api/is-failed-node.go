package serialapi

import (
	"errors"

	"gitlab.com/helioslabs/gozw/frame"
	"gitlab.com/helioslabs/gozw/protocol"
	"gitlab.com/helioslabs/gozw/session"
)

func (s *Layer) IsFailedNode(nodeID byte) (failed bool, err error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionID: protocol.FnIsNodeFailed,
		Payload:    []byte{nodeID},
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
