package commands

func NewDoorLockCommand() []byte {
	return []byte{
		CommandClassDoorLock,
		DoorLockOperationSet,
		byte(0),
	}
}
