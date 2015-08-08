package serialapi

import (
	"errors"
	"time"

	"github.com/helioslabs/gozw/zwave/frame"
	"github.com/helioslabs/gozw/zwave/protocol"
	"github.com/helioslabs/gozw/zwave/session"
)

func (s *SerialAPILayer) RemoveFailedNode(nodeId byte) (removed bool, err error) {

	done := make(chan frame.Frame)

	request := &session.Request{
		FunctionId:       protocol.FnRemoveFailingNode,
		Payload:          []byte{nodeId},
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
