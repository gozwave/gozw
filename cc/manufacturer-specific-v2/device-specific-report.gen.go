// THIS FILE IS AUTO-GENERATED BY ZWGEN
// DO NOT MODIFY

package manufacturerspecificv2

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/gozwave/gozw/cc"
)

const CommandDeviceSpecificReport cc.CommandID = 0x07

func init() {
	gob.Register(DeviceSpecificReport{})
	cc.Register(cc.CommandIdentifier{
		CommandClass: cc.CommandClassID(0x72),
		Command:      cc.CommandID(0x07),
		Version:      2,
	}, NewDeviceSpecificReport)
}

func NewDeviceSpecificReport() cc.Command {
	return &DeviceSpecificReport{}
}

// <no value>
type DeviceSpecificReport struct {
	Properties1 struct {
		DeviceIdType byte
	}

	Properties2 struct {
		DeviceIdDataLengthIndicator byte

		DeviceIdDataFormat byte
	}

	DeviceIdData []byte
}

func (cmd DeviceSpecificReport) CommandClassID() cc.CommandClassID {
	return 0x72
}

func (cmd DeviceSpecificReport) CommandID() cc.CommandID {
	return CommandDeviceSpecificReport
}

func (cmd DeviceSpecificReport) CommandIDString() string {
	return "DEVICE_SPECIFIC_REPORT"
}

func (cmd *DeviceSpecificReport) UnmarshalBinary(data []byte) error {
	// According to the docs, we must copy data if we wish to retain it after returning

	payload := make([]byte, len(data))
	copy(payload, data)

	if len(payload) < 2 {
		return errors.New("Payload length underflow")
	}

	i := 2

	if len(payload) <= i {
		return fmt.Errorf("slice index out of bounds (.Properties1) %d<=%d", len(payload), i)
	}

	cmd.Properties1.DeviceIdType = (payload[i] & 0x07)

	i += 1

	if len(payload) <= i {
		return fmt.Errorf("slice index out of bounds (.Properties2) %d<=%d", len(payload), i)
	}

	cmd.Properties2.DeviceIdDataLengthIndicator = (payload[i] & 0x1F)

	cmd.Properties2.DeviceIdDataFormat = (payload[i] & 0xE0) >> 5

	i += 1

	if len(payload) <= i {
		return fmt.Errorf("slice index out of bounds (.DeviceIdData) %d<=%d", len(payload), i)
	}

	{
		length := (payload[1+2] >> 0) & 0x1F
		cmd.DeviceIdData = payload[i : i+int(length)]
		i += int(length)
	}

	return nil
}

func (cmd *DeviceSpecificReport) MarshalBinary() (payload []byte, err error) {
	payload = make([]byte, 2)
	payload[0] = byte(cmd.CommandClassID())
	payload[1] = byte(cmd.CommandID())

	{
		var val byte

		val |= (cmd.Properties1.DeviceIdType) & byte(0x07)

		payload = append(payload, val)
	}

	{
		var val byte

		val |= (cmd.Properties2.DeviceIdDataLengthIndicator) & byte(0x1F)

		val |= (cmd.Properties2.DeviceIdDataFormat << byte(5)) & byte(0xE0)

		payload = append(payload, val)
	}

	if cmd.DeviceIdData != nil && len(cmd.DeviceIdData) > 0 {
		payload = append(payload, cmd.DeviceIdData...)
	}

	return
}
