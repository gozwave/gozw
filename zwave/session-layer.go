package zwave

import (
	"fmt"

	"github.com/bjyoungblood/gozw/zwave/command-class"
	"github.com/bjyoungblood/gozw/zwave/frame"
)

func (s *ZWaveSessionLayer) handleApplicationCommand(cmd *ApplicationCommandHandler, frame *frame.Frame) bool {
	cc := cmd.CommandData[0]

	if cc == commandclass.CommandClassSecurity {
		switch cmd.CommandData[1] {

		case commandclass.CommandSecurityMessageEncapsulation, commandclass.CommandSecurityMessageEncapsulationNonceGet:
			// @todo determine whether to bother with sequenced messages

			// 1. decrypt message
			// 2. if it's the first half of a sequenced message, wait for the second half
			// 2.5  if it's an EncapsulationGetNonce, then send a NonceReport back to the sender
			// 3. if it's the second half of a sequenced message, reassemble the payloads
			// 4. emit the payload back to the session layer

			data := commandclass.ParseSecurityMessageEncapsulation(cmd.CommandData)
			msg, err := s.securityLayer.DecryptMessage(data)

			if msg[0] == commandclass.CommandClassSecurity && msg[1] == commandclass.CommandNetworkKeyVerify {
				s.securityLayer.SecurityFrameHandler(cmd, frame)
				return true
			}

			if err != nil {
				fmt.Println("error handling encrypted message", err)
				return false
			}

			cmd.CommandData = msg
			cc = cmd.CommandData[0]

		case commandclass.CommandSecurityNonceGet,
			commandclass.CommandSecurityNonceReport,
			commandclass.CommandSecuritySchemeReport,
			commandclass.CommandNetworkKeyVerify:
			s.securityLayer.SecurityFrameHandler(cmd, frame)
			return true
		}
	}

	if callback, ok := s.applicationCommandHandlers[cc]; ok {
		go callback(cmd, frame)
		return true
	}

	return false
}
