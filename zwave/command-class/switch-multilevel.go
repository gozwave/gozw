package commandclass

const (
	SwitchMultilevelSet              byte = 0x01
	SwitchMultilevelGet                   = 0x02
	SwitchMultilevelReport                = 0x03
	SwitchMultilevelStartLevelChange      = 0x04
	SwitchMultilevelStopLevelChange       = 0x05
)

func NewSwitchMultilevelCommand(level byte) []byte {
	return []byte{
		CommandClassSensorMultilevel,
		SwitchMultilevelSet,
		level,
	}
}
