package serialapi

import (
	"errors"
	"time"

	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/protocol"
	"github.com/gozwave/gozw/session"
)

// RemoveFailedNode will remove a failed node.
func (s *Layer) RemoveFailedNode(nodeID byte) (removed bool, err error) {

	done := make(chan frame.Frame)

	request := &session.Request{
		FunctionID:       protocol.FnRemoveFailingNode,
		Payload:          []byte{nodeID},
		HasReturn:        true,
		ReceivesCallback: true,
		Timeout:          time.Second * 10,

		ReturnCallback: func(err error, ret *frame.Frame) bool {
			return true
		},

		Callback: func(cbFrame frame.Frame) {
			done <- cbFrame
		},
	}

	s.sessionLayer.MakeRequest(request)

	result := <-done

	switch result.Payload[2] {
	case protocol.NodeOk:
		return false, errors.New("Node was not failing")
	case protocol.FailedNodeRemoved:
		return true, nil
	case protocol.FailedNodeNotRemoved:
		return false, errors.New("Error removing node")
	}

	return false, errors.New("Unknown status")
}
