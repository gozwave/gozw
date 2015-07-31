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
	AddNode() (*AddRemoveNodeCallback, error)
	RemoveNode() (*AddRemoveNodeCallback, error)
	GetSerialApiCapabilities() (*SerialApiCapabilities, error)
	GetVersion() (version *Version, err error)
	MemoryGetId() (homeId uint32, nodeId uint8, err error)
	GetInitAppData() (*InitAppData, error)
	GetNodeProtocolInfo(nodeId uint8) (nodeInfo *NodeProtocolInfo, err error)
	SendData(nodeId byte, payload []byte) (txTime uint16, err error)
	IsFailedNode(nodeId byte) (failed bool, err error)
	RemoveFailedNode(nodeId byte) (removed bool, err error)
	RequestNodeInfo(nodeInfo byte) (err error)
	SoftReset()
}

type SerialAPILayer struct {
	sessionLayer        session.ISessionLayer
	controllerUpdates   chan ControllerUpdate
	applicationCommands chan ApplicationCommand
}

func NewSerialAPILayer(sessionLayer session.ISessionLayer) *SerialAPILayer {
	layer := &SerialAPILayer{
		sessionLayer:        sessionLayer,
		controllerUpdates:   make(chan ControllerUpdate, 10),
		applicationCommands: make(chan ApplicationCommand, 10),
	}

	go layer.handleUnsolicitedFrames()

	return layer
}

func (s *SerialAPILayer) ControllerUpdates() chan ControllerUpdate {
	return s.controllerUpdates
}

func (s *SerialAPILayer) ControllerCommands() chan ApplicationCommand {
	return s.applicationCommands
}

func (s *SerialAPILayer) handleUnsolicitedFrames() {
	for fr := range s.sessionLayer.UnsolicitedFramesChan() {
		switch fr.Payload[0] {
		case protocol.FnApplicationCommandHandler, protocol.FnApplicationCommandHandlerBridge:
			s.applicationCommands <- parseApplicationCommand(fr.Payload)
		case protocol.FnApplicationControllerUpdate:
			s.controllerUpdates <- parseControllerUpdate(fr.Payload)
		default:
			fmt.Println("Unknown unsolicited frame!")
			spew.Dump(fr)
		}
	}
}
