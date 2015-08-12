package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/helioslabs/gozw/zwave/application"
	"github.com/helioslabs/gozw/zwave/command-class"
	"github.com/helioslabs/gozw/zwave/frame"
	"github.com/helioslabs/gozw/zwave/serial-api"
	"github.com/helioslabs/gozw/zwave/session"
	"github.com/helioslabs/gozw/zwave/transport"
	"github.com/peterh/liner"
)

func main() {
	transport, err := transport.NewSerialPortTransport("/tmp/usbmodem", 115200)
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
		"(L) load all user codes for node",
		"(UN) request and print the number of supported user codes",
		"(UC) request a single user code",
		"(UCS) user code set",
		"(UCC) user code clear",
		"(ST) set temperature",
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

		case "UCS":
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
			userId, err := strconv.Atoi(input)
			if err != nil {
				spew.Dump(err)
				continue
			}

			code, _ := line.Prompt("code: ")
			if len(code) < 4 || len(code) > 8 {
				fmt.Println("Invalid code length")
				continue
			}

			lock.SetUserCode(byte(userId), []byte(code))

		case "UCC":
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
			userId, err := strconv.Atoi(input)
			if err != nil {
				spew.Dump(err)
				continue
			}

			lock.ClearUserCode(byte(userId))

		case "LS":
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

			spew.Dump(lock.GetLockStatus())

		case "ST":
			input, _ := line.Prompt("node id: ")
			nodeId, _ := strconv.Atoi(input)
			node, err := appLayer.Node(byte(nodeId))
			if err != nil {
				spew.Dump(err)
				continue
			}

			thermostat, err := node.GetThermostat()
			if err != nil {
				spew.Dump(err)
				continue
			}

			var setpointType commandclass.ThermostatSetpointType
			input, _ = line.Prompt("(c)ooling or (h)eating> ")
			switch input {
			case "c":
				setpointType = commandclass.ThermostatSetpointTypeCooling
			case "h":
				setpointType = commandclass.ThermostatSetpointTypeHeating
			default:
				fmt.Println("gg man")
				continue
			}

			input, _ = line.Prompt("temperature> ")
			temperature, err := strconv.Atoi(input)
			if err != nil {
				spew.Dump(err)
				continue
			}

			thermostat.SetpointSet(setpointType, float64(temperature))

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
			fmt.Printf("Home ID: 0x%x; Node ID: %d\n", appLayer.HomeID, appLayer.NodeID)
			fmt.Println("API Version:", appLayer.APIVersion)
			fmt.Println("Library:", appLayer.APILibraryType)
			fmt.Println("Version:", appLayer.Version)
			fmt.Println("API Type:", appLayer.APIType)
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
