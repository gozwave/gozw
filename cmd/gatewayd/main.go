package main

import (
	"fmt"

	"github.com/bjyoungblood/gozw/common"
	"github.com/bjyoungblood/gozw/gateway"
	"github.com/bjyoungblood/gozw/zwave"
)

func main() {

	config, err := common.LoadGatewayConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	serialPort, err := gateway.NewSerialPort(config)
	if err != nil {
		panic(err)
	}

	serialPort.Initialize()
	go serialPort.Run()

	go func() {
		version := serialPort.SendFrameSync(zwave.NewRequestFrame([]byte{0x15}))
		fmt.Println(version)
	}()

	for {
		// frame := <-serialPort.Incoming
		// packet := common.WirePacket{frame.Marshal()}
		fmt.Println("MAIN:", <-serialPort.Incoming)
	}

	// defer serialPort.Close()

}
