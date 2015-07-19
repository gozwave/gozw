package main

import (
	"fmt"
	"strconv"

	"github.com/bjyoungblood/gozw/zwave"
	"github.com/bjyoungblood/gozw/zwave/commandclass"
	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/transport"
	"github.com/peterh/liner"
)

func main() {
	transport, err := transport.NewSerialTransportLayer("/tmp/usbmodem", 115200)
	if err != nil {
		panic(err)
	}

	frameLayer := frame.NewFrameLayer(transport)
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
	fmt.Println("Node count:", len(manager.Nodes))

	manager.SendDataSecure(15, []byte{
		commandclass.CommandClassDoorLock,
		0x01, // door lock operation set
		0x00, // unsecured
	})

	// manager.SetApplicationNodeInformation()
	// manager.FactoryReset()

	for _, node := range manager.Nodes {
		fmt.Println(node.String())
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
		case "v":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			manager.SendData(uint8(nodeId), commandclass.NewVersionGet())
		case "q":
			return
		default:
			fmt.Println("invalid selection")
		}
	}

}
