package session

import (
	"encoding/hex"
	"errors"
	"log"
	"os"
	"time"

	"github.com/comail/colog"
	"github.com/helioslabs/gozw/frame"
	"github.com/helioslabs/gozw/protocol"
)

const (
	minSequenceNumber = 1
	maxSequenceNumber = 127
)

type ILayer interface {
	MakeRequest(request *Request)
	SendFrameDirect(req *frame.Frame)
	UnsolicitedFramesChan() chan frame.Frame
}

type Layer struct {
	frameLayer frame.ILayer

	UnsolicitedFrames chan frame.Frame

	lastRequestFuncID byte
	responses         chan frame.Frame

	// maps sequence number to callback
	sequenceNumber byte
	callbacks      map[byte]CallbackFunc

	logger *log.Logger

	requestQueue chan *Request
}

func NewSessionLayer(frameLayer frame.ILayer) *Layer {
	sessionLogger := colog.NewCoLog(os.Stdout, "session ", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	sessionLogger.ParseFields(true)

	session := &Layer{
		frameLayer: frameLayer,

		UnsolicitedFrames: make(chan frame.Frame, 10),

		lastRequestFuncID: 0,
		responses:         make(chan frame.Frame),

		sequenceNumber: 0,
		callbacks:      map[byte]CallbackFunc{},

		logger: sessionLogger.NewLogger(),

		requestQueue: make(chan *Request, 10),
	}

	go session.receiveThread()
	go session.sendThread()

	return session
}

func (s *Layer) MakeRequest(request *Request) {
	// Enqueue the request for processing
	s.requestQueue <- request
}

// Be careful with this. Should not be called outside of a callback
func (s *Layer) SendFrameDirect(req *frame.Frame) {
	s.frameLayer.Write(req)
}

func (s *Layer) UnsolicitedFramesChan() chan frame.Frame {
	return s.UnsolicitedFrames
}

func (s *Layer) receiveThread() {
	for frameIn := range s.frameLayer.GetOutputChannel() {
		if frameIn.IsResponse() {
			if frameIn.Payload[0] == s.lastRequestFuncID {
				select {
				case s.responses <- frameIn:
				default:
				}

				s.lastRequestFuncID = 0
			} else {
				s.logger.Println("warn: received an unexpected response frame: ", frameIn)
			}
		} else {
			var callbackID byte

			if s.lastRequestFuncID != 0 {
				s.logger.Println("REQUEST/RESPONSE COLLISION; SENDING CAN FRAME AND RETRYING PREVIOUS SEND")
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

				callbackID = frameIn.Payload[1]

				// These commands are never callbacks and shouldn't ever be handled as such
			case protocol.FnApplicationControllerUpdate,
				protocol.FnApplicationCommandHandler,
				protocol.FnApplicationCommandHandlerBridge:

				callbackID = 0

				// Log in case we need to set up a callback for a function
			default:
				s.logger.Println("warn: got unknown callback for func: ", hex.EncodeToString([]byte{frameIn.Payload[0]}))
				callbackID = 0
			}

			if callback, ok := s.callbacks[callbackID]; ok {
				go callback(frameIn)
			} else {
				s.UnsolicitedFrames <- frameIn
			}

		}
	}
}

// This function currently assumes that every single function that expects a callback
// sets the callback id as the last byte in the payload.
func (s *Layer) sendThread() {
	for request := range s.requestQueue {
		var seqNo byte

		if request.ReceivesCallback {
			seqNo = s.getSequenceNumber()
			request.Payload = append(request.Payload, seqNo)
			s.callbacks[seqNo] = request.Callback
		}

		if request.Payload == nil {
			request.Payload = []byte{}
		}

		var frame = frame.NewRequestFrame(append([]byte{request.FunctionID}, request.Payload...))
		attempts := 0

	retry:
		if request.HasReturn {
			s.lastRequestFuncID = request.FunctionID
		}

		s.frameLayer.Write(frame)

		if request.HasReturn {
			select {
			case response := <-s.responses:
				if response.IsCan() {
					// Hopefully we won't collide again if we wait for 10ms :)
					time.Sleep(100 * time.Millisecond)
					if attempts > 3 {
						s.logger.Println("alert: TOO MANY RETRIES")
						request.ReturnCallback(errors.New("Too many retries sending command"), nil)
						return
					}

					attempts++
					goto retry // https://xkcd.com/292/
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
				s.logger.Println("warn: session lock timeout")
			}
		}

	}
}

func (s *Layer) getSequenceNumber() byte {
	if s.sequenceNumber == maxSequenceNumber {
		s.sequenceNumber = minSequenceNumber
	} else {
		s.sequenceNumber = s.sequenceNumber + 1
	}

	return s.sequenceNumber
}
