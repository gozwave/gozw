package serialapi

import "github.com/gozwave/gozw/protocol"

// ApplicationCommand contains an application level command.
type ApplicationCommand struct {
	CommandID     byte
	ReceiveStatus byte
	DstNodeID     byte
	SrcNodeID     byte
	CmdLength     byte
	CommandData   []byte
	// @todo implement multicast functionality (maybe? only needed for bridge library)
}

func parseApplicationCommand(payload []byte) ApplicationCommand {
	if payload[0] == protocol.FnApplicationCommandHandler {
		return ApplicationCommand{
			CommandID:     payload[0],
			ReceiveStatus: payload[1],
			DstNodeID:     1, // always controller
			SrcNodeID:     payload[2],
			CmdLength:     payload[3],
			CommandData:   payload[4 : 4+payload[3]],
		}
	}

	return ApplicationCommand{
		CommandID:     payload[0],
		ReceiveStatus: payload[1],
		DstNodeID:     payload[2],
		SrcNodeID:     payload[3],
		CmdLength:     payload[4],
		CommandData:   payload[5 : 5+payload[4]],
	}
}
