package commandclass

import "errors"

const (
	CommandSecurityVersion                      uint8 = 0x01
	CommandNetworkKeySet                              = 0x06
	CommandNetworkKeyVerify                           = 0x07
	CommandSecurityCommandsSupportedGet               = 0x02
	CommandSecurityCommandsSupportedReport            = 0x03
	CommandSecurityMessageEncapsulation               = 0x81
	CommandSecurityMessageEncapsulationNonceGet       = 0xC1
	CommandSecurityNonceGet                           = 0x40
	CommandSecurityNonceReport                        = 0x80
	CommandSecuritySchemeGet                          = 0x04
	CommandSecuritySchemeInherit                      = 0x08
	CommandSecuritySchemeReport                       = 0x05
)

type SecuritySchemeGet struct {
	CommandClass             uint8
	Command                  uint8
	SupportedSecuritySchemes byte
}

type SecurityNonceGet struct {
	CommandClass uint8
	Command      uint8
}

type SecurityNonceReport struct {
	CommandClass uint8
	Command      uint8
	Nonce        []byte
}

func NewSecuritySchemeGet() []byte {
	return []byte{
		CommandClassSecurity,
		CommandSecuritySchemeGet,
		0x0,
	}
}

func NewSecurityNonceGet() []byte {
	return []byte{
		CommandClassSecurity,
		CommandSecurityNonceGet,
	}
}

func NewSecurityNonceReport() []byte {
	return []byte{
		CommandClassSecurity,
		CommandSecurityNonceReport,
	}
}

func ParseCommandClassSecurity(command []byte) (interface{}, error) {
	if command[0] != CommandClassSecurity {
		return nil, errors.New("Not security command class")
	}

	switch command[1] {
	case CommandSecurityNonceReport:
		return &SecurityNonceReport{
			CommandClass: command[0],
			Command:      command[1],
			Nonce:        command[2:],
		}, nil
	default:
		return nil, errors.New("Unhandled command")
	}
}
