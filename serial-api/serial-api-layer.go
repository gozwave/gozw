package serialapi

import (
	"log"
	"os"

	"github.com/comail/colog"
	"github.com/davecgh/go-spew/spew"
	"gitlab.com/helioslabs/gozw/protocol"
	"gitlab.com/helioslabs/gozw/session"
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
	logger              *log.Logger
}

func NewLayer(sessionLayer session.ILayer) *Layer {
	serialApiLogger := colog.NewCoLog(os.Stdout, "serial-api ", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	serialApiLogger.ParseFields(true)

	layer := &Layer{
		sessionLayer:        sessionLayer,
		controllerUpdates:   make(chan ControllerUpdate, 10),
		applicationCommands: make(chan ApplicationCommand, 10),
		logger:              serialApiLogger.NewLogger(),
	}

	go layer.handleUnsolicitedFrames()

	return layer
}

func (s *Layer) SetLogger(logger *log.Logger) {
	s.logger = logger
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
			s.logger.Println("warn: Unknown unsolicited frame!", spew.Sdump(fr))
		}
	}
}
