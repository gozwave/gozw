package serialapi

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/helioslabs/gozw/zwave/protocol"
	"github.com/helioslabs/gozw/zwave/session"
)

type ILayer interface {
	ControllerUpdates() chan ControllerUpdate
	ControllerCommands() chan ApplicationCommand
	AddNode() (*AddRemoveNodeCallback, error)
	RemoveNode() (*AddRemoveNodeCallback, error)
	GetSerialAPICapabilities() (*SerialAPICapabilities, error)
	GetVersion() (version *Version, err error)
	MemoryGetID() (homeID uint32, nodeID byte, err error)
	GetInitAppData() (*InitAppData, error)
	GetNodeProtocolInfo(nodeID byte) (nodeInfo *NodeProtocolInfo, err error)
	SendData(nodeID byte, payload []byte) (txTime uint16, err error)
	IsFailedNode(nodeID byte) (failed bool, err error)
	RemoveFailedNode(nodeID byte) (removed bool, err error)
	RequestNodeInfo(nodeInfo byte) (err error)
	SoftReset()
}

type Layer struct {
	sessionLayer        session.ILayer
	controllerUpdates   chan ControllerUpdate
	applicationCommands chan ApplicationCommand
}

func NewLayer(sessionLayer session.ILayer) *Layer {
	layer := &Layer{
		sessionLayer:        sessionLayer,
		controllerUpdates:   make(chan ControllerUpdate, 10),
		applicationCommands: make(chan ApplicationCommand, 10),
	}

	go layer.handleUnsolicitedFrames()

	return layer
}

func (s *Layer) ControllerUpdates() chan ControllerUpdate {
	return s.controllerUpdates
}

func (s *Layer) ControllerCommands() chan ApplicationCommand {
	return s.applicationCommands
}

func (s *Layer) handleUnsolicitedFrames() {
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
