package zwave

import "fmt"

import "github.com/bjyoungblood/gozw/zwave/commandclass"

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
	ApplicationVersion      byte
	ApplicationRevision     byte
	SupportedFunctions      []byte

	nodeList []uint8
	Nodes    map[uint8]*Node
}

func NewManager(session *SessionLayer) *Manager {
	manager := &Manager{
		session: session,
		Nodes:   map[uint8]*Node{},
	}

	session.manager = manager

	manager.init()

	return manager
}

func (m *Manager) init() {
	version, err := m.GetVersion()
	if err != nil {
		panic(err)
	}

	m.ApiVersion = version.ApiVersion
	m.ApiLibraryType = version.GetLibraryTypeString()

	ids, err := m.GetHomeId()
	if err != nil {
		panic(err)
	}

	m.HomeId = ids.HomeId
	m.NodeId = ids.NodeId

	appInfo, err := m.GetAppInfo()
	if err != nil {
		panic(err)
	}

	m.Version = appInfo.Version
	m.ApiType = appInfo.GetApiType()
	m.TimerFunctionsSupported = appInfo.TimerFunctionsSupported()
	m.IsPrimaryController = appInfo.IsPrimaryController()
	m.nodeList = appInfo.GetNodeIds()

	serialApi, err := m.GetSerialApiCapabilities()
	if err != nil {
		panic(err)
	}

	m.ApplicationVersion = serialApi.ApplicationVersion
	m.ApplicationRevision = serialApi.ApplicationRevision
	m.SupportedFunctions = serialApi.GetSupportedFunctions()

	m.loadNodes()

	go m.handleUnsolicitedFrames()

	m.session.SetSerialAPIReady(true)
}

func (m *Manager) Close() {
	m.session.SetSerialAPIReady(false)
}

func (m *Manager) SetApplicationNodeInformation() {
	m.session.ApplicationNodeInformation(
		ApplicationNodeInfoListening|ApplicationFreqListeningMode250ms|ApplicationNodeInfoOptionalFunctionality,
		GenericTypeGenericController,
		SpecificTypePortableSceneController,
		[]uint8{
			commandclass.CommandClassBasic,
		},
	)
}

// @todo temporary, remove me
func (m *Manager) SendDataSecure(nodeId uint8, data []byte) error {
	return m.session.SendDataSecure(nodeId, data)
}

func (m *Manager) FactoryReset() {
	m.session.SetDefault()
}

func (m *Manager) AddNode() {
	node, err := m.session.AddNodeToNetwork()
	if err != nil {
		fmt.Println(err)
		return
	}

	m.Nodes[node.NodeId] = node
	fmt.Println(node.String())
}

func (m *Manager) includeSecureNode(node *Node) {
	// fmt.Println(node.NodeId)
	// m.SendData(node.NodeId, commandclass.NewSecuritySchemeGet())
	// schemeReport := m.session.WaitForFrame()
	// fmt.Println(schemeReport.Payload)
}

func (m *Manager) RemoveNode() {
	node, err := m.session.RemoveNodeFromNetwork()
	if err != nil {
		fmt.Println(err)
		return
	}

	delete(m.Nodes, node.NodeId)
	fmt.Println(node.String())
}

func (m *Manager) GetHomeId() (*MemoryGetIdResponse, error) {
	return m.session.MemoryGetId()
}

func (m *Manager) GetAppInfo() (*NodeListResponse, error) {
	return m.session.GetInitAppData()
}

func (m *Manager) GetSerialApiCapabilities() (*SerialApiCapabilitiesResponse, error) {
	return m.session.GetSerialApiCapabilities()
}

func (m *Manager) GetVersion() (*VersionResponse, error) {
	return m.session.GetVersion()
}

func (m *Manager) SendData(nodeId uint8, data []byte) {
	resp, err := m.session.SendData(nodeId, data, false)
	fmt.Println(resp, err)
	if err != nil {
		fmt.Println("senddata error", resp)
	} else {
		fmt.Println("senddata reply", resp.Payload)
	}
}

func (m *Manager) loadNodes() {
	for _, nodeId := range m.nodeList {
		nodeInfo, _ := m.session.GetNodeProtocolInfo(nodeId)
		node := NewNode(m, nodeId)

		node.setFromNodeProtocolInfo(nodeInfo)

		m.Nodes[nodeId] = node
	}
}

func (m *Manager) handleUnsolicitedFrames() {
	for frame := range m.session.UnsolicitedFrames {
		switch frame.Payload[0] {
		case FnApplicationCommandHandlerBridge:
			cmd := ParseApplicationCommandHandler(frame.Payload)
			if cmd.CmdLength > 0 {
				fmt.Printf("Got %s: %v\n", commandclass.GetCommandClassString(cmd.CommandData[0]), cmd.CommandData)
			} else {
				fmt.Println("wat", cmd)
			}
		default:
			fmt.Println("Received unsolicited frame:", frame)
		}
	}
}
