// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package doorlocklogging

import (
	"encoding/gob"
	"errors"
)

func init() {
	gob.Register(RecordsSupportedReport{})
}

// <no value>
type RecordsSupportedReport struct {
	MaxRecordsStored byte
}

func (cmd RecordsSupportedReport) CommandClassID() byte {
	return 0x4C
}

func (cmd RecordsSupportedReport) CommandID() byte {
	return byte(CommandRecordsSupportedReport)
}

func (cmd *RecordsSupportedReport) UnmarshalBinary(data []byte) error {
	// According to the docs, we must copy data if we wish to retain it after returning

	payload := make([]byte, len(data))
	copy(payload, data)

	if len(payload) < 2 {
		return errors.New("Payload length underflow")
	}

	i := 2

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.MaxRecordsStored = payload[i]
	i++

	return nil
}

func (cmd *RecordsSupportedReport) MarshalBinary() (payload []byte, err error) {
	payload = make([]byte, 2)
	payload[0] = cmd.CommandClassID()
	payload[1] = cmd.CommandID()

	payload = append(payload, cmd.MaxRecordsStored)

	return
}