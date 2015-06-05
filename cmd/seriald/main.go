package main

import (
	"fmt"

	"github.com/bjyoungblood/gozw/zwave"
	"github.com/olebedev/config"
)

func loadConfigFromYaml(path string) (*zwave.SerialConfig, error) {
	config, err := config.ParseYamlFile(path)
	if err != nil {
		return nil, err
	}

	device, err := config.String("controller.device")
	if err != nil {
		return nil, err
	}

	baud, err := config.Int("controller.baud")
	if err != nil {
		return nil, err
	}

	zwaveConfig := zwave.SerialConfig{
		Device: device,
		Baud:   baud,
	}

	return &zwaveConfig, nil
}

func main() {

	config, err := loadConfigFromYaml("./zwconfig.yaml")
	if err != nil {
		panic(err)
	}

	serialPort, err := zwave.NewSerialPort(config)
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
		fmt.Println("MAIN:", <-serialPort.Incoming)
	}

	// defer serialPort.Close()

}
