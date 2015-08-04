package commandclass

const (
	CommandThermostatOperatingStateVersion byte = 0x01
	CommandThermostatOperatingStateGet          = 0x02
	CommandThermostatOperatingStateReport       = 0x03
)

const ThermostatOperatingStateMask byte = 0x0F

type ThermostatOperatingState byte

/* Values used for Thermostat Operating State Report command */
const (
	ThermostatOperatingStateIdle           ThermostatOperatingState = 0x00
	ThermostatOperatingStateHeating                                 = 0x01
	ThermostatOperatingStateCooling                                 = 0x02
	ThermostatOperatingStateFanOnly                                 = 0x03
	ThermostatOperatingStatePendingHeat                             = 0x04
	ThermostatOperatingStatePendingCool                             = 0x05
	ThermostatOperatingStateVentEconomizer                          = 0x06
)

func ParseThermostatOperatingStateReport(payload []byte) ThermostatOperatingState {
	return ThermostatOperatingState(payload[2] & ThermostatOperatingStateMask)
}
