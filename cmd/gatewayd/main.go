package main

import (
	"github.com/bjyoungblood/gozw/common"
	"github.com/bjyoungblood/gozw/gateway"
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

	serialPort.Write([]byte{
		0x01,
		0x03,
		0x00,
		0x02,
		0xfe,
	})

	serialPort.Initialize()
	// go serialPort.Run()

	// MAIN: &{1 26 0 [74 1 3 0 19 4 8 6 64 66 67 68 69 128 112 49 143 134 114 133 44 43 115 129] 62}
	// funcid = 1
	// status = 3
	// source = 0
	// blen = 19

	// serialPort.SendFrame(zwave.NewRequestFrame(zwave.RequestNodeInfo(30)), func(status uint8, response *zwave.ZFrame) {
	// 	fmt.Println("STATUS:", status)
	// })

	// serialPort.SendFrame(zwave.NewRequestFrame(zwave.NodeProtocolInfo(30)), func(status uint8, response *zwave.ZFrame) {
	// 	fmt.Println("STATUS:", status)
	// })

	// getInitData := functions.NewGetInitData()
	// serialPort.SendFrame(zwave.NewRequestFrame(getInitData.Marshal()), func(status uint8, response *zwave.ZFrame) {
	// 	nodeList := zwave.NodeList{}
	// 	nodeList.Unmarshal(response)
	// 	fmt.Println("node list:", nodeList.GetNodeIds())
	// })

	// for i := 2; i <= 32; i++ {
	// i := 4

	// nodeInfoFrame := functions.NewRequestNodeInfo(uint8(i))
	// res := serialPort.SendFrameSync(zwave.NewRequestFrame(nodeInfoFrame.Marshal()))
	// fmt.Println(res)
	//
	// isFailedFrame := functions.NewRemoveFailedNode(uint8(i))
	// isFailed := serialPort.SendFrameSync(zwave.NewRequestFrame(isFailedFrame.Marshal()))
	// fmt.Println("Node id", i, isFailed)
	// }

	// f := zwave.NewRequestFrame(zwave.SendData(uint8(30), zwave.NewDoorLockCommand()))
	// fmt.Println(f)
	// go serialPort.SendFrame(f, func(status uint8, response *zwave.ZFrame) {
	// 	fmt.Println("STATUS:", status)
	// })

	// f := zwave.NewRequestFrame(zwave.SendData(uint8(29), zwave.NewThermostatSetpointCommand()))
	// go serialPort.SendFrame(f, func(status uint8, response *zwave.ZFrame) {
	// 	fmt.Println("STATUS:", status)
	// })

	// go serialPort.SendFrame(zwave.NewRequestFrame(zwave.EnterInclusionMode()), func(status uint8, response *zwave.ZFrame) {
	// 	if status == int(zwave.FrameHeaderAck) {
	// 		fmt.Println("Entered inclusion mode")
	// 	} else {
	// 		fmt.Println("Could not enter inclusion mode")
	// 	}
	// })

	// go serialPort.SendFrame(zwave.NewRequestFrame(zwave.ExitInclusionMode()), func(status uint8, response *zwave.ZFrame) {
	// 	if status == int(zwave.FrameHeaderAck) {
	// 		fmt.Println("Exited inclusion mode")
	// 	} else {
	// 		fmt.Println("Could not exit inclusion mode")
	// 	}
	// })

	// go serialPort.SendFrame(zwave.NewRequestFrame(zwave.EnterExclusionMode()), func(status uint8, response *zwave.ZFrame) {
	// 	if status == int(zwave.FrameHeaderAck) {
	// 		fmt.Println("Entered exclusion mode")
	// 	} else {
	// 		fmt.Println("Could not enter exclusion mode")
	// 	}
	// })

	// go serialPort.SendFrame(zwave.NewRequestFrame(zwave.ExitExclusionMode()), func(status uint8, response *zwave.ZFrame) {
	// 	if status == int(zwave.FrameHeaderAck) {
	// 		fmt.Println("Exited exclusion mode")
	// 	} else {
	// 		fmt.Println("Could not exit exclusion mode")
	// 	}
	// })

	// for {
	// 	frame := <-serialPort.Incoming
	// 	// packet := common.WirePacket{frame.Marshal()}
	// 	fmt.Println("MAIN:", frame)
	//
	// 	if frame.Payload[0] == 0x4a && frame.Payload[2] == 5 {
	// 		fmt.Println("NODE ADDED")
	// 		addNodeEnd := functions.NewAddNode()
	// 		frame := addNodeEnd.Marshal()
	// 		go serialPort.SendFrame(zwave.NewRequestFrame(frame), func(status uint8, response *zwave.ZFrame) {
	// 			if status == zwave.FrameHeaderAck {
	// 				fmt.Println("Exited inclusion mode")
	// 			} else {
	// 				fmt.Println("Could not exit inclusion mode")
	// 			}
	// 		})
	// 	}
	//
	// 	if frame.Payload[0] == 0x4b && frame.Payload[2] == 5 {
	// 		fmt.Println("NODE REMOVED")
	// 		removeNodeEnd := functions.NewAddNode()
	// 		frame := removeNodeEnd.Marshal()
	// 		go serialPort.SendFrame(zwave.NewRequestFrame(frame), func(status uint8, response *zwave.ZFrame) {
	// 			if status == zwave.FrameHeaderAck {
	// 				fmt.Println("Exited exclusion mode")
	// 			} else {
	// 				fmt.Println("Could not exit exclusion mode")
	// 			}
	// 		})
	// 	}
	// }

	// defer serialPort.Close()

}
