package application

import (
	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/protocol"
)

type Thermostat struct {
	node *Node
}

func IsThermostat(node *Node) bool {
	if node.GenericDeviceClass != protocol.GenericTypeThermostat {
		return false
	}

	switch node.SpecificDeviceClass {
	case protocol.SpecificTypeSetbackScheduleThermostat,
		protocol.SpecificTypeSetbackThermostat,
		protocol.SpecificTypeSetpointThermostat,
		protocol.SpecificTypeThermostatGeneral,
		protocol.SpecificTypeThermostatGeneralV2,
		protocol.SpecificTypeThermostatHeating:
		return true
	default:
		// Not sure how to handle these other device types yet, since I don't have any
		return false
	}
}

func NewThermostat(node *Node) *Thermostat {
	return &Thermostat{
		node: node,
	}
}

func (t *Thermostat) initialize(node *Node) {
	t.node = node
}

func (t *Thermostat) SetpointSet(setpointType commandclass.ThermostatSetpointType, temperature float64) error {
	payload, err := commandclass.NewThermostatSetpointSet(setpointType, commandclass.Temperature{
		Scale: commandclass.SetpointScaleFarenheit,
		Value: temperature,
	})

	if err != nil {
		return err
	}

	return t.node.SendCommand(payload[0], payload[1], payload[2:]...)
}
