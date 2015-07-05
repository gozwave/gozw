package zwave

type ApplicationCommandHandler struct {
	CommandId     byte
	ReceiveStatus byte
	DstNodeId     uint8
	SrcNodeId     uint8
	CmdLength     uint8
	CommandData   []byte
	// @todo implement multicast functionality
}

func ParseApplicationCommandHandler(payload []byte) *ApplicationCommandHandler {
	if payload[0] == FnApplicationCommandHandler {
		return &ApplicationCommandHandler{
			CommandId:     payload[0],
			ReceiveStatus: payload[1],
			DstNodeId:     1, // always controller
			SrcNodeId:     payload[2],
			CmdLength:     payload[3],
			CommandData:   payload[4 : 4+payload[3]],
		}
	} else {
		return &ApplicationCommandHandler{
			CommandId:     payload[0],
			ReceiveStatus: payload[1],
			DstNodeId:     payload[2],
			SrcNodeId:     payload[3],
			CmdLength:     payload[4],
			CommandData:   payload[5 : 5+payload[4]],
		}
	}

}
