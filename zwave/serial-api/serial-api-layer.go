package serialapi

import (
	"errors"

	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/session"
)

type SerialAPILayer struct {
	sessionLayer session.SessionLayer
}

func NewSerialAPILayer(sessionLayer session.SessionLayer) *SerialAPILayer {
	return &SerialAPILayer{
		sessionLayer: sessionLayer,
	}
}

type Version struct {
	Version string
}

func (a *SerialAPILayer) GetVersion() (*Version, error) {

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

	a.sessionLayer.MakeRequest(request)
	versionFrame := <-done

	if versionFrame == nil {
		return nil, errors.New("Error getting version")
	}

	return &Version{Version: string(versionFrame.Payload)}, nil
}
