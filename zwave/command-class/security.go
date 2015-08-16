package commandclass

import "errors"

const (
	CommandSecurityVersion                      byte = 0x01
	CommandNetworkKeySet                             = 0x06
	CommandNetworkKeyVerify                          = 0x07
	CommandSecurityCommandsSupportedGet              = 0x02
	CommandSecurityCommandsSupportedReport           = 0x03
	CommandSecurityMessageEncapsulation              = 0x81
	CommandSecurityMessageEncapsulationNonceGet      = 0xC1
	CommandSecurityNonceGet                          = 0x40
	CommandSecurityNonceReport                       = 0x80
	CommandSecuritySchemeGet                         = 0x04
	CommandSecuritySchemeInherit                     = 0x08
	CommandSecuritySchemeReport                      = 0x05
)

const SecurityCommandsSupportedReportCommandClassMark = 0xEF

type SecuritySchemeGet struct {
	CommandClass             byte
	Command                  byte
	SupportedSecuritySchemes byte
}

type SecurityNonceGet struct {
	CommandClass byte
	Command      byte
}

type SecurityNonceReport struct {
	CommandClass byte
	Command      byte
	Nonce        []byte
}

type SecurityCommandsSupportedReport struct {
	CommandClass             byte
	Command                  byte
	RemainingFrames          byte
	SupportedCommandClasses  []byte
	ControlledCommandClasses []byte
}

type SecurityMessageEncapsulation struct {
	CommandClass     byte
	Command          byte
	SenderNonce      []byte
	EncryptedPayload []byte
	ReceiverNonceID  byte
	Hmac             []byte
}

func NewSecuritySchemeGet() []byte {
	return []byte{
		byte(Security),
		CommandSecuritySchemeGet,
		0x0,
	}
}

func NewSecurityNonceGet() []byte {
	return []byte{
		byte(Security),
		CommandSecurityNonceGet,
	}
}

func NewSecurityNonceReport(nonce []byte) []byte {
	buf := []byte{
		byte(Security),
		CommandSecurityNonceReport,
	}

	return append(buf, nonce...)
}

func NewSecurityNetworkKeySet(key []byte) []byte {
	buf := []byte{
		byte(Security),
		CommandNetworkKeySet,
	}

	return append(buf, key...)
}

func NewSecurityMessageEncapsulation(iv, payload, hmac []byte, receiverNonceID byte) []byte {
	buf := []byte{
		byte(Security),
		CommandSecurityMessageEncapsulation,
	}

	buf = append(buf, iv...)
	buf = append(buf, payload...)
	buf = append(buf, receiverNonceID)
	buf = append(buf, hmac...)

	return buf
}

func ParseSecurityCommandsSupportedReport(data []byte) *SecurityCommandsSupportedReport {
	cc := &SecurityCommandsSupportedReport{
		CommandClass:    data[0],
		Command:         data[1],
		RemainingFrames: data[2],
	}

	supportedCommandClasses := []byte{}
	controlledCommandClasses := []byte{}

	var i int
	for i = 3; i < len(data); i++ {
		if data[i] == SecurityCommandsSupportedReportCommandClassMark {
			break
		}

		supportedCommandClasses = append(supportedCommandClasses, data[i])
	}

	i++ // skip command class mark

	for i < len(data) {
		controlledCommandClasses = append(controlledCommandClasses, data[i])
	}

	cc.SupportedCommandClasses = supportedCommandClasses
	cc.ControlledCommandClasses = controlledCommandClasses

	return cc
}

func ParseSecurityMessageEncapsulation(data []byte) *SecurityMessageEncapsulation {
	payloadLen := len(data) - 19

	cmd := &SecurityMessageEncapsulation{
		CommandClass:     data[0],
		Command:          data[1],
		SenderNonce:      data[2:10],
		EncryptedPayload: data[10 : 10+payloadLen],
		ReceiverNonceID:  data[10+payloadLen],
		Hmac:             data[11+payloadLen:],
	}

	return cmd
}

func ParseSecurityNonceReport(command []byte) *SecurityNonceReport {
	return &SecurityNonceReport{
		CommandClass: command[0],
		Command:      command[1],
		Nonce:        command[2:],
	}
}

func ParseCommandClassSecurity(command []byte) (interface{}, error) {
	if command[0] != byte(Security) {
		return nil, errors.New("Not security command class")
	}

	switch command[1] {
	case CommandSecurityNonceReport:
		return ParseSecurityNonceReport(command), nil
	default:
		return nil, errors.New("Unhandled command")
	}
}
