package session

import (
	"time"

	"github.com/bjyoungblood/gozw/zwave/frame"
)

// type CallbackType int

type Request struct {
	FunctionId byte
	Payload    []byte

	HasReturn      bool
	ReturnCallback func(error, *frame.Frame) bool

	ReceivesCallback bool
	Callback         func(frame.Frame) bool
	Exclusive        bool
	Timeout          time.Duration
}
