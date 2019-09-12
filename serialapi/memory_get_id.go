package serialapi

import (
	"encoding/binary"
	"errors"

	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/protocol"
	"github.com/gozwave/gozw/session"
)

// MemoryGetID will get the home/node id.
func (s *Layer) MemoryGetID() (homeID uint32, nodeID byte, err error) {

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionID: protocol.FnMemoryGetID,
		HasReturn:  true,
		ReturnCallback: func(err error, ret *frame.Frame) bool {
			done <- ret
			return false
		},
	}

	s.sessionLayer.MakeRequest(request)
	ret := <-done

	if ret == nil {
		return 0, 0, errors.New("Error getting home/node id")
	}

	homeID = binary.BigEndian.Uint32(ret.Payload[1:5])
	nodeID = ret.Payload[5]

	return
}
