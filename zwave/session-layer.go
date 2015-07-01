package zwave

// @todo: ack timeouts, retransmission, backoff, etc.

type SessionLayer struct {
	frameLayer *FrameLayer
}

func NewSessionLayer(frameLayer *FrameLayer) *SessionLayer {
	return &SessionLayer{
		frameLayer: frameLayer,
	}
}

func (session *SessionLayer) ExecuteCommand(commandId uint8, payload []byte) *Frame {
	frame := NewRequestFrame()
	framePayload := &GenericPayload{
		CommandId: commandId,
		Payload:   payload,
	}

	frame.Payload = framePayload.Marshal()

	session.frameLayer.Write(frame)

	response := <-session.frameLayer.GetOutput()

	return response
}
