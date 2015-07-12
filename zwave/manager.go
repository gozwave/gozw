package zwave

import (
	"fmt"
	"sync"
)

import "github.com/bjyoungblood/gozw/zwave/commandclass"

type Manager struct {
	session SessionLayer

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

func NewManager(session SessionLayer) *Manager {
	manager := &Manager{
		session: session,
		Nodes:   map[uint8]*Node{},
	}

	session.SetManager(manager)

	manager.init()

	return manager
}

func (m *Manager) init() {
	version, err := m.session.GetVersion()
	if err != nil {
		panic(err)
	}

	m.ApiVersion = version.ApiVersion
	m.ApiLibraryType = version.GetLibraryTypeString()

	ids, err := m.session.MemoryGetId()
	if err != nil {
		panic(err)
	}

	m.HomeId = ids.HomeId
	m.NodeId = ids.NodeId

	appInfo, err := m.session.GetInitAppData()
	if err != nil {
		panic(err)
	}

	m.Version = appInfo.Version
	m.ApiType = appInfo.GetApiType()
	m.TimerFunctionsSupported = appInfo.TimerFunctionsSupported()
	m.IsPrimaryController = appInfo.IsPrimaryController()
	m.nodeList = appInfo.GetNodeIds()

	serialApi, err := m.session.GetSerialApiCapabilities()
	if err != nil {
		panic(err)
	}

	m.ApplicationVersion = serialApi.ApplicationVersion
	m.ApplicationRevision = serialApi.ApplicationRevision
	m.SupportedFunctions = serialApi.GetSupportedFunctions()

	go m.handleUnsolicitedFrames()

	m.session.registerApplicationCommandHandler(
		commandclass.CommandClassSecurity,
		m.handleSecurityCommands,
	)

	m.loadNodes()

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

func (m *Manager) RemoveNode() {
	node, err := m.session.RemoveNodeFromNetwork()
	if err != nil {
		fmt.Println(err)
		return
	}

	delete(m.Nodes, node.NodeId)
	fmt.Println(node.String())
}

func (m *Manager) SendData(nodeId uint8, data []byte) {
	resp, err := m.session.SendData(nodeId, data)
	fmt.Println(resp, err)
	if err != nil {
		fmt.Println("senddata error", resp)
	} else {
		fmt.Println("senddata reply", resp.Payload)
	}
}

func (m *Manager) removeFailedNode(nodeId uint8) {
	result, err := m.session.removeFailedNode(nodeId)
	if err != nil {
		fmt.Println("error removing failed node: ", err)
	} else {
		fmt.Println("remove failed node result:", result.Payload)
	}
}

func (m *Manager) loadNodes() {
	var wg sync.WaitGroup

	for _, nodeId := range m.nodeList {
		node := NewNode(m, nodeId)
		m.Nodes[nodeId] = node

		wait := node.Initialize()

		wg.Add(1)

		go func() {
			<-wait
			fmt.Println("init", node.NodeId)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (m *Manager) handleSecurityCommands(cmd *ApplicationCommandHandler, frame *Frame) {
	switch cmd.CommandData[1] {
	case commandclass.CommandSecurityCommandsSupportedReport:
		cc := commandclass.ParseSecurityCommandsSupportedReport(cmd.CommandData)
		if node, ok := m.Nodes[cmd.SrcNodeId]; ok {
			node.receiveSecurityCommandsSupportedReport(cc)
		}
		fmt.Println(cc.SupportedCommandClasses)

	default:
		fmt.Println(cmd)
	}
}

func (m *Manager) handleUnsolicitedFrames() {
	frames := m.session.GetUnsolicitedFrames()
	for frame := range frames {
		switch frame.Payload[0] {

		case FnApplicationCommandHandlerBridge, FnApplicationCommandHandler:
			cmd := ParseApplicationCommandHandler(frame.Payload)
			if cmd.CmdLength > 0 {
				fmt.Printf("Got %s: %v\n", commandclass.GetCommandClassString(cmd.CommandData[0]), cmd.CommandData)
			} else {
				fmt.Println("wat", cmd)
			}

		case FnApplicationControllerUpdate:
			cmd := ParseApplicationControllerUpdate(frame.Payload)
			if cmd.Status != UpdateStateNodeInfoReceived {
				fmt.Printf("Node update (%d): %s\n", cmd.NodeId, cmd.GetStatusString())
			} else {
				if node, ok := m.Nodes[cmd.NodeId]; ok {
					node.setFromApplicationControllerUpdate(cmd)
				}
			}

		default:
			fmt.Println("Received unsolicited frame:", frame)

		}
	}
}
