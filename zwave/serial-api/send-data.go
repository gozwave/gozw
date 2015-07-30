package serialapi

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/protocol"
	"github.com/bjyoungblood/gozw/zwave/session"
)

type TransmitStatus struct {
	Status uint8
	TxTime uint16
}

func (s *SerialAPILayer) SendData(nodeId byte, payload []byte) (txTime uint16, err error) {

	transmitDone := make(chan bool)
	retStatus := make(chan error)
	txStatus := make(chan TransmitStatus)

	payload = append([]byte{nodeId, uint8(len(payload))}, payload...)
	payload = append(payload, protocol.TransmitOptionAck)

	request := &session.Request{
		FunctionId:       protocol.FnSendData,
		Payload:          payload,
		HasReturn:        true,
		ReceivesCallback: true,
		Lock:             true,
		Release:          transmitDone,
		Timeout:          10 * time.Second,

		ReturnCallback: func(err error, ret *frame.Frame) bool {
			if ret.Payload[1] == 0 {
				transmitDone <- true
				retStatus <- errors.New("SendData: transmit buffer overflow")
			} else {
				retStatus <- nil
			}

			return true
		},

		Callback: func(cbFrame frame.Frame) {
			status := TransmitStatus{}
			status.Status = cbFrame.Payload[2]
			if len(cbFrame.Payload) == 5 {
				status.TxTime = binary.BigEndian.Uint16(cbFrame.Payload[3:5])
			}

			txStatus <- status
		},
	}

	s.sessionLayer.MakeRequest(request)

	err = <-retStatus
	if err != nil {
		return
	}

	status := <-txStatus
	switch status.Status {
	case protocol.TransmitCompleteOk:
		return status.TxTime, nil
	case protocol.TransmitCompleteNoAck:
		return status.TxTime, errors.New("Transmit complete: no ack from destination")
	case protocol.TransmitCompleteFail:
		return status.TxTime, errors.New("Transmit failure: network busy/jammed")
	case protocol.TransmitRoutingNotIdle:
		return status.TxTime, errors.New("Transmit failure: routing not idle")
	case protocol.TransmitCompleteNoRoute:
		return status.TxTime, errors.New("Transmit complete: no route")
	default:
		return status.TxTime, fmt.Errorf("Unknown tranmission status: %d", status.Status)
	}
}
