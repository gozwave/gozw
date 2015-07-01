package zwave

import (
	"runtime"
	"sync"
)

// @todo: ack timeouts, retransmission, backoff, etc.

type SessionLayer struct {
	frameLayer *FrameLayer

	writeLock *sync.Mutex
}

func NewSessionLayer(frameLayer *FrameLayer) *SessionLayer {
	return &SessionLayer{
		frameLayer: frameLayer,
		writeLock:  &sync.Mutex{},
	}
}

func (session *SessionLayer) ExecuteCommand(commandId uint8, payload []byte) *Frame {
	frame := NewRequestFrame()
	framePayload := &GenericPayload{
		CommandId: commandId,
		Payload:   payload,
	}

	frame.Payload = framePayload.Marshal()

	session.writeLock.Lock()
	session.frameLayer.Write(frame)
	response := <-session.frameLayer.GetOutput()
	session.writeLock.Unlock()
	runtime.Gosched()

	return response
}
