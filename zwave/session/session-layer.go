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

type ISessionLayer interface {
	MakeRequest(request *Request)
	SendFrameDirect(req *frame.Frame)
	UnsolicitedFramesChan() chan frame.Frame
}

type SessionLayer struct {
	frameLayer frame.IFrameLayer

	UnsolicitedFrames chan frame.Frame

	lastRequestFuncId byte
	responses         chan frame.Frame

	// maps sequence number to callback
	sequenceNumber byte
	callbacks      map[byte]CallbackFunc

	requestQueue chan *Request
}

func NewSessionLayer(frameLayer frame.IFrameLayer) *SessionLayer {
	session := &SessionLayer{
		frameLayer: frameLayer,

		UnsolicitedFrames: make(chan frame.Frame, 10),

		lastRequestFuncId: 0,
		responses:         make(chan frame.Frame),

		sequenceNumber: 0,
		callbacks:      map[byte]CallbackFunc{},

		requestQueue: make(chan *Request, 10),
	}

	go session.receiveThread()
	go session.sendThread()

	return session
}

func (s *SessionLayer) MakeRequest(request *Request) {
	// Enqueue the request for processing
	s.requestQueue <- request
}

// Be careful with this. Should not be called outside of a callback
func (s *SessionLayer) SendFrameDirect(req *frame.Frame) {
	s.frameLayer.Write(req)
}

func (s *SessionLayer) UnsolicitedFramesChan() chan frame.Frame {
	return s.UnsolicitedFrames
}

func (s *SessionLayer) receiveThread() {
	for frameIn := range s.frameLayer.GetOutputChannel() {
		if frameIn.IsResponse() {
			if frameIn.Payload[0] == s.lastRequestFuncId {
				select {
				case s.responses <- frameIn:
				default:
				}

				s.lastRequestFuncId = 0
			} else {
				fmt.Println("Received an unexpected response frame: ", frameIn)
			}
		} else {
			var callbackId byte

			if s.lastRequestFuncId != 0 {
				fmt.Println("REQUEST/RESPONSE COLLISION; SENDING CAN FRAME AND RETRYING PREVIOUS SEND")
				s.frameLayer.Write(frame.NewCanFrame())
				select {
				case s.responses <- *frame.NewCanFrame():
				default:
				}
				return
			}

			switch frameIn.Payload[0] {

			// These commands, when received as requests, are always callbacks and will
			// have the callback id as the first byte after the function id
			case protocol.FnAddNodeToNetwork,
				protocol.FnRemoveNodeFromNetwork,
				protocol.FnSendData,
				protocol.FnSetDefault,
				protocol.FnRequestNetworkUpdate,
				protocol.FnRemoveFailingNode:

				callbackId = frameIn.Payload[1]

				// These commands are never callbacks and shouldn't ever be handled as such
			case protocol.FnApplicationControllerUpdate,
				protocol.FnApplicationCommandHandler,
				protocol.FnApplicationCommandHandlerBridge:

				callbackId = 0

				// Log in case we need to set up a callback for a function
			default:
				fmt.Println("session-layer: got unknown callback for func: ", hex.EncodeToString([]byte{frameIn.Payload[0]}))
				callbackId = 0
			}

			if callback, ok := s.callbacks[callbackId]; ok {
				go callback(frameIn)
			} else {
				s.UnsolicitedFrames <- frameIn
			}

		}
	}
}

// This function currently assumes that every single function that expects a callback
// sets the callback id as the last byte in the payload.
func (s *SessionLayer) sendThread() {
	for request := range s.requestQueue {
		var seqNo byte = 0

		if request.ReceivesCallback {
			seqNo = s.getSequenceNumber()
			request.Payload = append(request.Payload, seqNo)
			s.callbacks[seqNo] = request.Callback
		}

		if request.Payload == nil {
			request.Payload = []byte{}
		}

		var frame = frame.NewRequestFrame(append([]byte{request.FunctionId}, request.Payload...))
		attempts := 0

	retry:
		if request.HasReturn {
			s.lastRequestFuncId = request.FunctionId
		}

		s.frameLayer.Write(frame)

		if request.HasReturn {
			select {
			case response := <-s.responses:
				if response.IsCan() {
					// Hopefully we won't collide again if we wait for 10ms :)
					time.Sleep(100 * time.Millisecond)
					if attempts > 3 {
						fmt.Println("TOO MANY RETRIES")
						request.ReturnCallback(errors.New("Too many retries sending command"), nil)
						return
					}

					attempts += 1
					goto retry // lol not even joking
				}

				if request.ReturnCallback(nil, &response) == false {
					continue
				}

			case <-time.After(10 * time.Second):
				if request.ReturnCallback(errors.New("Response timeout"), nil) == false {
					continue
				}
			}
		}

		if request.ReceivesCallback && request.Lock {
			select {
			case <-request.Release:
			case <-time.After(request.Timeout):
				fmt.Println("session lock timeout")
			}
		}

	}
}

func (s *SessionLayer) getSequenceNumber() byte {
	if s.sequenceNumber == MaxSequenceNumber {
		s.sequenceNumber = MinSequenceNumber
	} else {
		s.sequenceNumber = s.sequenceNumber + 1
	}

	return s.sequenceNumber
}
