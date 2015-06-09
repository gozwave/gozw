package commands

const (
	CommandClassDoorLock             = 0x62
	CommandClassThermostatSetpointV3 = 0x43
	CommandClassSecurity             = 0x98
)

const (
	// Door lock
	DoorLockOperationSet = 0x01

	// Thermostat
	ThermostatSetpointSet = 0x01

	// Security
	NetworkKeySet                        = 0x06
	NetworkKeyVerify                     = 0x07
	SecurityCommandsSupportedGet         = 0x02
	SecurityMessageEncapsulation         = 0x81
	SecurityMessageEncapsulationNonceGet = 0xC1
	SecurityNonceGet                     = 0x40
	SecurityNonceReport                  = 0x80
	SecuritySchemeGet                    = 0x04
)

const (
	DoorLockStatusUnlocked = 0x00
	DoorLockStatusLocked   = 0xFF
)
