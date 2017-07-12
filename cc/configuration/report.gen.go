// THIS FILE IS AUTO-GENERATED BY ZWGEN
// DO NOT MODIFY

package configuration

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/gozwave/gozw/cc"
)

const CommandReport cc.CommandID = 0x06

func init() {
	gob.Register(Report{})
	cc.Register(cc.CommandIdentifier{
		CommandClass: cc.CommandClassID(0x70),
		Command:      cc.CommandID(0x06),
		Version:      1,
	}, NewReport)
}

func NewReport() cc.Command {
	return &Report{}
}

// <no value>
type Report struct {
	ParameterNumber byte

	Level struct {
		Size byte
	}

	ConfigurationValue []byte
}

func (cmd Report) CommandClassID() cc.CommandClassID {
	return 0x70
}

func (cmd Report) CommandID() cc.CommandID {
	return CommandReport
}

func (cmd Report) CommandIDString() string {
	return "CONFIGURATION_REPORT"
}

func (cmd *Report) UnmarshalBinary(data []byte) error {
	// According to the docs, we must copy data if we wish to retain it after returning

	payload := make([]byte, len(data))
	copy(payload, data)

	if len(payload) < 2 {
		return errors.New("Payload length underflow")
	}

	i := 2

	if len(payload) <= i {
		return fmt.Errorf("slice index out of bounds (.ParameterNumber) %d<=%d", len(payload), i)
	}

	cmd.ParameterNumber = payload[i]
	i++

	if len(payload) <= i {
		return fmt.Errorf("slice index out of bounds (.Level) %d<=%d", len(payload), i)
	}

	cmd.Level.Size = (payload[i] & 0x07)

	i += 1

	if len(payload) <= i {
		return fmt.Errorf("slice index out of bounds (.ConfigurationValue) %d<=%d", len(payload), i)
	}

	{
		length := (payload[1+2] >> 0) & 0x07
		cmd.ConfigurationValue = payload[i : i+int(length)]
		i += int(length)
	}

	return nil
}

func (cmd *Report) MarshalBinary() (payload []byte, err error) {
	payload = make([]byte, 2)
	payload[0] = byte(cmd.CommandClassID())
	payload[1] = byte(cmd.CommandID())

	payload = append(payload, cmd.ParameterNumber)

	{
		var val byte

		val |= (cmd.Level.Size) & byte(0x07)

		payload = append(payload, val)
	}

	if cmd.ConfigurationValue != nil && len(cmd.ConfigurationValue) > 0 {
		payload = append(payload, cmd.ConfigurationValue...)
	}

	return
}
