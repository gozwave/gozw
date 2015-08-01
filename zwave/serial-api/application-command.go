package serialapi

import "github.com/bjyoungblood/gozw/zwave/protocol"

type ApplicationCommand struct {
	CommandId     byte
	ReceiveStatus byte
	DstNodeId     byte
	SrcNodeId     byte
	CmdLength     byte
	CommandData   []byte
	// @todo implement multicast functionality (maybe? only needed for bridge library)
}

func parseApplicationCommand(payload []byte) ApplicationCommand {
	if payload[0] == protocol.FnApplicationCommandHandler {
		return ApplicationCommand{
			CommandId:     payload[0],
			ReceiveStatus: payload[1],
			DstNodeId:     1, // always controller
			SrcNodeId:     payload[2],
			CmdLength:     payload[3],
			CommandData:   payload[4 : 4+payload[3]],
		}
	} else {
		return ApplicationCommand{
			CommandId:     payload[0],
			ReceiveStatus: payload[1],
			DstNodeId:     payload[2],
			SrcNodeId:     payload[3],
			CmdLength:     payload[4],
			CommandData:   payload[5 : 5+payload[4]],
		}
	}

}
