// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package multichannelassociationv2

import "errors"

// <no value>

type MultiChannelAssociationReport struct {
	GroupingIdentifier byte

	MaxNodesSupported byte

	ReportsToFollow byte

	NodeId []byte
}

func (cmd *MultiChannelAssociationReport) UnmarshalBinary(payload []byte) error {
	i := 0

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.GroupingIdentifier = payload[i]
	i++

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.MaxNodesSupported = payload[i]
	i++

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.ReportsToFollow = payload[i]
	i++

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	{
		fieldStart := i
		for ; i < len(payload) && payload[i] != 0x00; i++ {
		}
		cmd.NodeId = payload[fieldStart:i]
	}

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	i += 1 // skipping MARKER

	return nil
}

func (cmd *MultiChannelAssociationReport) MarshalBinary() (payload []byte, err error) {

	payload = append(payload, cmd.GroupingIdentifier)

	payload = append(payload, cmd.MaxNodesSupported)

	payload = append(payload, cmd.ReportsToFollow)

	{
		if cmd.NodeId != nil && len(cmd.NodeId) > 0 {
			payload = append(payload, cmd.NodeId...)
		}
		payload = append(payload, 0x00)
	}

	payload = append(payload, 0x00) // marker

	return
}
