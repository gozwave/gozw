package commandclass

const (
	CommandThermostatModeSet             byte = 0x01
	CommandThermostatModeGet                  = 0x02
	CommandThermostatModeReport               = 0x03
	CommandThermostatModeSupportedGet         = 0x04
	CommandThermostatModeSupportedReport      = 0x05
)

const ThermostatModeMask byte = 0x1F

type ThermostatMode byte

const (
	ThermostatModeModeOff            ThermostatMode = 0x00
	ThermostatModeModeHeat                          = 0x01
	ThermostatModeModeCool                          = 0x02
	ThermostatModeModeAuto                          = 0x03
	ThermostatModeModeAuxiliaryHeat                 = 0x04
	ThermostatModeModeResume                        = 0x05
	ThermostatModeModeFanOnly                       = 0x06
	ThermostatModeModeFurnace                       = 0x07
	ThermostatModeModeDryAir                        = 0x08
	ThermostatModeModeMoistAir                      = 0x09
	ThermostatModeModeAutoChangeover                = 0x0A
)

func NewThermostatModeSet(mode ThermostatMode) []byte {
	return []byte{
		CommandClassThermostatMode,
		CommandThermostatModeSet,
		byte(mode) & ThermostatModeMask,
	}
}

func ParseThermostatModeReport(payload []byte) ThermostatMode {
	return ThermostatMode(payload[2] & ThermostatModeMask)
}
