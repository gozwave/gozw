package serialapi

import (
	"errors"

	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/protocol"
	"github.com/gozwave/gozw/session"
)

// IsFailedNode Will return if a node has failed.
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
