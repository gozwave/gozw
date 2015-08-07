package main

import (
	"os"
	"os/signal"

	"github.com/bjyoungblood/gozw/gateway"
)

func main() {
	opts := gateway.GatewayOptions{
		CommNetType: "unix",
		CommAddress: "/tmp/arc",

		ZWaveSerialPort: "/tmp/usbmodem",
		BaudRate:        115200,
	}

	gw, err := gateway.NewGateway(opts)
	if err != nil {
		panic(err)
	}

	defer gw.Shutdown()

	gw.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

}
