package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bjyoungblood/gozw/zwave"
	cc "github.com/bjyoungblood/gozw/zwave/commandclass"
)

func main() {

	transport, err := zwave.NewTransportLayer("/tmp/usbmodem", 115200)
	if err != nil {
		panic(err)
	}

	frameLayer := zwave.NewFrameLayer(transport)
	sessionLayer := zwave.NewSessionLayer(frameLayer)
	manager := zwave.NewManager(sessionLayer)

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

	manager.SendData(3, cc.NewSwitchMultilevelCommand(0))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	// manager.RemoveNode()
	// manager.AddNode()

}
