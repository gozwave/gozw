package application

import (
	"fmt"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
	"github.com/davecgh/go-spew/spew"
)

type DoorLock struct {
	UserCodes map[byte]commandclass.UserCodeReport

	LockStatus byte

	node *Node
}

func NewDoorLock(node *Node) *DoorLock {
	return &DoorLock{
		UserCodes: map[byte]commandclass.UserCodeReport{},
		node:      node,
	}
}

func IsDoorLock(node *Node) bool {
	if node.GenericDeviceClass != protocol.GenericTypeEntryControl {
		return false
	}

	switch node.SpecificDeviceClass {
	case protocol.SpecificTypeDoorLock,
		protocol.SpecificTypeAdvancedDoorLock,
		protocol.SpecificTypeSecureKeypadDoorLock,
		protocol.SpecificTypeSecureKeypadDoorLockDeadbolt:
		return true
	default:
		// Not sure how to handle these other device types yet, since I don't have any
		return false
	}
}

func (d *DoorLock) handleDoorLockCommandClass(cmd serialapi.ApplicationCommand) {
	fmt.Println("Door lock cc")
	spew.Dump(cmd)
}

func (d *DoorLock) handleUserCodeCommandClass(cmd serialapi.ApplicationCommand) {
	if cmd.CommandData[1] == commandclass.CommandUserCodeReport {
		d.receiveUserCodeReport(commandclass.ParseUserCodeReport(cmd.CommandData))
	} else {
		spew.Dump(cmd.CommandData)
	}
}

func (d *DoorLock) receiveUserCodeReport(code commandclass.UserCodeReport) {
	fmt.Println("user code")
	if code.UserStatus == 0x0 { // code slot is available; don't save
		return
	}

	d.UserCodes[code.UserIdentifier] = code
	spew.Dump(code)
	d.node.saveToDb()
}
