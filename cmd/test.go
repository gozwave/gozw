package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/comail/colog"
	"github.com/davecgh/go-spew/spew"
	"github.com/gozwave/gozw/application"
	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/session"
	"github.com/gozwave/gozw/transport"
	"github.com/peterh/liner"
)

func init() {
	colog.Register()
	colog.ParseFields(true)
}

func main() {
	transport, err := transport.NewSerialPortTransport("/dev/tty.usbmodem1461", 115200)
	if err != nil {
		panic(err)
	}

	frameLayer := frame.NewFrameLayer(transport)
	sessionLayer := session.NewSessionLayer(frameLayer)
	apiLayer := serialapi.NewLayer(sessionLayer)
	appLayer, err := application.NewLayer(apiLayer)
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
		"(M) load manufacturer-specific data for node",
		"(PV) print the result of the above",
		"(NIF) request node information frame from node",
		"(F)ailed node removal",
		"(p)rint network info",
		"(ON) turn light on",
		"(ON) turn light off",
		"(LOCK) lock door",
		"(UNLOCK) lock door",
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
		case "M":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			spew.Dump(node.LoadManufacturerInfo())

		case "PV":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			for id, cc := range node.CommandClasses {
				fmt.Printf(
					"%s: %d\n",
					id,
					cc.Version,
				)
			}

		case "ON":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			spew.Dump(node.SendCommand(&switchbinary.Set{
				SwitchValue: 0x01,
			}))
		case "OFF":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			spew.Dump(node.SendCommand(&switchbinary.Set{
				SwitchValue: 0x00,
			}))

		case "UNLOCK":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			spew.Dump(node.SendCommand(&doorlock.OperationSet{
				DoorLockMode: 0x00,
			}))

		case "LOCK":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			spew.Dump(node.SendCommand(&doorlock.OperationSet{
				DoorLockMode: 0xFF,
			}))

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
			fmt.Printf("Home ID: 0x%x; Node ID: %d\n", appLayer.Controller.HomeID, appLayer.Controller.NodeID)
			fmt.Println("API Version:", appLayer.Controller.APIVersion)
			fmt.Println("Library:", appLayer.Controller.APILibraryType)
			fmt.Println("Version:", appLayer.Controller.Version)
			fmt.Println("API Type:", appLayer.Controller.APIType)
			fmt.Println("Is Primary Controller:", appLayer.Controller.IsPrimaryController)
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
