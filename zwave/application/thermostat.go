package application

import (
	"fmt"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
	"github.com/davecgh/go-spew/spew"
)

type Thermostat struct {
	node *Node

	CoolingSetpoint commandclass.ThermostatSetpoint
	HeatingSetpoint commandclass.ThermostatSetpoint
	Mode            commandclass.ThermostatMode
	OperatingState  commandclass.ThermostatOperatingState
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

func (t *Thermostat) SetpointGet(setpointType commandclass.ThermostatSetpointType, temperature float64) error {
	return t.node.SendCommand(
		commandclass.CommandClassThermostatSetpoint,
		commandclass.CommandThermostatSetpointGet,
		byte(setpointType),
	)
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

func (t *Thermostat) ModeGet() error {
	return t.node.SendCommand(
		commandclass.CommandClassThermostatMode,
		commandclass.CommandThermostatModeGet,
	)
}

func (t *Thermostat) OperatingStateGet() error {
	return t.node.SendCommand(
		commandclass.CommandClassThermostatOperatingState,
		commandclass.CommandThermostatOperatingStateGet,
	)
}

func (t *Thermostat) handleThermostatOperatingStateCommandClass(cmd serialapi.ApplicationCommand) {
	if cmd.CommandData[1] == commandclass.CommandThermostatOperatingStateReport {
		t.receiveOperatingStateReport(commandclass.ParseThermostatOperatingStateReport(cmd.CommandData))
	} else {
		spew.Dump(cmd)
	}
}

func (t *Thermostat) receiveOperatingStateReport(operatingState commandclass.ThermostatOperatingState) {
	t.OperatingState = operatingState
	switch t.OperatingState {
	case commandclass.ThermostatOperatingStateIdle:
		fmt.Println("Thermostat operating state: Idle")
	case commandclass.ThermostatOperatingStateHeating:
		fmt.Println("Thermostat operating state: Heating")
	case commandclass.ThermostatOperatingStateCooling:
		fmt.Println("Thermostat operating state: Cooling")
	case commandclass.ThermostatOperatingStateFanOnly:
		fmt.Println("Thermostat operating state: Fan Only")
	case commandclass.ThermostatOperatingStatePendingHeat:
		fmt.Println("Thermostat operating state: Pending Heat")
	case commandclass.ThermostatOperatingStatePendingCool:
		fmt.Println("Thermostat operating state: Pending Cool")
	case commandclass.ThermostatOperatingStateVentEconomizer:
		fmt.Println("Thermostat operating state: Vent Economizer")
	default:
		fmt.Printf("Thermostat operating state: unknown (%d)\n", t.Mode)
	}

	t.node.saveToDb()
}

func (t *Thermostat) handleThermostatModeCommandClass(cmd serialapi.ApplicationCommand) {
	if cmd.CommandData[1] == commandclass.CommandThermostatModeReport {
		t.receiveModeReport(commandclass.ParseThermostatModeReport(cmd.CommandData))
	} else {
		spew.Dump(cmd)
	}
}

func (t *Thermostat) receiveModeReport(mode commandclass.ThermostatMode) {
	t.Mode = mode
	switch t.Mode {
	case commandclass.ThermostatModeModeOff:
		fmt.Println("Thermostat mode: Off")
	case commandclass.ThermostatModeModeHeat:
		fmt.Println("Thermostat mode: Heat")
	case commandclass.ThermostatModeModeCool:
		fmt.Println("Thermostat mode: Cool")
	case commandclass.ThermostatModeModeAuto:
		fmt.Println("Thermostat mode: Auto")
	case commandclass.ThermostatModeModeAuxiliaryHeat:
		fmt.Println("Thermostat mode: Auxiliary Heat")
	case commandclass.ThermostatModeModeResume:
		fmt.Println("Thermostat mode: Resume")
	case commandclass.ThermostatModeModeFanOnly:
		fmt.Println("Thermostat mode: Fan Only")
	case commandclass.ThermostatModeModeFurnace:
		fmt.Println("Thermostat mode: Furnace")
	case commandclass.ThermostatModeModeDryAir:
		fmt.Println("Thermostat mode: Dry Air")
	case commandclass.ThermostatModeModeMoistAir:
		fmt.Println("Thermostat mode: Moist Air")
	case commandclass.ThermostatModeModeAutoChangeover:
		fmt.Println("Thermostat mode: Auto Changeover")
	default:
		fmt.Printf("Thermostat mode: unknown (%d)\n", t.Mode)
	}

	t.node.saveToDb()
}

func (t *Thermostat) handleThermostatSetpointCommandClass(cmd serialapi.ApplicationCommand) {
	if cmd.CommandData[1] == commandclass.CommandThermostatSetpointReport {
		report := commandclass.ParseThermostatSetpointReport(cmd.CommandData)
		t.receiveSetpointReport(report)
	} else {
		spew.Dump(cmd)
	}
}

func (t *Thermostat) receiveSetpointReport(setpoint commandclass.ThermostatSetpoint) {
	switch setpoint.Type {
	case commandclass.ThermostatSetpointTypeCooling:
		t.CoolingSetpoint = setpoint
		temperature, err := setpoint.GetTemperature()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("New cooling setpoint:", temperature.Value)
	case commandclass.ThermostatSetpointTypeHeating:
		t.HeatingSetpoint = setpoint
		temperature, err := setpoint.GetTemperature()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("New heating setpoint:", temperature.Value)
	default:
		fmt.Println("Unknown setpoint update")
		spew.Dump(setpoint)
		return
	}

	t.node.saveToDb()
}
