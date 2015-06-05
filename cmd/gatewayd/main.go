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

	// MAIN: &{1 26 0 [74 1 3 0 19 4 8 6 64 66 67 68 69 128 112 49 143 134 114 133 44 43 115 129] 62}
	// funcid = 1
	// status = 3
	// source = 0
	// blen = 19
	//

	// serialPort.SendFrame(zwave.NewRequestFrame(zwave.GetNodeList()), func(status int) {
	// 	fmt.Println("STATUS:", status)
	// })

	f := zwave.NewRequestFrame(zwave.SendData(uint8(29), zwave.NewThermostatSetpointCommand()))
	go serialPort.SendFrame(f, func(status int) {
		fmt.Println("STATUS:", status)
	})

	// go serialPort.SendFrame(zwave.NewRequestFrame(zwave.EnterInclusionMode()), func(status int) {
	// 	if status == int(zwave.FrameHeaderAck) {
	// 		fmt.Println("Entered inclusion mode")
	// 	} else {
	// 		fmt.Println("Could not enter inclusion mode")
	// 	}
	// })

	// go serialPort.SendFrame(zwave.NewRequestFrame(zwave.ExitInclusionMode()), func(status int) {
	// 	if status == int(zwave.FrameHeaderAck) {
	// 		fmt.Println("Exited inclusion mode")
	// 	} else {
	// 		fmt.Println("Could not exit inclusion mode")
	// 	}
	// })

	// go serialPort.SendFrame(zwave.NewRequestFrame(zwave.EnterExclusionMode()), func(status int) {
	// 	if status == int(zwave.FrameHeaderAck) {
	// 		fmt.Println("Entered exclusion mode")
	// 	} else {
	// 		fmt.Println("Could not enter exclusion mode")
	// 	}
	// })

	// go serialPort.SendFrame(zwave.NewRequestFrame(zwave.ExitExclusionMode()), func(status int) {
	// 	if status == int(zwave.FrameHeaderAck) {
	// 		fmt.Println("Exited exclusion mode")
	// 	} else {
	// 		fmt.Println("Could not exit exclusion mode")
	// 	}
	// })

	for {
		frame := <-serialPort.Incoming
		// packet := common.WirePacket{frame.Marshal()}
		fmt.Println("MAIN:", frame)

		if frame.Payload[0] == 0x02 && frame.IsResponse() {
			fmt.Println("Frame payload contains a node list")
			nodeList := zwave.NodeList{}
			nodeList.Unmarshal(frame)
			fmt.Println(nodeList.GetNodeIds())
		}

		if frame.Payload[0] == 0x4a && frame.Payload[2] == 5 {
			fmt.Println("NODE ADDED")
			go serialPort.SendFrame(zwave.NewRequestFrame(zwave.ExitInclusionMode()), func(status int) {
				if status == int(zwave.FrameHeaderAck) {
					fmt.Println("Exited inclusion mode")
				} else {
					fmt.Println("Could not exit inclusion mode")
				}
			})
		}

		if frame.Payload[0] == 0x4b && frame.Payload[2] == 5 {
			fmt.Println("NODE REMOVED")
			go serialPort.SendFrame(zwave.NewRequestFrame(zwave.ExitExclusionMode()), func(status int) {
				if status == int(zwave.FrameHeaderAck) {
					fmt.Println("Exited exclusion mode")
				} else {
					fmt.Println("Could not exit exclusion mode")
				}
			})
		}
	}

	// defer serialPort.Close()

}
