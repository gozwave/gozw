package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/comail/colog"
	"github.com/gozwave/gozw/application"
	"github.com/gozwave/gozw/frame"
	serialapi "github.com/gozwave/gozw/serial-api"
	"github.com/gozwave/gozw/session"
	"github.com/gozwave/gozw/transport"
	"github.com/gozwave/gozw/util"
	"go.bug.st/serial.v1/enumerator"
)

func init() {
	colog.Register()
	colog.ParseFields(true)
}

func main() {
	var port = flag.String("port", findPort(), "Device path (ex /dev/ttyACM0)")
	flag.Parse()

	if *port == "" {
		fmt.Println("No known device found, please provide the port with --port=")
		return
	}

	transport, err := transport.NewSerialPortTransport(*port, 115200)
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

	fmt.Printf("Home ID: 0x%x; Node ID: %d\n", appLayer.Controller.HomeID, appLayer.Controller.NodeID)
	fmt.Println("API Version:", appLayer.Controller.APIVersion)
	fmt.Println("Library:", appLayer.Controller.APILibraryType)
	fmt.Println("Version:", appLayer.Controller.Version)
	fmt.Println("API Type:", appLayer.Controller.APIType)
	fmt.Println("Is Primary Controller:", appLayer.Controller.IsPrimaryController)
	fmt.Println("Node count:", len(appLayer.Nodes()))
	fmt.Println("------------------------------------------------")

	<-time.After(time.Second * 10)

	for _, node := range appLayer.Nodes() {
		if node.NodeID == 1 {
			continue
		}

		fmt.Println(node.String())
	}

	select {}
}

func findPort() string {
	knownDevices := util.KnownUsbDevices()

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}

	for _, port := range ports {
		id := fmt.Sprintf("%s:%s", port.VID, port.PID)

		if dev, ok := knownDevices[id]; ok {
			fmt.Printf("Found port: %s (%s)\n", port.Name, dev)
			return port.Name
		}
	}

	return ""
}
