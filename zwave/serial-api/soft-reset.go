package serialapi

import (
	"time"

	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/session"
)

// WARNING: This can (and often will) cause the device to get a new USB address,
// rendering the serial port's file descriptor invalid.
func (s *SerialAPILayer) SoftReset() {

	request := &session.Request{
		FunctionId: protocol.FnSerialApiSoftReset,
		Payload:    []byte{},
		HasReturn:  false,
	}

	s.sessionLayer.MakeRequest(request)

	time.Sleep(1500 * time.Millisecond)

}
