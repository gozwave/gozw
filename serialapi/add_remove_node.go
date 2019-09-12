package serialapi

import (
	"errors"
	"fmt"
	"time"

	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/protocol"
	"github.com/gozwave/gozw/session"
	"go.uber.org/zap"
)

// AddNode will put the controller into add node mode and handle operations for adding a node.
func (s *Layer) AddNode() (*AddRemoveNodeCallback, error) {

	var newNode *AddRemoveNodeCallback

	addNodeDone := make(chan bool)
	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionID: protocol.FnAddNodeToNetwork,
		Payload:    []byte{protocol.AddNodeAny | protocol.AddNodeOptionNetworkWide | protocol.AddNodeOptionNormalPower},

		HasReturn:        false,
		ReceivesCallback: true,
		Lock:             true,
		Timeout:          60 * time.Second,
		Release:          addNodeDone,

		Callback: func(cbFrame frame.Frame) {
			cbData := parseAddRemoveNodeCallback(cbFrame.Payload)

			switch cbData.Status {
			case protocol.AddNodeStatusLearnReady:
				s.l.Debug("ADD NODE: learn ready")

			case protocol.AddNodeStatusNodeFound:
				s.l.Debug("ADD NODE: node found")

			case protocol.AddNodeStatusAddingSlave:
				s.l.Debug("ADD NODE: adding slave node")
				newNode = cbData

			case protocol.AddNodeStatusAddingController:
				// hey, i just met you, and this is crazy
				// but it could happen, so implement me maybe
				s.l.Debug("ADD NODE: adding controller node")
				newNode = cbData

			case protocol.AddNodeStatusProtocolDone:
				s.l.Debug("ADD NODE: protocol done")
				reply := addRemoveStatusFrame(
					protocol.FnAddNodeToNetwork,
					protocol.AddNodeStop,
					cbData.CallbackID,
				)
				s.sessionLayer.SendFrameDirect(reply)

			case protocol.AddNodeStatusDone:
				s.l.Debug("ADD NODE: done")
				reply := addRemoveStatusFrame(
					protocol.FnAddNodeToNetwork,
					protocol.AddNodeStop,
					0,
				)
				s.sessionLayer.SendFrameDirect(reply)

				addNodeDone <- true
				close(addNodeDone)
				done <- &cbFrame

			case protocol.AddNodeStatusFailed:
				s.l.Error("ADD NODE: failed")
				reply := addRemoveStatusFrame(
					protocol.FnAddNodeToNetwork,
					protocol.AddNodeStop,
					cbData.CallbackID,
				)
				s.sessionLayer.SendFrameDirect(reply)

			default:
				s.l.Warn("ADD NODE: unknown status", zap.String("status", fmt.Sprint(cbData.Status)))
			}
		},
	}

	s.sessionLayer.MakeRequest(request)
	ret := <-done

	if ret == nil {
		return nil, errors.New("Error adding node")
	}

	return newNode, nil

}

// RemoveNode will put the controller into remove node mode  and handle all operations.
func (s *Layer) RemoveNode() (*AddRemoveNodeCallback, error) {

	var removedNode *AddRemoveNodeCallback

	removeNodeDone := make(chan bool)
	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionID: protocol.FnRemoveNodeFromNetwork,
		Payload:    []byte{protocol.RemoveNodeAny | protocol.RemoveNodeOptionNetworkWide | protocol.RemoveNodeOptionNormalPower},

		HasReturn:        false,
		ReceivesCallback: true,
		Lock:             true,
		Timeout:          60 * time.Second,
		Release:          removeNodeDone,

		Callback: func(cbFrame frame.Frame) {
			cbData := parseAddRemoveNodeCallback(cbFrame.Payload)

			switch cbData.Status {
			case protocol.RemoveNodeStatusLearnReady:
				s.l.Debug("REMOVE NODE: learn ready")

			case protocol.RemoveNodeStatusNodeFound:
				s.l.Debug("REMOVE NODE: node found")

			case protocol.RemoveNodeStatusRemovingSlave:
				s.l.Debug("REMOVE NODE: removing slave node")
				removedNode = cbData

			case protocol.RemoveNodeStatusRemovingController:
				// hey, i just met you, and this is crazy
				// but it could happen, so implement me maybe
				s.l.Debug("REMOVE NODE: removing controller node")
				removedNode = cbData

			case protocol.RemoveNodeStatusProtocolDone:
				s.l.Debug("REMOVE NODE: protocol done")
				reply := addRemoveStatusFrame(
					protocol.FnRemoveNodeFromNetwork,
					protocol.RemoveNodeStop,
					cbData.CallbackID,
				)
				s.sessionLayer.SendFrameDirect(reply)

			case protocol.RemoveNodeStatusDone:
				s.l.Debug("REMOVE NODE: done")
				reply := addRemoveStatusFrame(
					protocol.FnRemoveNodeFromNetwork,
					protocol.RemoveNodeStop,
					0,
				)
				s.sessionLayer.SendFrameDirect(reply)

				removeNodeDone <- true
				close(removeNodeDone)
				done <- &cbFrame

			case protocol.RemoveNodeStatusFailed:
				s.l.Error("REMOVE NODE: failed")
				reply := addRemoveStatusFrame(
					protocol.FnRemoveNodeFromNetwork,
					protocol.RemoveNodeStop,
					cbData.CallbackID,
				)
				s.sessionLayer.SendFrameDirect(reply)

			default:
				s.l.Warn("REMOVE NODE: unknown status", zap.String("status", fmt.Sprint(cbData.Status)))
			}
		},
	}

	s.sessionLayer.MakeRequest(request)
	ret := <-done

	if ret == nil {
		return nil, errors.New("Error removing node")
	}

	return removedNode, nil

}

func addRemoveStatusFrame(functionID, status, callbackID byte) *frame.Frame {
	return frame.NewRequestFrame([]byte{
		functionID,
		status,
		callbackID,
	})
}

// AddRemoveNodeCallback contains callback info for add/remove node.
type AddRemoveNodeCallback struct {
	CommandID      byte
	CallbackID     byte
	Status         byte
	Source         byte
	Length         byte
	Basic          byte
	Generic        byte
	Specific       byte
	CommandClasses []byte
}

func parseAddRemoveNodeCallback(payload []byte) *AddRemoveNodeCallback {
	val := &AddRemoveNodeCallback{
		CommandID:  payload[0],
		CallbackID: payload[1],
		Status:     payload[2],
		Source:     payload[3],
		Length:     payload[4],
	}

	if val.Length == 0 {
		return val
	}

	if val.Length >= 1 {
		val.Basic = payload[5]
	}

	if val.Length >= 2 {
		val.Generic = payload[6]
	}

	if val.Length >= 3 {
		val.Specific = payload[7]
	}

	if val.Length >= 4 {
		val.CommandClasses = payload[8:]
	}

	return val
}
