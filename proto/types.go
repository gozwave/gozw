package proto

import (
	"encoding/gob"

	"github.com/bjyoungblood/gozw/zwave/command-class"
)

func init() {
	gob.Register(Event{})
	gob.Register(IdentEvent{})
	gob.Register(NodeEvent{})
	gob.Register(DoorLockEvent{})
	gob.Register(UserCodeEvent{})
	gob.Register(ErrorEvent{})
	gob.Register(commandclass.AlarmReport{})
}

type Event struct {
	Payload interface{}
}

type GatewayStatus byte

const (
	Offline GatewayStatus = iota
	Initializing
	Online
	Error
)

type IdentEvent struct {
	HomeId uint32
}

type NodeEvent struct {
	NodeId byte
	Event  interface{}
}

type DoorLockEvent struct {
	NodeId     byte
	LockStatus byte
}

type UserCodeEvent struct {
	NodeId         byte
	Status         byte
	UserIdentifier byte
	UserCode       string
}

type ErrorEvent struct {
	Error error
}
