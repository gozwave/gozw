package serialapi

import (
	"errors"

	"github.com/helioslabs/gozw/zwave/frame"
	"github.com/helioslabs/gozw/zwave/protocol"
	"github.com/helioslabs/gozw/zwave/session"
)

func (s *SerialAPILayer) RequestNodeInfo(nodeId byte) error {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionId: protocol.FnRequestNodeInfo,
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
		return errors.New("Error requesting node information frame")
	}

	status := ret.Payload[1]

	if status == 0 {
		return errors.New("Failed putting node info request in transmit queue")
	}

	return nil
}
