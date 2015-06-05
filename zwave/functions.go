package zwave

const (
	ZwNodeList              = 0x02
	ZwSendData              = 0x13
	ZwAddNodeToNetwork      = 0x4a
	ZwRemoveNodeFromNetwork = 0x4b
	ZwSerialAPIReady        = 0xEF
)

const (
	AddNodeAny                  = 1
	AddNodeController           = 2
	AddNodeSlave                = 3
	AddNodeExisting             = 4
	AddNodeStop                 = 5
	AddNodeStopFailed           = 6
	AddNodeStatusSecurityFailed = 9
)

const (
	AddNodeOptionNormalPower = 0x80
	AddNodeOptionNetworkWide = 0x40
)

const (
	RemoveNodeAny        = AddNodeAny
	RemoveNodeController = AddNodeController
	RemoveNodeSlave      = AddNodeSlave
	RemoveNodeStop       = AddNodeStop
)

const (
	RemoveNodeOptionNormalPower = AddNodeOptionNormalPower
	RemoveNodeOptionNetworkWide = AddNodeOptionNetworkWide
)

func SendData(nodeId uint8, data []byte) []byte {
	buf := []byte{
		ZwSendData,
		byte(nodeId),
		byte(len(data)),
	}

	// @todo transport options
	buf = append(buf, data...)
	buf = append(buf, 0x1, 0x1)

	return buf
}

func ReadyCommand() []byte {
	return []byte{ZwSerialAPIReady, 0x01}
}

func GetNodeList() []byte {
	return []byte{ZwNodeList}
}

func EnterInclusionMode() []byte {
	return []byte{
		ZwAddNodeToNetwork,
		AddNodeAny | AddNodeOptionNormalPower | AddNodeOptionNetworkWide,
		0x01,
	}
}

func ExitInclusionMode() []byte {
	return []byte{
		ZwRemoveNodeFromNetwork,
		AddNodeStop,
		0x01,
	}
}

func EnterExclusionMode() []byte {
	return []byte{
		ZwRemoveNodeFromNetwork,
		RemoveNodeAny | RemoveNodeOptionNormalPower | RemoveNodeOptionNetworkWide,
		0x01,
	}
}

func ExitExclusionMode() []byte {
	return []byte{
		ZwRemoveNodeFromNetwork,
		RemoveNodeAny | RemoveNodeOptionNormalPower | RemoveNodeOptionNetworkWide,
		0x01,
	}
}
