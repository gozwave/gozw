package zwave

import "fmt"

type SessionLayer struct {
	frameLayer *FrameLayer
}

func NewSessionLayer(frameLayer *FrameLayer) *SessionLayer {
	return &SessionLayer{
		frameLayer: frameLayer,
	}
}

func (session *SessionLayer) ExecuteCommand(commandId uint8, payload []byte) {
	frame := NewRequestFrame(append([]byte{commandId}, payload...))

	session.frameLayer.Write(frame)

	response := <-session.frameLayer.GetOutput()

	fmt.Println(response)
}
