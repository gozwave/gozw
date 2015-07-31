package commandclass

const (
	CommandAlarmVersion uint8 = 0x01
	CommandAlarmGet           = 0x04
	CommandAlarmReport        = 0x05
)

type AlarmReport struct {
	Type  uint8
	Level uint8
}

// @todo handle multiple versions, since this one has a lot
func ParseAlarmReport(payload []byte) AlarmReport {
	return AlarmReport{
		Type:  payload[2],
		Level: payload[3],
	}
}
