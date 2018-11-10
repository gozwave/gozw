package serialapi

import (
	"context"

	"github.com/gozwave/gozw/protocol"
	"github.com/gozwave/gozw/session"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

// ILayer is an interface for the serialapi layer.
type ILayer interface {
	ControllerUpdates() chan ControllerUpdate
	ControllerCommands() chan ApplicationCommand
	AddNode() (*AddRemoveNodeCallback, error)
	RemoveNode() (*AddRemoveNodeCallback, error)
	GetCapabilities() (*Capabilities, error)
	GetVersion() (version *Version, err error)
	MemoryGetID() (homeID uint32, nodeID byte, err error)
	GetInitAppData() (*InitAppData, error)
	GetNodeProtocolInfo(nodeID byte) (nodeInfo *NodeProtocolInfo, err error)
	SendData(nodeID byte, payload []byte) (txTime uint16, err error)
	IsFailedNode(nodeID byte) (failed bool, err error)
	RemoveFailedNode(nodeID byte) (removed bool, err error)
	RequestNodeInfo(nodeInfo byte) (*NodeInfoFrame, error)
	SoftReset()
}

// Layer contains the serial api layer.
type Layer struct {
	sessionLayer        session.ILayer
	controllerUpdates   chan ControllerUpdate
	applicationCommands chan ApplicationCommand
	l                   *zap.Logger
	ctx                 context.Context
}

// NewLayer returns a new serialapi layer.
func NewLayer(ctx context.Context, sessionLayer session.ILayer, logger *zap.Logger) *Layer {

	layer := &Layer{
		sessionLayer:        sessionLayer,
		controllerUpdates:   make(chan ControllerUpdate, 10),
		applicationCommands: make(chan ApplicationCommand, 10),
		l:                   logger,
		ctx:                 ctx,
	}

	go layer.handleUnsolicitedFrames()

	return layer
}

// func (s *Layer) SetLogger(logger *log.Logger) {
// 	s.logger = logger
// }

// ControllerUpdates returns a channel for controller updates.
func (s *Layer) ControllerUpdates() chan ControllerUpdate {
	return s.controllerUpdates
}

// ControllerCommands returns a channel for application commands.
func (s *Layer) ControllerCommands() chan ApplicationCommand {
	return s.applicationCommands
}

func (s *Layer) handleUnsolicitedFrames() {
	for {
		select {
		case fr := <-s.sessionLayer.UnsolicitedFramesChan():
			switch fr.Payload[0] {
			case protocol.FnApplicationCommandHandler, protocol.FnApplicationCommandHandlerBridge:
				s.applicationCommands <- parseApplicationCommand(fr.Payload)
			case protocol.FnApplicationControllerUpdate:
				s.controllerUpdates <- parseControllerUpdate(fr.Payload)
			default:
				s.l.Warn("Unknown unsolicited frame!", zap.String("frame_info", spew.Sdump(fr)))
			}

		case <-s.ctx.Done():
			s.l.Info("closing unsolicited frames handler")
			return
		}
	}
}
