package zwave

const (
	ApplicationNodeInfoNotListening          = 0x00
	ApplicationNodeInfoListening             = 0x01
	ApplicationNodeInfoOptionalFunctionality = 0x02
	ApplicationFreqListeningMode1000ms       = 0x10
	ApplicationFreqListeningMode250ms        = 0x20
)

const (
	UpdateStateNodeInfoReceived  = 0x84
	UpdateStateNodeInfoReqDone   = 0x82
	UpdateStateNodeInfoReqFailed = 0x81
	UpdateStateRoutingPending    = 0x80
	UpdateStateNewIdAssigned     = 0x40
	UpdateStateDeleteDone        = 0x20
	UpdateStateSucId             = 0x10
)

const (
	AddNodeAny        = 1
	AddNodeController = 2
	AddNodeSlave      = 3
	AddNodeExisting   = 4
	AddNodeStop       = 5
	AddNodeStopFailed = 6
)

const (
	AddNodeOptionNormalPower = 0x80
	AddNodeOptionNetworkWide = 0x40
)

const (
	AddNodeStatusLearnReady       = 1
	AddNodeStatusNodeFound        = 2
	AddNodeStatusAddingSlave      = 3
	AddNodeStatusAddingController = 4
	AddNodeStatusProtocolDone     = 5
	AddNodeStatusDone             = 6
	AddNodeStatusFailed           = 7
	AddNodeStatusSecurityFailed   = 9
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

const (
	RemoveNodeStatusLearnReady         = AddNodeStatusLearnReady
	RemoveNodeStatusNodeFound          = AddNodeStatusNodeFound
	RemoveNodeStatusRemovingSlave      = AddNodeStatusAddingSlave
	RemoveNodeStatusRemovingController = AddNodeStatusAddingController
	RemoveNodeStatusProtocolDone       = AddNodeStatusProtocolDone
	RemoveNodeStatusDone               = AddNodeStatusDone
	RemoveNodeStatusFailed             = AddNodeStatusFailed
)

const (
	LibraryControllerStatic = 0x01
	LibraryController       = 0x02
	LibrarySlaveEnhanced    = 0x03
	LibrarySlave            = 0x04
	LibraryInstaller        = 0x05
	LibrarySlaveRouting     = 0x06
	LibraryControllerBridge = 0x07
	LibraryDUT              = 0x08
	LibraryAvRemote         = 0x0A
	LibraryAvDevice         = 0x0B
)
