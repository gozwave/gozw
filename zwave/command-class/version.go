package commandclass

const (
	CommandVersionGet                byte = 0x11
	CommandVersionReport                  = 0x12
	CommandVersionCommandClassGet         = 0x13
	CommandVersionCommandClassReport      = 0x14
)

type VersionCommandClassReport struct {
	CommandClass byte
	Version      byte
}

func NewVersionGet() []byte {
	return []byte{
		CommandClassVersion,
		CommandVersionGet,
	}
}

func ParseVersionCommandClassReport(payload []byte) VersionCommandClassReport {
	return VersionCommandClassReport{
		CommandClass: payload[2],
		Version:      payload[3],
	}
}
