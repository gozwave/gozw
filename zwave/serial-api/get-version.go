package serialapi

import (
	"errors"

	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/session"
)

type Version struct {
	Version string
}

func (s *SerialAPILayer) GetVersion() (*Version, error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionId: protocol.FnGetVersion,
		Payload:    []byte{},
		HasReturn:  true,
		ReturnCallback: func(err error, ret *frame.Frame) bool {
			done <- ret
			return false
		},
	}

	s.sessionLayer.MakeRequest(request)
	versionFrame := <-done

	if versionFrame == nil {
		return nil, errors.New("Error getting version")
	}

	return &Version{Version: string(versionFrame.Payload)}, nil

}
