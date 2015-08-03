package commandclass

const (
	CommandDoorLockOperationSet        = 0x01
	CommandDoorLockOperationGet        = 0x02
	CommandDoorLockOperationReport     = 0x03
	CommandDoorLockConfigurationSet    = 0x04
	CommandDoorLockConfigurationGet    = 0x05
	CommandDoorLockConfigurationReport = 0x06
)

type DoorLockOperationReport struct {
	DoorLockMode      byte
	OutsideHandleMode byte // @todo
	InsideHandleMode  byte // @todo
	DoorCondition     byte
	TimeoutMinutes    byte
	TimeoutSeconds    byte
}

func ParseDoorLockOperationReport(payload []byte) DoorLockOperationReport {
	return DoorLockOperationReport{
		DoorLockMode:   payload[2],
		DoorCondition:  payload[4],
		TimeoutMinutes: payload[5],
		TimeoutSeconds: payload[6],
	}
}
