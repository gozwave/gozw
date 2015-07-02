package zwave

import "fmt"

type AddRemoveNodeCallback struct {
	CommandId      byte
	CallbackId     byte
	Status         byte
	Source         byte
	Length         byte
	Basic          uint8
	Generic        uint8
	Specific       uint8
	CommandClasses []byte
}

func ParseAddNodeCallback(payload []byte) *AddRemoveNodeCallback {
	val := &AddRemoveNodeCallback{
		CommandId:  payload[0],
		CallbackId: payload[1],
		Status:     payload[2],
		Source:     payload[3],
		Length:     payload[4],
	}

	fmt.Println(payload)

	if val.Length == 0 {
		return val
	}

	if val.Length >= 1 {
		val.Basic = payload[5]
	}

	if val.Length >= 2 {
		val.Generic = payload[6]
	}

	if val.Length >= 3 {
		val.Specific = payload[7]
	}

	if val.Length >= 4 {
		val.CommandClasses = payload[8:]
	}

	return val
}
