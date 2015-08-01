package commandclass

const (
	CommandVersionGet                byte = 0x11
	CommandVersionReport                  = 0x12
	CommandVersionCommandClassGet         = 0x13
	CommandVersionCommandClassReport      = 0x14
)

func NewVersionGet() []byte {
	return []byte{
		CommandClassVersion,
		CommandVersionGet,
	}
}
