package serialapi

import (
	"fmt"

	"github.com/bjyoungblood/gozw/zwave/session"
	"github.com/davecgh/go-spew/spew"
)

type SerialAPILayer struct {
	sessionLayer session.SessionLayer
}

func NewSerialAPILayer(sessionLayer session.SessionLayer) *SerialAPILayer {
	layer := &SerialAPILayer{
		sessionLayer: sessionLayer,
	}

	go func() {
		for fr := range sessionLayer.UnsolicitedFramesChan() {
			fmt.Println("unsolicited:")
			spew.Dump(fr.Payload)
		}
	}()

	return layer
}
