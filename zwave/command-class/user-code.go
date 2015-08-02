package commandclass

const (
	CommandUserCodeVersion   byte = 0x01
	CommandUserCodeGet            = 0x02
	CommandUserCodeReport         = 0x03
	CommandUserCodeSet            = 0x01
	CommandUsersNumberGet         = 0x04
	CommandUsersNumberReport      = 0x05
)

const (
	UserCodeReportAvailableNotSet         byte = 0x00
	UserCodeReportOccupied                     = 0x01
	UserCodeReportReservedByAdministrator      = 0x02
	UserCodeReportStatusNotAvailable           = 0xFE
)

const (
	UserCodeSetAvailableNotSet         byte = 0x00
	UserCodeSetOccupied                     = 0x01
	UserCodeSetReservedByAdministrator      = 0x02
	UserCodeSetStatusNotAvailable           = 0xFE
)

type UserCodeReport struct {
	UserIdentifier byte
	UserStatus     byte
	Code           string
}

func ParseUserCodeReport(payload []byte) UserCodeReport {
	return UserCodeReport{
		UserIdentifier: payload[2],
		UserStatus:     payload[3],
		Code:           string(payload[4:]),
	}
}

func ParseUsersNumberReport(payload []byte) byte {
	return payload[2]
}
