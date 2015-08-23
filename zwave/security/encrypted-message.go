package security

import (
	"errors"

	"github.com/helioslabs/gozw/zwave/command-class"
	"github.com/helioslabs/gozw/zwave/command-class/security"
)

type EncryptedMessage struct {
	SenderNonce      []byte
	EncryptedPayload []byte
	ReceiverNonceID  byte
	HMAC             []byte
}

func (cmd EncryptedMessage) CommandClassID() byte {
	return byte(commandclass.Security)
}

func (cmd EncryptedMessage) CommandID() byte {
	return byte(security.CommandMessageEncapsulation)
}

func (cmd *EncryptedMessage) UnmarshalBinary(data []byte) error {
	// According to the docs, we must copy data if we wish to retain it after returning

	if len(data) < 19 {
		return errors.New("Payload length underflow")
	}

	payload := make([]byte, len(data))
	copy(payload, data)

	cmd.SenderNonce = payload[2:10]
	cmd.EncryptedPayload = payload[10 : len(payload)-9]
	cmd.ReceiverNonceID = payload[10+len(cmd.EncryptedPayload)]
	cmd.HMAC = payload[11+len(cmd.EncryptedPayload):]

	return nil
}

func (cmd *EncryptedMessage) MarshalBinary() (payload []byte, err error) {
	payload = make([]byte, 0)

	payload = append(payload, cmd.CommandClassID())
	payload = append(payload, cmd.CommandID())
	payload = append(payload, cmd.SenderNonce...)
	payload = append(payload, cmd.EncryptedPayload...)
	payload = append(payload, cmd.ReceiverNonceID)
	payload = append(payload, cmd.HMAC...)

	return payload, nil
}
