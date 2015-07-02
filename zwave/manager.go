package zwave

import "fmt"
import cc "github.com/bjyoungblood/gozw/zwave/commandclass"

type Manager struct {
	session *SessionLayer

	ApiVersion     string
	ApiLibraryType string

	HomeId uint32
	NodeId byte

	Version                 byte
	ApiType                 string
	TimerFunctionsSupported bool
	IsPrimaryController     bool
	NodeList                []uint8
	ApplicationVersion      byte
	ApplicationRevision     byte
	SupportedFunctions      []byte
}

func NewManager(session *SessionLayer) *Manager {
	manager := &Manager{
		session: session,
	}

	manager.init()

	return manager
}

func (m *Manager) init() {
	version := m.GetVersion()
	m.ApiVersion = version.ApiVersion
	m.ApiLibraryType = version.GetLibraryTypeString()

	ids := m.GetHomeId()
	m.HomeId = ids.HomeId
	m.NodeId = ids.NodeId

	appInfo := m.GetAppInfo()
	m.Version = appInfo.Version
	m.ApiType = appInfo.GetApiType()
	m.TimerFunctionsSupported = appInfo.TimerFunctionsSupported()
	m.IsPrimaryController = appInfo.IsPrimaryController()
	m.NodeList = appInfo.GetNodeIds()

	serialApi := m.GetSerialApiCapabilities()
	m.ApplicationVersion = serialApi.ApplicationVersion
	m.ApplicationRevision = serialApi.ApplicationRevision
	m.SupportedFunctions = serialApi.GetSupportedFunctions()

	m.setSerialApiReady(true)
}

func (m *Manager) Close() {
	m.setSerialApiReady(false)
}

func (m *Manager) SetApplicationNodeInformation() {
	m.session.ExecuteCommandNoWait(FnApplicationNodeInformation, []byte{
		ApplicationNodeInfoListening | ApplicationFreqListeningMode250ms | ApplicationNodeInfoOptionalFunctionality,
		GenericTypeGenericController,
		SpecificTypePortableSceneController,
		0x01,
		cc.CommandClassBasic,
	})
}

func (m *Manager) FactoryReset() {
	m.session.ExecuteCommand(FnSetDefault, []byte{0x01})
}

func (m *Manager) AddNode() {
	m.session.ExecuteCommandNoWait(FnAddNodeToNetwork, []byte{
		AddNodeAny | AddNodeOptionNetworkWide | AddNodeOptionNormalPower,
		0x01,
	})

	for {
		frame := m.session.WaitForFrame()
		callback := ParseFunctionPayload(frame.Payload).(*AddRemoveNodeCallback)

		switch {
		case callback.Status == AddNodeStatusLearnReady:
			// timeout after waiting for this for 10 seconds
			// then perform ADD_NODE_STOP
			fmt.Print("Add node ready... ")
		case callback.Status == AddNodeStatusNodeFound:
			// recommended timeout interval for receiving this is 60 seconds
			// can either timeout or manually stop inclusion process
			fmt.Print("found node... ")
		case callback.Status == AddNodeStatusAddingSlave:
			// recommended timeout interval after NodeFound is 60 seconds
			// if timeout, must call ADD_NODE_STOP with callback and wait for ADD_NODE_STATUS_DONE
			fmt.Print("slave... ")
		case callback.Status == AddNodeStatusProtocolDone:
			// must call ADD_NODE_STOP and wait for AddNodeStatusDone
			// must timeout after period depending on network size and composition (see 4.4.1.3.3)
			// when timing out, must call ADD_NODE_STOP with callback
			fmt.Println("protocol done")
			m.session.ExecuteCommandNoWait(FnAddNodeToNetwork, []byte{
				AddNodeStop,
				0x02,
			})
		case callback.Status == AddNodeStatusDone:
			// Must call ADD_NODE_STOP with null callback (Serial API: is this 0x0 or omit the field?)
			fmt.Println("done")
			m.session.ExecuteCommandNoWait(FnAddNodeToNetwork, []byte{
				AddNodeStop,
				0x0,
			})

			return
		default:
			fmt.Println("unknown frame", callback)
		}

	}
}

func (m *Manager) RemoveNode() {
	m.session.ExecuteCommandNoWait(FnRemoveNodeFromNetwork, []byte{
		RemoveNodeAny | RemoveNodeOptionNetworkWide | RemoveNodeOptionNormalPower,
		0x01,
	})

	for {
		frame := m.session.WaitForFrame()
		callback := ParseFunctionPayload(frame.Payload).(*AddRemoveNodeCallback)

		switch {
		case callback.Status == RemoveNodeStatusLearnReady:
			fmt.Print("Remove node ready... ")
		case callback.Status == RemoveNodeStatusNodeFound:
			fmt.Print("found node... ")
		case callback.Status == RemoveNodeStatusRemovingSlave:
			fmt.Print("slave... ")
		case callback.Status == RemoveNodeStatusProtocolDone:
			fmt.Println("protocol done")
			m.session.ExecuteCommandNoWait(FnRemoveNodeFromNetwork, []byte{
				RemoveNodeStop,
				0x02,
			})
		case callback.Status == RemoveNodeStatusDone:
			fmt.Println("done")
			m.session.ExecuteCommandNoWait(FnRemoveNodeFromNetwork, []byte{
				RemoveNodeStop,
				0x0,
			})

			return
		default:
			fmt.Println("unknown frame", callback)
		}

	}
}

func (m *Manager) GetHomeId() *MemoryGetIdResponse {
	resp := m.session.ExecuteCommand(FnMemoryGetId, []byte{})
	return ParseFunctionPayload(resp.Payload).(*MemoryGetIdResponse)
}

func (m *Manager) GetAppInfo() *NodeListResponse {
	resp := m.session.ExecuteCommand(FnGetInitAppData, []byte{})
	return ParseFunctionPayload(resp.Payload).(*NodeListResponse)
}

func (m *Manager) GetSerialApiCapabilities() *SerialApiCapabilitiesResponse {
	resp := m.session.ExecuteCommand(FnSerialApiCapabilities, []byte{})
	return ParseFunctionPayload(resp.Payload).(*SerialApiCapabilitiesResponse)
}

func (m *Manager) GetVersion() *VersionResponse {
	resp := m.session.ExecuteCommand(FnGetVersion, []byte{})
	return ParseFunctionPayload(resp.Payload).(*VersionResponse)
}

func (m *Manager) GetNodeProtocolInfo(nodeId uint8) *NodeProtocolInfoResponse {
	resp := m.session.ExecuteCommand(FnGetNodeProtocolInfo, []byte{nodeId})
	return ParseFunctionPayload(resp.Payload).(*NodeProtocolInfoResponse)
}

func (m *Manager) SendData(nodeId uint8, data []byte) {
	payload := []byte{
		nodeId,
		uint8(len(data)),
	}

	payload = append(payload, data...)
	payload = append(payload, TransmitOptionAck|TransmitOptionAutoRoute)
	payload = append(payload, 0x01)

	resp := m.session.ExecuteCommand(FnSendData, payload)
	fmt.Println("senddata reply", resp.Payload)
}

func (m *Manager) setSerialApiReady(ready bool) {
	var rdy byte
	if ready {
		rdy = 1
	} else {
		rdy = 0
	}

	m.session.ExecuteCommandNoWait(FnSerialAPIReady, []byte{rdy})
}
