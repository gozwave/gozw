package protocol

const (
	FnSerialAPIGetInitAppData             byte = 0x02
	FnSerialAPIApplicationNodeInformation      = 0x03
	FnApplicationCommandHandler                = 0x04
	FnGetControllerCapabilities                = 0x05
	FnSerialAPIGetTimeouts                     = 0x06
	FnSerialAPIGetCapabilities                 = 0x07
	FnSerialAPISoftReset                       = 0x08
	FnGetProtocolVersion                       = 0x09
	FnSendNodeInformation                      = 0x12
	FnSendData                                 = 0x13
	FnSendDataMulti                            = 0x14
	FnGetVersion                               = 0x15
	FnSendDataAbort                            = 0x16
	FnRFPowerLevelSet                          = 0x17
	FnSendDataMeta                             = 0x18
	FnSetRoutingInfo                           = 0x1B
	FnRFPowerLevelRediscoverySet               = 0x1E
	FnMemoryGetID                              = 0x20
	FnGetNodeProtocolInfo                      = 0x41
	FnSetDefault                               = 0x42
	FnAssignReturnRoute                        = 0x46
	FnDeleteReturnRoute                        = 0x47
	FnRequestNodeNeighborUpdate                = 0x48
	FnApplicationControllerUpdate              = 0x49
	FnAddNodeToNetwork                         = 0x4a
	FnRemoveNodeFromNetwork                    = 0x4b
	FnRequestNetworkUpdate                     = 0x53
	FnRequestNodeInfo                          = 0x60
	FnRemoveFailingNode                        = 0x61
	FnIsNodeFailed                             = 0x62
	FnApplicationCommandHandlerBridge          = 0xA8
	FnSerialAPIReady                           = 0xEF
)

const (
	LibraryControllerStatic byte = 0x01
	LibraryController            = 0x02
	LibrarySlaveEnhanced         = 0x03
	LibrarySlave                 = 0x04
	LibraryInstaller             = 0x05
	LibrarySlaveRouting          = 0x06
	LibraryControllerBridge      = 0x07
	LibraryDUT                   = 0x08
	LibraryAvRemote              = 0x0A
	LibraryAvDevice              = 0x0B
)

const (
	UpdateStateNodeInfoReceived  byte = 0x84
	UpdateStateNodeInfoReqDone        = 0x82
	UpdateStateNodeInfoReqFailed      = 0x81
	UpdateStateRoutingPending         = 0x80
	UpdateStateNewIDAssigned          = 0x40
	UpdateStateDeleteDone             = 0x20
	UpdateStateSucID                  = 0x10
)

const (
	NodeOk                  byte = 0
	FailedNodeRemoved            = 1
	FailedNodeNotRemoved         = 2
	FailedNodeReplace            = 3
	FailedNodeReplaceDone        = 4
	FailedNodeReplaceFailed      = 5
)
