// THIS FILE IS AUTO-GENERATED BY ZWGEN
// DO NOT MODIFY

package applicationstatus

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/gozwave/gozw/cc"
)

const CommandApplicationRejectedRequest cc.CommandID = 0x02

func init() {
	gob.Register(ApplicationRejectedRequest{})
	cc.Register(cc.CommandIdentifier{
		CommandClass: cc.CommandClassID(0x22),
		Command:      cc.CommandID(0x02),
		Version:      1,
	}, NewApplicationRejectedRequest)
}

func NewApplicationRejectedRequest() cc.Command {
	return &ApplicationRejectedRequest{}
}

// <no value>
type ApplicationRejectedRequest struct {
	Status byte
}

func (cmd ApplicationRejectedRequest) CommandClassID() cc.CommandClassID {
	return 0x22
}

func (cmd ApplicationRejectedRequest) CommandID() cc.CommandID {
	return CommandApplicationRejectedRequest
}

func (cmd ApplicationRejectedRequest) CommandIDString() string {
	return "APPLICATION_REJECTED_REQUEST"
}

func (cmd *ApplicationRejectedRequest) UnmarshalBinary(data []byte) error {
	// According to the docs, we must copy data if we wish to retain it after returning

	payload := make([]byte, len(data))
	copy(payload, data)

	if len(payload) < 2 {
		return errors.New("Payload length underflow")
	}

	i := 2

	if len(payload) <= i {
		return fmt.Errorf("slice index out of bounds (.Status) %d<=%d", len(payload), i)
	}

	cmd.Status = payload[i]
	i++

	return nil
}

func (cmd *ApplicationRejectedRequest) MarshalBinary() (payload []byte, err error) {
	payload = make([]byte, 2)
	payload[0] = byte(cmd.CommandClassID())
	payload[1] = byte(cmd.CommandID())

	payload = append(payload, cmd.Status)

	return
}
