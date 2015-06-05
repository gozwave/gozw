package zwave

import "fmt"

const (
	CommandClassThermostatSetpointV3 = 0x43
)

const (
	ThermostatSetpointSet = 0x01
)

func NewThermostatSetpointCommand() []byte {

	var precision, scale, size uint8

	precision = 0
	scale = 1 << 3
	size = 1 // 8-bit signed

	// var value uint8 = 0x44

	buf := []byte{
		CommandClassThermostatSetpointV3,
		ThermostatSetpointSet,
		0x01, // heating
		precision | scale | size, // intentional bitwise OR
		0x48,
	}

	fmt.Println(buf)

	return buf
}
