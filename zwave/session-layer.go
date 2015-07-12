package zwave

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/bjyoungblood/gozw/zwave/commandclass"
)

const (
	MinSequenceNumber = 1
	MaxSequenceNumber = 127
)

const (
	SendDataInsecure = false
	SendDataSecure   = true
)

const (
	AddRemoveNodeReadyTimeout = time.Second * 10
	AddRemoveNodeFoundTimeout = time.Second * 30 // @todo recommended is 60, but 30 is nicer for testing
)

type CommandClassHandlerCallback func(*ApplicationCommandHandler, *Frame)

type Callback func(Frame)
type CallbackResult struct {
	frame *Frame
	err   error
}

type AddRemoveNodeResult struct {
	node *Node
	err  error
}

type SessionLayer interface {
	SetManager(manager *Manager)

	GetUnsolicitedFrames() chan Frame
	ApplicationNodeInformation(deviceOptions uint8, genericType uint8, specificType uint8, supportedCommandClasses []uint8)
	SetDefault()
	AddNodeToNetwork() (*Node, error)
	RemoveNodeFromNetwork() (*Node, error)
	GetVersion() (*VersionResponse, error)
	MemoryGetId() (*MemoryGetIdResponse, error)
	GetInitAppData() (*NodeListResponse, error)
	GetSerialApiCapabilities() (*SerialApiCapabilitiesResponse, error)
	GetNodeProtocolInfo(nodeId uint8) (*NodeProtocolInfoResponse, error)
	SetSerialAPIReady(ready bool)
	SendData(nodeId uint8, data []byte) (*Frame, error)
	SendDataSecure(nodeId uint8, data []byte) error

	registerApplicationCommandHandler(commandClass uint8, callback CommandClassHandlerCallback)
	sendDataUnsafe(nodeId uint8, data []byte) (*Frame, error)
	requestNodeInformationFrame(nodeId uint8) error
	isNodeFailing(nodeId uint8) (bool, error)
	removeFailedNode(nodeId uint8) (*Frame, error)
}

// @todo: ack timeouts, retransmission, backoff, etc.

type ZWaveSessionLayer struct {
	manager       *Manager
	frameLayer    FrameLayer
	securityLayer SecurityLayer

	UnsolicitedFrames chan Frame

	lastRequestedFn uint8
	responses       chan Frame

	// maps sequence number to callback
	callbacks map[uint8]Callback

	// maps command class to callback
	applicationCommandHandlers map[uint8]CommandClassHandlerCallback

	currentSequenceNumber uint8

	execLock *sync.Mutex
}

func NewSessionLayer(frameLayer *SerialFrameLayer) *ZWaveSessionLayer {
	session := &ZWaveSessionLayer{
		frameLayer: frameLayer,

		UnsolicitedFrames: make(chan Frame),

		lastRequestedFn: 0,
		responses:       make(chan Frame),

		callbacks: map[uint8]Callback{},

		applicationCommandHandlers: map[uint8]CommandClassHandlerCallback{},

		execLock: &sync.Mutex{},
	}

	session.securityLayer = NewSecurityLayer(session)

	go session.readFrames()

	return session
}

func (s *ZWaveSessionLayer) GetUnsolicitedFrames() chan Frame {
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

	s.write(NewRequestFrame(payload))

}

func (s *ZWaveSessionLayer) SetDefault() {
	done := make(chan bool)

	s.execLock.Lock()
	defer s.execLock.Unlock()
	defer runtime.Gosched()

	callback := func(frame Frame) {
		done <- true
	}

	seqNo := s.registerCallback(callback)
	defer s.unregisterCallback(seqNo)

	payload := []byte{
		FnSetDefault,
		seqNo,
	}

	s.write(NewRequestFrame(payload))

	// @todo timeout
	<-done
}

func (s *ZWaveSessionLayer) AddNodeToNetwork() (*Node, error) {
	done := make(chan *AddRemoveNodeResult)

	// Just keep in mind that this will lock all other z-wave traffic until the
	// process completes (it's supposed to, but you should be aware anyway)
	s.execLock.Lock()
	defer s.execLock.Unlock()
	defer runtime.Gosched()

	var newNode *Node = nil
	timeout := time.NewTimer(0)

	callback := func(cb Frame) {
		payload := ParseAddNodeCallback(cb.Payload)
		fmt.Println("ADD NODE: FRAME", payload)
		// at each step, don't forget to kick the timer to the correct value
		switch {
		case payload.Status == AddNodeStatusLearnReady:
			timeout.Reset(AddRemoveNodeFoundTimeout)
			fmt.Println("ADD NODE: learn ready")
		case payload.Status == AddNodeStatusNodeFound:
			timeout.Reset(AddRemoveNodeFoundTimeout)
			fmt.Println("ADD NODE: node found")
		case payload.Status == AddNodeStatusAddingSlave:
			timeout.Reset(AddRemoveNodeFoundTimeout)
			newNode = NewNodeFromAddNodeCallback(s.manager, payload)
			fmt.Println("ADD NODE: adding slave node")
		case payload.Status == AddNodeStatusAddingController:
			// hey, i just met you, and this is crazy
			// but it could happen, so implement me maybe
			timeout.Reset(AddRemoveNodeFoundTimeout)
			fmt.Println("ADD NODE: adding controller node")
		case payload.Status == AddNodeStatusProtocolDone:
			fmt.Println("ADD NODE: protocol done")

			s.AddRemoveNodeStop(FnAddNodeToNetwork)

			if newNode.IsSecure() {
				s.securityLayer.includeSecureNode(newNode)
			}

			if newNode != nil {
				done <- &AddRemoveNodeResult{
					node: newNode,
					err:  nil,
				}
			} else {
				done <- &AddRemoveNodeResult{
					node: nil,
					err:  errors.New("Unknown error adding node (this should not happen)"),
				}
			}
		case payload.Status == AddNodeStatusFailed:
			fmt.Println("ADD NODE: failed")
			s.AddRemoveNodeStop(FnAddNodeToNetwork)
			done <- &AddRemoveNodeResult{
				node: nil,
				err:  errors.New("Failed to add node correctly"),
			}
		}
	}

	seqNo := s.registerCallback(callback)
	defer s.unregisterCallback(seqNo)

	frame := NewRequestFrame([]byte{
		FnAddNodeToNetwork,
		AddNodeAny | AddNodeOptionNetworkWide | AddNodeOptionNormalPower,
		seqNo,
	})

	s.write(frame)
	timeout.Reset(AddRemoveNodeReadyTimeout)

	select {
	case result := <-done:

		return result.node, result.err
	case <-timeout.C:
		return nil, errors.New("Timed out adding node")
	}
}

// Believe it or not, this function is exempt from execLock
func (s *ZWaveSessionLayer) AddRemoveNodeStop(funcId uint8) {
	done := make(chan bool)

	callback := func(cb Frame) {
		payload := ParseAddNodeCallback(cb.Payload)
		fmt.Println("ADD NODE2: FRAME", payload)
		switch {
		case payload.Status == AddNodeStatusDone:
			fmt.Println("ADD NODE: done")
		default:
			fmt.Printf("ADD NODE: unexpected status 0x%x\n", payload.Status)
		}

		// this should happen regardless of what the previous status was
		s.write(NewRequestFrame([]byte{
			funcId,
			AddNodeStop,
			0x00, // @todo should this be 0x0 or omitted entirely?
		}))

		done <- true
	}

	seqNo := s.registerCallback(callback)
	defer s.unregisterCallback(seqNo)
	fmt.Println("remove node stop seq", seqNo)

	frame := NewRequestFrame([]byte{
		funcId,
		AddNodeStop,
		seqNo,
	})

	s.write(frame)
	timeout := time.NewTimer(AddRemoveNodeReadyTimeout)

	// both cases are noops
	select {
	case <-done:
	case <-timeout.C:
	}
}

// @todo remove node is currently breaking things after it runs
func (s *ZWaveSessionLayer) RemoveNodeFromNetwork() (*Node, error) {
	done := make(chan *AddRemoveNodeResult)

	// Just keep in mind that this will lock all other z-wave traffic until the
	// process completes (it's supposed to, but you should be aware anyway)
	s.execLock.Lock()
	defer s.execLock.Unlock()
	defer runtime.Gosched()

	var newNode *Node = nil
	timeout := time.NewTimer(0)

	callback := func(cb Frame) {
		payload := ParseAddNodeCallback(cb.Payload)
		fmt.Println("REMOVE NODE: FRAME", payload)
		// at each step, don't forget to kick the timer to the correct value
		switch {
		case payload.Status == RemoveNodeStatusLearnReady:
			timeout.Reset(AddRemoveNodeFoundTimeout)
			fmt.Println("REMOVE NODE: learn ready")

		case payload.Status == RemoveNodeStatusNodeFound:
			timeout.Reset(AddRemoveNodeFoundTimeout)
			fmt.Println("REMOVE NODE: node found")

		case payload.Status == RemoveNodeStatusRemovingSlave || payload.Status == RemoveNodeStatusRemovingController:
			timeout.Reset(AddRemoveNodeFoundTimeout)
			newNode = NewNodeFromAddNodeCallback(s.manager, payload)
			fmt.Println("REMOVE NODE: remove node")

			s.AddRemoveNodeStop(FnRemoveNodeFromNetwork)
			if newNode != nil {
				done <- &AddRemoveNodeResult{
					node: newNode,
					err:  nil,
				}
			} else {
				done <- &AddRemoveNodeResult{
					node: nil,
					err:  errors.New("Unknown error removing node (this should not happen)"),
				}
			}

		case payload.Status == RemoveNodeStatusFailed:
			fmt.Println("REMOVE NODE: failed")
			s.AddRemoveNodeStop(FnRemoveNodeFromNetwork)
			done <- &AddRemoveNodeResult{
				node: nil,
				err:  errors.New("Failed to remove node correctly"),
			}
		}
	}

	seqNo := s.registerCallback(callback)
	fmt.Println("remove node seq", seqNo)
	defer s.unregisterCallback(seqNo)

	frame := NewRequestFrame([]byte{
		FnRemoveNodeFromNetwork,
		RemoveNodeAny,
		seqNo,
	})

	s.write(frame)
	timeout.Reset(AddRemoveNodeReadyTimeout)

	select {
	case result := <-done:
		return result.node, result.err
	case <-timeout.C:
		return nil, errors.New("Timed out removing node")
	}
}

func (s *ZWaveSessionLayer) GetVersion() (*VersionResponse, error) {
	response, err := s.writeSimple(FnGetVersion, []byte{})
	if err != nil {
		return nil, err
	}

	return ParseVersionResponse(response.Payload), nil
}

func (s *ZWaveSessionLayer) MemoryGetId() (*MemoryGetIdResponse, error) {
	response, err := s.writeSimple(FnMemoryGetId, []byte{})
	if err != nil {
		return nil, err
	}

	return ParseMemoryGetIdResponse(response.Payload), nil
}

func (s *ZWaveSessionLayer) GetInitAppData() (*NodeListResponse, error) {
	response, err := s.writeSimple(FnGetInitAppData, []byte{})
	if err != nil {
		return nil, err
	}

	return ParseNodeListResponse(response.Payload), nil
}

func (s *ZWaveSessionLayer) isNodeFailing(nodeId uint8) (bool, error) {
	response, err := s.writeSimple(FnIsNodeFailed, []byte{nodeId})
	if err != nil {
		return false, err
	}

	if response.Payload[1] == 0 {
		return false, nil
	}

	return true, nil
}

func (s *ZWaveSessionLayer) GetSerialApiCapabilities() (*SerialApiCapabilitiesResponse, error) {
	response, err := s.writeSimple(FnSerialApiCapabilities, []byte{})
	if err != nil {
		return nil, err
	}

	return ParseSerialApiCapabilitiesResponse(response.Payload), nil
}

func (s *ZWaveSessionLayer) GetNodeProtocolInfo(nodeId uint8) (*NodeProtocolInfoResponse, error) {
	response, err := s.writeSimple(FnGetNodeProtocolInfo, []byte{nodeId})
	if err != nil {
		return nil, err
	}

	return ParseNodeProtocolInfoResponse(response.Payload), nil
}

func (s *ZWaveSessionLayer) SetSerialAPIReady(ready bool) {
	var rdy byte
	if ready {
		rdy = 1
	} else {
		rdy = 0
	}

	s.execLock.Lock()
	defer s.execLock.Unlock()
	defer runtime.Gosched()

	payload := []byte{FnSerialAPIReady, rdy}
	s.write(NewRequestFrame(payload))
}

func (s *ZWaveSessionLayer) requestNodeInformationFrame(nodeId uint8) error {
	_, err := s.writeSimple(FnRequestNodeInfo, []byte{nodeId})
	return err
}

func (s *ZWaveSessionLayer) SendData(nodeId uint8, data []byte) (*Frame, error) {
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

func (s *ZWaveSessionLayer) sendDataUnsafe(nodeId uint8, data []byte) (*Frame, error) {
	done := make(chan CallbackResult)

	callback := func(callbackFrame Frame) {
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

	frame := NewRequestFrame(payload)

	s.write(frame)

	result := <-done
	return result.frame, result.err
}

// Writes a single byte function call and awaits the response frame
func (s *ZWaveSessionLayer) writeSimple(funcId uint8, payload []byte) (*Frame, error) {
	s.execLock.Lock()
	defer s.execLock.Unlock()
	defer runtime.Gosched()

	payload = append([]byte{funcId}, payload...)

	s.write(NewRequestFrame(payload))

	// @todo correct timeout implementation?
	select {
	case response := <-s.responses:
		return &response, nil
	case <-time.After(10 * time.Second):
		return nil, errors.New("Request timeout")
	}
}

func (s *ZWaveSessionLayer) removeFailedNode(nodeId uint8) (*Frame, error) {
	done := make(chan CallbackResult)

	callback := func(callbackFrame Frame) {
		done <- CallbackResult{
			frame: &callbackFrame,
			err:   nil,
		}
	}

	seqNo := s.registerCallback(callback)
	defer s.unregisterCallback(seqNo)

	payload := []byte{
		FnRemoveFailingNode,
		nodeId,
		seqNo,
	}

	frame := NewRequestFrame(payload)

	s.write(frame)

	result := <-done
	return result.frame, result.err
}

// DO NOT CALL THIS WITHOUT ALREADY HAVING OBTAINED THE LOCK
func (s *ZWaveSessionLayer) write(frame *Frame) {
	s.lastRequestedFn = frame.Payload[0]
	s.frameLayer.Write(frame)
}

func (s *ZWaveSessionLayer) readFrames() {
	for frame := range s.frameLayer.GetOutputChannel() {
		s.processFrame(frame)
	}
}

// for each incoming frame, determine how to handle it based on whether it is a
// return value (response), a callback (request), or an unsolicited frame (request)
func (s *ZWaveSessionLayer) processFrame(frame Frame) {
	if frame.IsResponse() {
		// handle frame as a response
		if frame.Payload[0] == s.lastRequestedFn {

			// performs a non-blocking send
			select {
			case s.responses <- frame:
				// noop
			default:
				// noop
			}

			// Clear last requested function
			s.lastRequestedFn = 0

		} else {
			fmt.Println("Received an unexpected response frame: ", frame)
		}
	} else {
		// handle frame as a callback

		var callbackId uint8

		// find the callback id in the payload
		switch frame.Payload[0] {
		case FnAddNodeToNetwork, FnRemoveNodeFromNetwork, FnSendData:
			callbackId = frame.Payload[1]

		case FnApplicationCommandHandler, FnApplicationCommandHandlerBridge:
			cmd := ParseApplicationCommandHandler(frame.Payload)

			handled := s.handleApplicationCommand(cmd, &frame)
			if handled {
				return
			}

			// never a callback
			callbackId = 0

		case FnApplicationControllerUpdate:
			// never a callback
			callbackId = 0

		default:
			fmt.Println("session-layer: Potentially missed callback!")
			callbackId = 0
		}

		// if we have a registered callback, remove it from the list of registered
		// callbacks, and call it (asynchronously, in case it makes further calls)
		if callback, ok := s.callbacks[callbackId]; ok {
			go callback(frame)
		} else {
			s.UnsolicitedFrames <- frame
		}
	}
}

func (s *ZWaveSessionLayer) handleApplicationCommand(cmd *ApplicationCommandHandler, frame *Frame) bool {
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

// This will prevent these command classes from reaching the application layer
// until they are unregistered
func (s *ZWaveSessionLayer) registerApplicationCommandHandler(commandClass uint8, callback CommandClassHandlerCallback) {
	s.applicationCommandHandlers[commandClass] = callback
}

func (s *ZWaveSessionLayer) unregisterApplicationCommandHandler(commandClass uint8) {
	delete(s.applicationCommandHandlers, commandClass)
}

func (s *ZWaveSessionLayer) registerCallback(callback Callback) uint8 {
	seqNo := s.getSequenceNumber()

	s.callbacks[seqNo] = callback

	return seqNo
}

func (s *ZWaveSessionLayer) unregisterCallback(seqNo uint8) {
	delete(s.callbacks, seqNo)
}

func (s *ZWaveSessionLayer) getSequenceNumber() uint8 {
	if s.currentSequenceNumber == MaxSequenceNumber {
		s.currentSequenceNumber = MinSequenceNumber
	} else {
		s.currentSequenceNumber = s.currentSequenceNumber + 1
	}

	return s.currentSequenceNumber
}
