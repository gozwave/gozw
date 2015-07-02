package main

import (
	"fmt"
	"strconv"

	"github.com/bjyoungblood/gozw/zwave"
	"github.com/bjyoungblood/gozw/zwave/commandclass"
	"github.com/peterh/liner"
)

func main() {

	transport, err := zwave.NewTransportLayer("/tmp/usbmodem", 115200)
	if err != nil {
		panic(err)
	}

	frameLayer := zwave.NewFrameLayer(transport)
	sessionLayer := zwave.NewSessionLayer(frameLayer)
	manager := zwave.NewManager(sessionLayer)

	defer manager.Close()

	fmt.Printf("Home ID: 0x%x; Node ID: %d\n", manager.HomeId, manager.NodeId)
	fmt.Println("API Version:", manager.ApiVersion)
	fmt.Println("Library:", manager.ApiLibraryType)
	fmt.Println("Version:", manager.Version)
	fmt.Println("API Type:", manager.ApiType)
	fmt.Println("Timer Functions Supported:", manager.TimerFunctionsSupported)
	fmt.Println("Is Primary Controller:", manager.IsPrimaryController)
	fmt.Println("Nodes:", manager.NodeList)

	// manager.SetApplicationNodeInformation()
	// manager.FactoryReset()

	for _, i := range manager.NodeList {
		nodeInfo := manager.GetNodeProtocolInfo(i)

		fmt.Printf("Node %d: \n", i)
		fmt.Printf("  Is listening? %t\n", nodeInfo.IsListening())
		fmt.Printf("  Basic device class: %s\n", nodeInfo.GetBasicDeviceClassName())
		fmt.Printf("  Generic device class: %s\n", nodeInfo.GetGenericDeviceClassName())
		fmt.Printf("  Specific device class: %s\n", nodeInfo.GetSpecificDeviceClassName())
		fmt.Printf("  Raw: %v\n\n", nodeInfo)
	}

	// manager.SendData(3, cc.NewSwitchMultilevelCommand(0))

	line := liner.NewLiner()
	defer line.Close()

	for {
		cmd, _ := line.Prompt("(a)dd node\n(r)emove node\n(g)et nonce\n(q)uit\n> ")
		switch cmd {
		case "a":
			manager.AddNode()
		case "r":
			manager.RemoveNode()
		case "s":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			manager.SendData(uint8(nodeId), commandclass.NewSecuritySchemeGet())
		case "g":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			manager.SendData(uint8(nodeId), commandclass.NewSecurityNonceGet())
		case "q":
			return
		default:
			fmt.Println("invalid selection")
		}
	}

}
