package serialapi

import "github.com/bjyoungblood/gozw/zwave/session"

type SerialAPILayer struct {
	sessionLayer session.SessionLayer
}

func NewSerialAPILayer(sessionLayer session.SessionLayer) *SerialAPILayer {
	return &SerialAPILayer{
		sessionLayer: sessionLayer,
	}
}
