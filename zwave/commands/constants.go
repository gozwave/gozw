package commands

const (
	CommandClassDoorLock             = 0x62
	CommandClassThermostatSetpointV3 = 0x43
)

const (
	DoorLockOperationSet  = 0x01
	ThermostatSetpointSet = 0x01
)

const (
	DoorLockStatusUnlocked = 0x00
	DoorLockStatusLocked   = 0xFF
)
