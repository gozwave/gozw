package main

import (
	"fmt"
	"strconv"

	"github.com/bjyoungblood/gozw/zwave/application"
	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
	"github.com/bjyoungblood/gozw/zwave/session"
	"github.com/bjyoungblood/gozw/zwave/transport"
	"github.com/davecgh/go-spew/spew"
	"github.com/peterh/liner"
)

func main() {
	transport, err := transport.NewSerialTransportLayer("/tmp/usbmodem", 115200)
	if err != nil {
		panic(err)
	}

	frameLayer := frame.NewFrameLayer(transport)
	sessionLayer := session.NewSessionLayer(frameLayer)
	apiLayer := serialapi.NewSerialAPILayer(sessionLayer)
	appLayer, err := application.NewApplicationLayer(apiLayer)
	if err != nil {
		panic(err)
	}

	// spew.Dump(applicationLayer.Nodes)

	// apiLayer.SoftReset()

	// spew.Dump(apiLayer.GetVersion())
	// nodeList, err := apiLayer.GetNodeList()
	// fmt.Println(nodeList.GetNodeIds())
	// spew.Dump(apiLayer.MemoryGetId())
	// spew.Dump(apiLayer.GetNodeProtocolInfo(27))
	// spew.Dump(apiLayer.GetSerialApiCapabilities())

	// txTime, err := apiLayer.SendData(27, commandclass.NewVersionGet())
	// fmt.Println("TX: ", txTime)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if apiLayer.AddNode() != nil {
	// 	apiLayer.AddNode()
	// }

	// sessionLayer := zwave.NewSessionLayer(frameLayer)
	// manager := zwave.NewManager(sessionLayer)

	// defer manager.Close()

	fmt.Printf("Home ID: 0x%x; Node ID: %d\n", appLayer.HomeId, appLayer.NodeId)
	fmt.Println("API Version:", appLayer.ApiVersion)
	fmt.Println("Library:", appLayer.ApiLibraryType)
	fmt.Println("Version:", appLayer.Version)
	fmt.Println("API Type:", appLayer.ApiType)
	fmt.Println("Is Primary Controller:", appLayer.IsPrimaryController)
	fmt.Println("Node count:", len(appLayer.Nodes()))
	//
	// appLayer.SendDataSecure(42, []byte{
	// 	commandclass.CommandClassDoorLock,
	// 	0x01, // door lock operation set
	// 	0xFF, // unsecured
	// })

	// manager.SetApplicationNodeInformation()
	// manager.FactoryReset()

	for _, node := range appLayer.Nodes() {
		fmt.Println(node.String())
	}

	// manager.SendData(3, cc.NewSwitchMultilevelCommand(0))

	line := liner.NewLiner()
	defer line.Close()

	for {
		cmd, _ := line.Prompt("(a)dd node\n(r)emove node\n(g)et nonce\n(q)uit\n> ")
		switch cmd {
		case "a":
			spew.Dump(appLayer.AddNode())
		case "r":
			spew.Dump(appLayer.RemoveNode())
		case "L":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			spew.Dump(node.LoadAllUserCodes())
		case "F":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			spew.Dump(appLayer.RemoveFailedNode(byte(nodeId)))
		// case "s":
		// 	input, _ := line.Prompt("node id: ")
		// 	nodeId, _ := strconv.Atoi(input)
		// 	manager.SendData(uint8(nodeId), commandclass.NewSecuritySchemeGet())
		// case "g":
		// 	input, _ := line.Prompt("node id: ")
		// 	nodeId, _ := strconv.Atoi(input)
		// 	manager.SendData(uint8(nodeId), commandclass.NewSecurityNonceGet())
		// case "v":
		// 	input, _ := line.Prompt("node id: ")
		// 	nodeId, _ := strconv.Atoi(input)
		// 	manager.SendData(uint8(nodeId), commandclass.NewVersionGet())
		case "q":
			appLayer.Shutdown()
			return
		default:
			fmt.Println("invalid selection")
		}
	}

}
