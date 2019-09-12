package serialapi

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/protocol"
	"github.com/gozwave/gozw/session"
)

type transmitStatus struct {
	Status byte
	TxTime uint16
}

// SendData will send data to a node.
func (s *Layer) SendData(nodeID byte, payload []byte) (txTime uint16, err error) {

	transmitDone := make(chan bool)
	retStatus := make(chan error)
	txStatus := make(chan transmitStatus)

	payload = append([]byte{nodeID, byte(len(payload))}, payload...)
	payload = append(payload, protocol.TransmitOptionAck)

	request := &session.Request{
		FunctionID:       protocol.FnSendData,
		Payload:          payload,
		HasReturn:        true,
		ReceivesCallback: true,
		Lock:             true,
		Release:          transmitDone,
		Timeout:          10 * time.Second,

		ReturnCallback: func(err error, ret *frame.Frame) bool {
			if err != nil {
				transmitDone <- true
				retStatus <- err
				return false
			}

			if ret.Payload[1] == 0 {
				transmitDone <- true
				retStatus <- errors.New("SendData: transmit buffer overflow")
			} else {
				retStatus <- nil
			}

			return true
		},

		Callback: func(cbFrame frame.Frame) {
			status := transmitStatus{}
			status.Status = cbFrame.Payload[2]
			if len(cbFrame.Payload) == 5 {
				status.TxTime = binary.BigEndian.Uint16(cbFrame.Payload[3:5])
			}

			transmitDone <- true
			txStatus <- status
		},
	}

	s.sessionLayer.MakeRequest(request)

	err = <-retStatus
	if err != nil {
		return 0, err
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
