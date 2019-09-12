package session

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/protocol"
	"go.uber.org/zap"
)

const (
	minSequenceNumber = 1
	maxSequenceNumber = 127
)

// ILayer is an interface for  a  session layer.
type ILayer interface {
	MakeRequest(request *Request)
	SendFrameDirect(req *frame.Frame)
	UnsolicitedFramesChan() chan frame.Frame
}

// Layer contains a session layer.
type Layer struct {
	frameLayer        frame.ILayer
	UnsolicitedFrames chan frame.Frame
	lastRequestFuncID byte
	responses         chan frame.Frame
	sequenceNumber    byte
	callbacks         map[byte]CallbackFunc
	requestQueue      chan *Request
	l                 *zap.Logger
	ctx               context.Context
}

// NewSessionLayer will return a new session layer.
func NewSessionLayer(ctx context.Context, frameLayer frame.ILayer, logger *zap.Logger) *Layer {
	session := &Layer{
		frameLayer:        frameLayer,
		UnsolicitedFrames: make(chan frame.Frame, 10),
		lastRequestFuncID: 0,
		responses:         make(chan frame.Frame, 1),
		sequenceNumber:    0,
		callbacks:         map[byte]CallbackFunc{},
		requestQueue:      make(chan *Request, 10),
		l:                 logger,
		ctx:               ctx,
	}

	go session.receiveThread()
	go session.sendThread()

	return session
}

// MakeRequest will queue a request.
func (s *Layer) MakeRequest(request *Request) {
	s.requestQueue <- request
}

// SendFrameDirect should only be called inside a callback.
func (s *Layer) SendFrameDirect(req *frame.Frame) {
	s.frameLayer.Write(req)
}

// UnsolicitedFramesChan will return the unsolicited frames channel.
func (s *Layer) UnsolicitedFramesChan() chan frame.Frame {
	return s.UnsolicitedFrames
}

func (s *Layer) receiveThread() {
	for {
		select {
		case frameIn := <-s.frameLayer.GetOutputChannel():
			s.l.Debug("frame recieved")

			if frameIn.IsResponse() {
				s.l.Debug("was response")

				if frameIn.Payload[0] == s.lastRequestFuncID {
					select {
					case s.responses <- frameIn:
					default:
					}

					s.lastRequestFuncID = 0
				} else {
					s.l.Warn("received an unexpected response frame",
						zap.String("expected", fmt.Sprint(s.lastRequestFuncID)),
						zap.String("actual", fmt.Sprint(frameIn.Payload[0])),
					)
				}
			} else {
				var callbackID byte

				if s.lastRequestFuncID != 0 {
					s.l.Warn("REQUEST/RESPONSE COLLISION; SENDING CAN FRAME AND RETRYING PREVIOUS SEND")
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
					s.l.Warn("got unknown callback for func: ", zap.String("callback", hex.EncodeToString([]byte{frameIn.Payload[0]})))
					callbackID = 0
				}

				if callback, ok := s.callbacks[callbackID]; ok {
					go callback(frameIn)
				} else {
					s.UnsolicitedFrames <- frameIn
				}

			}
		case <-s.ctx.Done():
			s.l.Info("stopping session receive thread")
			return
		}
	}
}

// This function currently assumes that every single function that expects a callback
// sets the callback id as the last byte in the payload.
func (s *Layer) sendThread() {
	for {
		select {
		case request := <-s.requestQueue:
			var seqNo byte

			s.l.Debug("received request")

			if request.ReceivesCallback {
				seqNo = s.getSequenceNumber()
				request.Payload = append(request.Payload, seqNo)
				s.callbacks[seqNo] = request.Callback
			}

			if request.Payload == nil {
				request.Payload = []byte{}
			}

			s.l.Debug("creating request frame")

			var frame = frame.NewRequestFrame(append([]byte{request.FunctionID}, request.Payload...))
			attempts := 0

		retry:
			if request.HasReturn {
				s.lastRequestFuncID = request.FunctionID
			}

			s.l.Debug("writing frame")

			s.frameLayer.Write(frame)

			if request.HasReturn {
				select {
				case response := <-s.responses:
					if response.IsCan() {
						// Hopefully we won't collide again if we wait for 10ms :)
						time.Sleep(100 * time.Millisecond)
						if attempts > 3 {
							s.l.Error("too many retries")
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
					s.l.Warn("session lock timeout")
				}
			}
		case <-s.ctx.Done():
			s.l.Info("stopping session send thread")
			return
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
