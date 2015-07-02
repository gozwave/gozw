package zwave

import (
	"fmt"
	"runtime"
	"sync"
)

// @todo: ack timeouts, retransmission, backoff, etc.

type SessionLayer struct {
	frameLayer *FrameLayer

	writeLock *sync.Mutex
	readState chan bool

	ApplicationFrames chan Frame
}

func NewSessionLayer(frameLayer *FrameLayer) *SessionLayer {
	session := &SessionLayer{
		frameLayer:        frameLayer,
		writeLock:         &sync.Mutex{},
		readState:         make(chan bool),
		ApplicationFrames: make(chan Frame),
	}

	go session.read()

	return session
}

func (session *SessionLayer) WaitForFrame() Frame {
	return <-session.frameLayer.frameOutput
}

func (session *SessionLayer) ExecuteCommand(commandId uint8, payload []byte) Frame {
	frame := NewRequestFrame()
	framePayload := &GenericPayload{
		CommandId: commandId,
		Payload:   payload,
	}

	frame.Payload = framePayload.Marshal()

	session.writeLock.Lock()
	session.pauseReads()
	session.frameLayer.Write(frame)
	// @todo handle timeouts, transmission, incorrect responses, etc.
	response := <-session.frameLayer.frameOutput
	session.resumeReads()
	session.writeLock.Unlock()
	runtime.Gosched()

	return response
}

func (session *SessionLayer) ExecuteCommandNoWait(commandId uint8, payload []byte) {
	frame := NewRequestFrame()
	framePayload := &GenericPayload{
		CommandId: commandId,
		Payload:   payload,
	}

	frame.Payload = framePayload.Marshal()

	session.writeLock.Lock()
	session.frameLayer.Write(frame)
	session.writeLock.Unlock()
	runtime.Gosched()
}

func (session *SessionLayer) pauseReads() {
	session.readState <- false
}

func (session *SessionLayer) resumeReads() {
	session.readState <- true
}

func (session *SessionLayer) read() {
	for {
	read:
		select {
		case continueReading := <-session.readState:
			if !continueReading {
				for continueReading = range session.readState {
					if continueReading {
						goto read
					}
				}
			}

		case frame := <-session.frameLayer.frameOutput:
			fmt.Println("Application frame:", frame)
			session.ApplicationFrames <- frame
		}
	}
}
