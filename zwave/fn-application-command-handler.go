package zwave

type ApplicationCommandHandlerBridge struct {
	CommandId     byte
	ReceiveStatus byte
	DstNodeId     uint8
	SrcNodeId     uint8
	CmdLength     uint8
	CommandData   []byte
	// @todo implement multicast functionality
}

func ParseApplicationCommandHandlerBridge(payload []byte) *ApplicationCommandHandlerBridge {
	return &ApplicationCommandHandlerBridge{
		CommandId:     payload[0],
		ReceiveStatus: payload[1],
		DstNodeId:     payload[2],
		SrcNodeId:     payload[3],
		CmdLength:     payload[4],
		CommandData:   payload[5 : 5+payload[4]],
	}
}
