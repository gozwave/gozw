package serialapi

import (
	"time"

	"github.com/helioslabs/gozw/zwave/protocol"
	"github.com/helioslabs/gozw/zwave/session"
)

// WARNING: This can (and often will) cause the device to get a new USB address,
// rendering the serial port's file descriptor invalid.
func (s *Layer) SoftReset() {

	request := &session.Request{
		FunctionID: protocol.FnSerialAPISoftReset,
		HasReturn:  false,
	}

	s.sessionLayer.MakeRequest(request)

	time.Sleep(1500 * time.Millisecond)

}
