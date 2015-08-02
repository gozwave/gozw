package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
	"github.com/davecgh/go-spew/spew"
)

type DoorLock struct {
	UsersNumber byte
	UserCodes   map[byte]commandclass.UserCodeReport

	LockStatus byte

	node *Node

	receiveUsersNumber chan byte
}

func NewDoorLock(node *Node) *DoorLock {
	return &DoorLock{
		UserCodes:          map[byte]commandclass.UserCodeReport{},
		node:               node,
		receiveUsersNumber: make(chan byte),
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

func (d *DoorLock) initialize(node *Node) {
	d.node = node
	d.receiveUsersNumber = make(chan byte)

	if d.UserCodes == nil {
		d.UserCodes = make(map[byte]commandclass.UserCodeReport)
	}
}

func (d *DoorLock) GetSupportedUserCount() (byte, error) {
	if d.UsersNumber != 0 {
		return d.UsersNumber, nil
	}

	d.node.SendCommand(
		commandclass.CommandClassUserCode,
		commandclass.CommandUsersNumberGet,
	)

	select {
	case <-d.receiveUsersNumber:
		return d.UsersNumber, nil
	case <-time.After(time.Second * 5):
		return 0, errors.New("Timed out waiting for report")
	}
}

func (d *DoorLock) LoadUserCode(userId byte) error {
	return d.node.SendCommand(
		commandclass.CommandClassUserCode,
		commandclass.CommandUserCodeGet,
		userId,
	)
}

func (d *DoorLock) LoadAllUserCodes() error {
	var i byte

	max, err := d.GetSupportedUserCount()
	if err != nil {
		return err
	}

	for i = 0; i < max; i++ {
		err := d.LoadUserCode(i)
		time.Sleep(1 * time.Second)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DoorLock) handleAlarmCommandClass(cmd serialapi.ApplicationCommand) {
	// This is special handling code that will probably only work with yale locks
	notif := commandclass.ParseAlarmReport(cmd.CommandData)
	switch notif.Type {
	case 0x70:
		if notif.Level == 0x00 {
			fmt.Println("Master code changed")
		} else {
			fmt.Println("User added", notif.Level)
			d.LoadUserCode(notif.Level)
		}
	case 0xA1:
		if notif.Level == 0x01 {
			fmt.Println("Keypad limit exceeded")
		} else {
			fmt.Println("Physical tampering")
		}
	case 0x16:
		fmt.Println("Manual unlock")
	case 0x19:
		fmt.Println("RF operate unlock")
	case 0x15:
		fmt.Println("Manual lock")
	case 0x18:
		fmt.Println("RF operate lock")
	case 0x12:
		fmt.Println("keypad lock by user", notif.Level)
		d.LoadUserCode(notif.Level)
	case 0x13:
		fmt.Println("keypad unlock by user", notif.Level)
		d.LoadUserCode(notif.Level)
	case 0x09:
		fmt.Println("deadbolt jammed")
	case 0xA9:
		fmt.Println("dead battery; lock inoperable")
	case 0xA8:
		fmt.Println("critical battery")
	case 0xA7:
		fmt.Println("low battery")
	case 0x1B:
		fmt.Println("auto re-lock syscle completed")
	case 0x71:
		fmt.Println("duplicate pin code error")
	case 0x82:
		fmt.Println("power restored")
	case 0x21:
		fmt.Println("user deleted", notif.Level)
	}
}

func (d *DoorLock) handleDoorLockCommandClass(cmd serialapi.ApplicationCommand) {
	fmt.Println("Door lock cc")
	spew.Dump(cmd)
}

func (d *DoorLock) handleUserCodeCommandClass(cmd serialapi.ApplicationCommand) {
	switch cmd.CommandData[1] {
	case commandclass.CommandUserCodeReport:
		d.receiveUserCodeReport(commandclass.ParseUserCodeReport(cmd.CommandData))
	case commandclass.CommandUsersNumberReport:
		d.receiveUsersNumberReport(commandclass.ParseUsersNumberReport(cmd.CommandData))
	default:
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

func (d *DoorLock) receiveUsersNumberReport(number byte) {
	d.UsersNumber = number
	d.receiveUsersNumber <- number
}
