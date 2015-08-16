// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package networkmanagementproxy

import "errors"

// <no value>

type NodeInfoCachedReport struct {
	SeqNo byte

	Properties1 struct {
		Age byte

		Status byte
	}

	Properties2 struct {
		Capability byte

		Listening bool
	}

	Properties3 struct {
		Security byte

		Sensor byte

		Opt bool
	}

	BasicDeviceClass byte

	GenericDeviceClass byte

	SpecificDeviceClass byte

	NonSecureCommandClass []byte

	SecurityScheme0CommandClass []byte
}

func (cmd *NodeInfoCachedReport) UnmarshalBinary(payload []byte) error {
	i := 0

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.SeqNo = payload[i]
	i++

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.Properties1.Age = (payload[i] & 0x0F)

	cmd.Properties1.Status = (payload[i] & 0xF0) >> 4

	i += 1

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.Properties2.Capability = (payload[i] & 0x7F)

	if payload[i]&0x80 == 0x80 {
		cmd.Properties2.Listening = true
	} else {
		cmd.Properties2.Listening = false
	}

	i += 1

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.Properties3.Security = (payload[i] & 0x0F)

	cmd.Properties3.Sensor = (payload[i] & 0x70) >> 4

	if payload[i]&0x80 == 0x80 {
		cmd.Properties3.Opt = true
	} else {
		cmd.Properties3.Opt = false
	}

	i += 1

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.BasicDeviceClass = payload[i]
	i++

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.GenericDeviceClass = payload[i]
	i++

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.SpecificDeviceClass = payload[i]
	i++

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	{
		fieldStart := i
		for ; i < len(payload) && payload[i] != 0xF1; i++ {
		}
		cmd.NonSecureCommandClass = payload[fieldStart:i]
	}

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	i += 1 // skipping MARKER

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.SecurityScheme0CommandClass = payload[i:]

	return nil
}

func (cmd *NodeInfoCachedReport) MarshalBinary() (payload []byte, err error) {

	payload = append(payload, cmd.SeqNo)

	{
		var val byte

		val |= (cmd.Properties1.Age) & byte(0x0F)

		val |= (cmd.Properties1.Status << byte(4)) & byte(0xF0)

		payload = append(payload, val)
	}

	{
		var val byte

		val |= (cmd.Properties2.Capability) & byte(0x7F)

		if cmd.Properties2.Listening {
			val |= byte(0x80) // flip bits on
		} else {
			val &= ^byte(0x80) // flip bits off
		}

		payload = append(payload, val)
	}

	{
		var val byte

		val |= (cmd.Properties3.Security) & byte(0x0F)

		val |= (cmd.Properties3.Sensor << byte(4)) & byte(0x70)

		if cmd.Properties3.Opt {
			val |= byte(0x80) // flip bits on
		} else {
			val &= ^byte(0x80) // flip bits off
		}

		payload = append(payload, val)
	}

	payload = append(payload, cmd.BasicDeviceClass)

	payload = append(payload, cmd.GenericDeviceClass)

	payload = append(payload, cmd.SpecificDeviceClass)

	{
		if cmd.NonSecureCommandClass != nil && len(cmd.NonSecureCommandClass) > 0 {
			payload = append(payload, cmd.NonSecureCommandClass...)
		}
		payload = append(payload, 0xF1)
	}

	payload = append(payload, 0xF1) // marker

	payload = append(payload, cmd.SecurityScheme0CommandClass...)

	return
}
