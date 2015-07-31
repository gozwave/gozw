package commandclass

const (
	CommandUserCodeVersion   uint = 0x01
	CommandUserCodeGet            = 0x02
	CommandUserCodeReport         = 0x03
	CommandUserCodeSet            = 0x01
	CommandUsersNumberGet         = 0x04
	CommandUsersNumberReport      = 0x05
)

const (
	UserCodeReportAvailableNotSet         uint8 = 0x00
	UserCodeReportOccupied                      = 0x01
	UserCodeReportReservedByAdministrator       = 0x02
	UserCodeReportStatusNotAvailable            = 0xFE
)

const (
	UserCodeSetAvailableNotSet         uint8 = 0x00
	UserCodeSetOccupied                      = 0x01
	UserCodeSetReservedByAdministrator       = 0x02
	UserCodeSetStatusNotAvailable            = 0xFE
)

type UserCodeReport struct {
	UserIdentifier uint8
	UserStatus     uint8
	Code           string
}

func ParseUserCodeReport(payload []byte) UserCodeReport {
	return UserCodeReport{
		UserIdentifier: payload[2],
		UserStatus:     payload[3],
		Code:           string(payload[4:]),
	}
}
