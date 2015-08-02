package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bjyoungblood/gozw/zwave/application"
	"github.com/bjyoungblood/gozw/zwave/command-class"
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

	defer appLayer.Shutdown()

	line := liner.NewLiner()
	defer line.Close()

	commands := strings.Join([]string{
		"(a)dd node",
		"(r)emove node",
		"(V) load command class versions for node",
		"(PV) print the result of the above",
		"(L) load all user codes for node",
		"(UN) request and print the number of supported user codes",
		"(UC) request a single user code",
		"(NIF) request node information frame from node",
		"(F)ailed node removal",
		"(p)rint network info",
		"(q)uit",
	}, "\n")

	fmt.Println(commands)

	for {
		cmd, _ := line.Prompt("> ")
		switch cmd {
		case "a":
			spew.Dump(appLayer.AddNode())
		case "r":
			spew.Dump(appLayer.RemoveNode())
		case "V":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			spew.Dump(node.LoadCommandClassVersions())
		case "PV":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			for cc, _ := range node.SupportedCommandClasses {
				fmt.Printf(
					"%s: %d\n",
					commandclass.GetCommandClassString(cc),
					node.CommandClassVersions[cc],
				)
			}

			for cc, _ := range node.SecureSupportedCommandClasses {
				fmt.Printf(
					"%s: %d\n",
					commandclass.GetCommandClassString(cc),
					node.CommandClassVersions[cc],
				)
			}

		case "L":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			lock, err := node.GetDoorLock()
			if err != nil {
				spew.Dump(err)
				continue
			}

			lock.LoadAllUserCodes()
		case "UN":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			lock, err := node.GetDoorLock()
			if err != nil {
				spew.Dump(err)
				continue
			}

			count, err := lock.GetSupportedUserCount()
			if err != nil {
				spew.Dump(err)
				continue
			}

			fmt.Printf("Supported users: %d\n", count)
		case "UC":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			lock, err := node.GetDoorLock()
			if err != nil {
				spew.Dump(err)
				continue
			}

			input, _ = line.Prompt("user id: ")
			userId, _ := strconv.Atoi(input)

			lock.LoadUserCode(byte(userId))
		case "NIF":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, _ := appLayer.Node(byte(nodeId))
			spew.Dump(node.RequestNodeInformationFrame())
		case "F":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			spew.Dump(appLayer.RemoveFailedNode(byte(nodeId)))
		case "p":
			fmt.Printf("Home ID: 0x%x; Node ID: %d\n", appLayer.HomeId, appLayer.NodeId)
			fmt.Println("API Version:", appLayer.ApiVersion)
			fmt.Println("Library:", appLayer.ApiLibraryType)
			fmt.Println("Version:", appLayer.Version)
			fmt.Println("API Type:", appLayer.ApiType)
			fmt.Println("Is Primary Controller:", appLayer.IsPrimaryController)
			fmt.Println("Node count:", len(appLayer.Nodes()))

			for _, node := range appLayer.Nodes() {
				fmt.Println(node.String())
			}
		case "q":
			return
		default:
			fmt.Printf("invalid selection\n\n")
			fmt.Println(commands)
		}
	}

}
