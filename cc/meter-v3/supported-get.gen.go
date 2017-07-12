// THIS FILE IS AUTO-GENERATED BY ZWGEN
// DO NOT MODIFY

package meterv3

import (
	"encoding/gob"

	"github.com/gozwave/gozw/cc"
)

const CommandSupportedGet cc.CommandID = 0x03

func init() {
	gob.Register(SupportedGet{})
	cc.Register(cc.CommandIdentifier{
		CommandClass: cc.CommandClassID(0x32),
		Command:      cc.CommandID(0x03),
		Version:      3,
	}, NewSupportedGet)
}

func NewSupportedGet() cc.Command {
	return &SupportedGet{}
}

// <no value>
type SupportedGet struct {
}

func (cmd SupportedGet) CommandClassID() cc.CommandClassID {
	return 0x32
}

func (cmd SupportedGet) CommandID() cc.CommandID {
	return CommandSupportedGet
}

func (cmd SupportedGet) CommandIDString() string {
	return "METER_SUPPORTED_GET"
}

func (cmd *SupportedGet) UnmarshalBinary(data []byte) error {
	// According to the docs, we must copy data if we wish to retain it after returning

	return nil
}

func (cmd *SupportedGet) MarshalBinary() (payload []byte, err error) {
	payload = make([]byte, 2)
	payload[0] = byte(cmd.CommandClassID())
	payload[1] = byte(cmd.CommandID())

	return
}
