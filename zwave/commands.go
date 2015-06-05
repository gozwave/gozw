package zwave

const CommandClassNetworkManagementProxy uint8 = 0x52
const NodeInfoCachedGet uint8 = 0x03

func MakeNodeInfoCachedGetPacket() []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, CommandClassNetworkManagementProxy)
	bytes = append(bytes, NodeInfoCachedGet)
	bytes = append(bytes, uint8(0x20))
	bytes = append(bytes, uint8(0x02))
	bytes = append(bytes, uint8(0x00))

	return bytes
}
