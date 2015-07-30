package session

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/protocol"
)

const (
	MinSequenceNumber = 1
	MaxSequenceNumber = 127
)

type Callback func(frame.Frame)

type SessionLayer interface {
	MakeRequest(request *Request)
}

type ZWaveSessionLayer struct {
	frameLayer frame.FrameLayer

	UnsolicitedFrames chan frame.Frame

	lastRequestFuncId uint8
	responses         chan frame.Frame

	// maps sequence number to callback
	sequenceNumber uint8
	callbacks      map[uint8]Callback

	requestQueue chan *Request
}

func NewSessionLayer(frameLayer frame.FrameLayer) *ZWaveSessionLayer {
	session := &ZWaveSessionLayer{
		frameLayer: frameLayer,

		UnsolicitedFrames: make(chan frame.Frame, 10),

		lastRequestFuncId: 0,
		responses:         make(chan frame.Frame),

		sequenceNumber: 0,
		callbacks:      map[uint8]Callback{},

		requestQueue: make(chan *Request, 10),
	}

	go session.receiveThread()
	go session.sendThread()

	return session
}

func (s *ZWaveSessionLayer) MakeRequest(request *Request) {
	// Enqueue the request for processing
	s.requestQueue <- request
}

func (s *ZWaveSessionLayer) receiveThread() {
	for frame := range s.frameLayer.GetOutputChannel() {
		if frame.IsResponse() {
			if frame.Payload[0] == s.lastRequestFuncId {
				select {
				case s.responses <- frame:
				default:
				}

				s.lastRequestFuncId = 0
			} else {
				fmt.Println("Received an unexpected response frame: ", frame)
			}
		} else {
			var callbackId uint8

			switch frame.Payload[0] {

			// These commands, when received as requests, are always callbacks and will
			// have the callback id as the first byte after the function id
			case protocol.FnAddNodeToNetwork,
				protocol.FnRemoveNodeFromNetwork,
				protocol.FnSendData,
				protocol.FnSetDefault,
				protocol.FnRequestNetworkUpdate,
				protocol.FnRemoveFailingNode:

				callbackId = frame.Payload[1]

				// These commands are never callbacks and shouldn't ever be handled as such
			case protocol.FnApplicationControllerUpdate,
				protocol.FnApplicationCommandHandler,
				protocol.FnApplicationCommandHandlerBridge:

				callbackId = 0

				// Log in case we need to set up a callback for a function
			default:
				fmt.Println("session-layer: got unknown callback for func: ", hex.EncodeToString([]byte{frame.Payload[0]}))
				callbackId = 0
			}

			if callback, ok := s.callbacks[callbackId]; ok {
				go callback(frame)
			} else {
				s.UnsolicitedFrames <- frame
			}

		}
	}
}

// This function currently assumes that every single function that expects a callback
// sets the callback id as the last byte in the payload.
func (s *ZWaveSessionLayer) sendThread() {
	for request := range s.requestQueue {
		var seqNo uint8 = 0

		if request.ReceivesCallback {
			seqNo = s.getSequenceNumber()
			request.Payload = append(request.Payload, seqNo)
		}

		done := make(chan bool)
		if request.ReceivesCallback && request.Exclusive {
			oldCb := request.Callback
			request.Callback = func(cbFrame frame.Frame) bool {
				if oldCb(cbFrame) == false {
					done <- true
					return false
				}

				return true
			}
		}

		if request.Payload == nil {
			request.Payload = []byte{}
		}

		var frame = frame.NewRequestFrame(append([]byte{request.FunctionId}, request.Payload...))

		s.lastRequestFuncId = request.FunctionId
		s.frameLayer.Write(frame)

		if request.HasReturn {
			select {
			case response := <-s.responses:
				if request.ReturnCallback(nil, &response) == false {
					continue
				}

			case <-time.After(10 * time.Second):
				if request.ReturnCallback(errors.New("Response timeout"), nil) == false {
					continue
				}
			}
		}

		if request.ReceivesCallback && request.Exclusive {
			select {
			case <-done:
			case <-time.After(request.Timeout):
			}
		}

	}
}

func (s *ZWaveSessionLayer) getSequenceNumber() uint8 {
	if s.sequenceNumber == MaxSequenceNumber {
		s.sequenceNumber = MinSequenceNumber
	} else {
		s.sequenceNumber = s.sequenceNumber + 1
	}

	return s.sequenceNumber
}
