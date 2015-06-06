package commands

import "fmt"

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
		64,
	}

	fmt.Println(buf)

	return buf
}
