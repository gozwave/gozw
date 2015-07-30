package zwave

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/bjyoungblood/gozw/zwave/commandclass"
	"github.com/bjyoungblood/gozw/zwave/frame"
)

type SessionLayer interface {
	SetManager(manager *Manager)

	GetUnsolicitedFrames() chan frame.Frame
	ApplicationNodeInformation(deviceOptions uint8, genericType uint8, specificType uint8, supportedCommandClasses []uint8)
	SetSerialAPIReady(ready bool)
	SendData(nodeId uint8, data []byte) (*frame.Frame, error)
	SendDataSecure(nodeId uint8, data []byte) error

	sendDataUnsafe(nodeId uint8, data []byte) (*frame.Frame, error)
}

// @todo: ack timeouts, retransmission, backoff, etc.

type ZWaveSessionLayer struct {
	manager       *Manager
	frameLayer    frame.FrameLayer
	securityLayer SecurityLayer

	UnsolicitedFrames chan frame.Frame

	lastRequestedFn uint8
	responses       chan frame.Frame

	// maps sequence number to callback
	callbacks map[uint8]Callback

	// maps command class to callback
	applicationCommandHandlers map[uint8]CommandClassHandlerCallback

	currentSequenceNumber uint8

	execLock *sync.Mutex
}

func NewSessionLayer(frameLayer *frame.SerialFrameLayer) *ZWaveSessionLayer {
	session := &ZWaveSessionLayer{
		frameLayer: frameLayer,

		UnsolicitedFrames: make(chan frame.Frame),

		lastRequestedFn: 0,
		responses:       make(chan frame.Frame),

		callbacks: map[uint8]Callback{},

		applicationCommandHandlers: map[uint8]CommandClassHandlerCallback{},

		execLock: &sync.Mutex{},
	}

	session.securityLayer = NewSecurityLayer(session)

	go session.readFrames()

	return session
}

func (s *ZWaveSessionLayer) GetUnsolicitedFrames() chan frame.Frame {
	return s.UnsolicitedFrames
}

func (s *ZWaveSessionLayer) SetManager(manager *Manager) {
	s.manager = manager
}

func (s *ZWaveSessionLayer) ApplicationNodeInformation(
	deviceOptions uint8,
	genericType uint8,
	specificType uint8,
	supportedCommandClasses []uint8,
) {

	s.execLock.Lock()
	defer s.execLock.Unlock()
	defer runtime.Gosched()

	payload := []byte{
		FnApplicationNodeInformation,
		deviceOptions,
		genericType,
		specificType,
		uint8(len(supportedCommandClasses)),
	}

	payload = append(payload, supportedCommandClasses...)

	s.write(frame.NewRequestFrame(payload))

}

func (s *ZWaveSessionLayer) SendData(nodeId uint8, data []byte) (*frame.Frame, error) {
	s.execLock.Lock()
	defer s.execLock.Unlock()
	defer runtime.Gosched()

	return s.sendDataUnsafe(nodeId, data)
}

func (s *ZWaveSessionLayer) SendDataSecure(nodeId uint8, data []byte) error {
	s.execLock.Lock()
	defer s.execLock.Unlock()
	defer runtime.Gosched()

	return s.securityLayer.sendDataSecure(nodeId, data, false)
}

func (s *ZWaveSessionLayer) sendDataUnsafe(nodeId uint8, data []byte) (*frame.Frame, error) {
	done := make(chan CallbackResult)

	callback := func(callbackFrame frame.Frame) {
		// @todo implement me better maybe
		done <- CallbackResult{
			frame: &callbackFrame,
			err:   nil,
		}
	}

	seqNo := s.registerCallback(callback)
	defer s.unregisterCallback(seqNo)

	payload := []byte{
		FnSendData,
		nodeId,
		uint8(len(data)),
	}

	payload = append(payload, data...)
	payload = append(payload, TransmitOptionAck) // @todo implement ability to choose options
	payload = append(payload, seqNo)

	frame := frame.NewRequestFrame(payload)

	s.write(frame)

	result := <-done
	return result.frame, result.err
}

func (s *ZWaveSessionLayer) handleApplicationCommand(cmd *ApplicationCommandHandler, frame *frame.Frame) bool {
	cc := cmd.CommandData[0]

	if cc == commandclass.CommandClassSecurity {
		switch cmd.CommandData[1] {

		case commandclass.CommandSecurityMessageEncapsulation, commandclass.CommandSecurityMessageEncapsulationNonceGet:
			// @todo determine whether to bother with sequenced messages

			// 1. decrypt message
			// 2. if it's the first half of a sequenced message, wait for the second half
			// 2.5  if it's an EncapsulationGetNonce, then send a NonceReport back to the sender
			// 3. if it's the second half of a sequenced message, reassemble the payloads
			// 4. emit the payload back to the session layer

			data := commandclass.ParseSecurityMessageEncapsulation(cmd.CommandData)
			msg, err := s.securityLayer.DecryptMessage(data)

			if msg[0] == commandclass.CommandClassSecurity && msg[1] == commandclass.CommandNetworkKeyVerify {
				s.securityLayer.SecurityFrameHandler(cmd, frame)
				return true
			}

			if err != nil {
				fmt.Println("error handling encrypted message", err)
				return false
			}

			cmd.CommandData = msg
			cc = cmd.CommandData[0]

		case commandclass.CommandSecurityNonceGet,
			commandclass.CommandSecurityNonceReport,
			commandclass.CommandSecuritySchemeReport,
			commandclass.CommandNetworkKeyVerify:
			s.securityLayer.SecurityFrameHandler(cmd, frame)
			return true
		}
	}

	if callback, ok := s.applicationCommandHandlers[cc]; ok {
		go callback(cmd, frame)
		return true
	}

	return false
}
