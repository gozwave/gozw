package protocol

const (
	AddNodeAny        uint8 = 1
	AddNodeController       = 2
	AddNodeSlave            = 3
	AddNodeExisting         = 4
	AddNodeStop             = 5
	AddNodeStopFailed       = 6
)

const (
	AddNodeOptionNormalPower uint8 = 0x80
	AddNodeOptionNetworkWide       = 0x40
)

const (
	AddNodeStatusLearnReady       uint8 = 1
	AddNodeStatusNodeFound              = 2
	AddNodeStatusAddingSlave            = 3
	AddNodeStatusAddingController       = 4
	AddNodeStatusProtocolDone           = 5
	AddNodeStatusDone                   = 6
	AddNodeStatusFailed                 = 7
	AddNodeStatusSecurityFailed         = 9
)

const (
	RemoveNodeAny        uint8 = AddNodeAny
	RemoveNodeController       = AddNodeController
	RemoveNodeSlave            = AddNodeSlave
	RemoveNodeStop             = AddNodeStop
)

const (
	RemoveNodeOptionNormalPower uint8 = AddNodeOptionNormalPower
	RemoveNodeOptionNetworkWide       = AddNodeOptionNetworkWide
)

const (
	RemoveNodeStatusLearnReady         uint8 = AddNodeStatusLearnReady
	RemoveNodeStatusNodeFound                = AddNodeStatusNodeFound
	RemoveNodeStatusRemovingSlave            = AddNodeStatusAddingSlave
	RemoveNodeStatusRemovingController       = AddNodeStatusAddingController
	RemoveNodeStatusProtocolDone             = AddNodeStatusProtocolDone
	RemoveNodeStatusDone                     = AddNodeStatusDone
	RemoveNodeStatusFailed                   = AddNodeStatusFailed
)
