package serialapi

import (
	"fmt"

	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/session"
	"github.com/davecgh/go-spew/spew"
)

type ISerialAPILayer interface {
	ControllerUpdates() chan ControllerUpdate
	ControllerCommands() chan ApplicationCommand
}

type SerialAPILayer struct {
	controllerUpdates  chan ControllerUpdate
	controllerCommands chan ApplicationCommand
	sessionLayer        session.ISessionLayer
}

func NewSerialAPILayer(sessionLayer session.ISessionLayer) *SerialAPILayer {
	layer := &SerialAPILayer{
		sessionLayer:       sessionLayer,
		controllerUpdates:  make(chan ControllerUpdate, 10),
		controllerCommands: make(chan ApplicationCommand, 10),
	}

	go layer.handleUnsolicitedFrames()

	return layer
}

func (s *SerialAPILayer) ControllerUpdates() chan ControllerUpdate {
	return s.controllerUpdates
}

func (s *SerialAPILayer) ControllerCommands() chan ApplicationCommand {
	return s.controllerCommands
}

func (s *SerialAPILayer) handleUnsolicitedFrames() {
	for fr := range s.sessionLayer.UnsolicitedFramesChan() {
		switch fr.Payload[0] {
		case protocol.FnApplicationCommandHandler, protocol.FnApplicationCommandHandlerBridge:
			s.controllerCommands <- parseApplicationCommand(fr.Payload)
		case protocol.FnApplicationControllerUpdate:
			s.controllerUpdates <- parseControllerUpdate(fr.Payload)
		default:
			fmt.Println("Unknown unsolicited frame!")
			spew.Dump(fr)
		}
	}
}
