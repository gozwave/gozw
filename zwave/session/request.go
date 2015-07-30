package session

import (
	"time"

	"github.com/bjyoungblood/gozw/zwave/frame"
)

// type CallbackType int

type CallbackFunc func(frame.Frame)

type Request struct {
	FunctionId byte
	Payload    []byte

	HasReturn      bool
	ReturnCallback func(error, *frame.Frame) bool

	ReceivesCallback bool
	Callback         CallbackFunc
	Lock             bool
	Release          chan bool
	Timeout          time.Duration
}
