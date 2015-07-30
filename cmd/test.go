package main

import (
	"fmt"
	"time"

	"github.com/bjyoungblood/gozw/zwave/commandclass"
	"github.com/bjyoungblood/gozw/zwave/frame"
	"github.com/bjyoungblood/gozw/zwave/serial-api"
	"github.com/bjyoungblood/gozw/zwave/session"
	"github.com/bjyoungblood/gozw/zwave/transport"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	transport, err := transport.NewSerialTransportLayer("/tmp/usbmodem", 115200)
	if err != nil {
		panic(err)
	}

	frameLayer := frame.NewFrameLayer(transport)
	sessionLayer := session.NewSessionLayer(frameLayer)
	apiLayer := serialapi.NewSerialAPILayer(sessionLayer)

	// apiLayer.SoftReset()

	// spew.Dump(apiLayer.GetVersion())
	nodeList, err := apiLayer.GetNodeList()
	fmt.Println(nodeList.GetNodeIds())
	// spew.Dump(apiLayer.MemoryGetId())
	spew.Dump(apiLayer.GetNodeProtocolInfo(27))
	// spew.Dump(apiLayer.GetSerialApiCapabilities())

	time.Sleep(time.Second * 1)

	txTime, err := apiLayer.SendData(27, commandclass.NewVersionGet())
	fmt.Println("TX: ", txTime)
	if err != nil {
		fmt.Println(err)
	}

	<-time.After(30 * time.Second)

	// if apiLayer.AddNode() != nil {
	// 	apiLayer.AddNode()
	// }

	// sessionLayer := zwave.NewSessionLayer(frameLayer)
	// manager := zwave.NewManager(sessionLayer)

	// defer manager.Close()

	// fmt.Printf("Home ID: 0x%x; Node ID: %d\n", manager.HomeId, manager.NodeId)
	// fmt.Println("API Version:", manager.ApiVersion)
	// fmt.Println("Library:", manager.ApiLibraryType)
	// fmt.Println("Version:", manager.Version)
	// fmt.Println("API Type:", manager.ApiType)
	// fmt.Println("Timer Functions Supported:", manager.TimerFunctionsSupported)
	// fmt.Println("Is Primary Controller:", manager.IsPrimaryController)
	// fmt.Println("Node count:", len(manager.Nodes))
	//
	// manager.SendDataSecure(15, []byte{
	// 	commandclass.CommandClassDoorLock,
	// 	0x01, // door lock operation set
	// 	0x00, // unsecured
	// })

	// manager.SetApplicationNodeInformation()
	// manager.FactoryReset()

	// for _, node := range manager.Nodes {
	// 	fmt.Println(node.String())
	// }

	// manager.SendData(3, cc.NewSwitchMultilevelCommand(0))

	// line := liner.NewLiner()
	// defer line.Close()
	//
	// for {
	// 	cmd, _ := line.Prompt("(a)dd node\n(r)emove node\n(g)et nonce\n(q)uit\n> ")
	// 	switch cmd {
	// 	case "a":
	// 		manager.AddNode()
	// 	case "r":
	// 		manager.RemoveNode()
	// 	case "s":
	// 		input, _ := line.Prompt("node id: ")
	// 		nodeId, _ := strconv.Atoi(input)
	// 		manager.SendData(uint8(nodeId), commandclass.NewSecuritySchemeGet())
	// 	case "g":
	// 		input, _ := line.Prompt("node id: ")
	// 		nodeId, _ := strconv.Atoi(input)
	// 		manager.SendData(uint8(nodeId), commandclass.NewSecurityNonceGet())
	// 	case "v":
	// 		input, _ := line.Prompt("node id: ")
	// 		nodeId, _ := strconv.Atoi(input)
	// 		manager.SendData(uint8(nodeId), commandclass.NewVersionGet())
	// 	case "q":
	// 		return
	// 	default:
	// 		fmt.Println("invalid selection")
	// 	}
	// }

}
