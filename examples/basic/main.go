package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gozwave/gozw"
	switchbinary "github.com/gozwave/gozw/cc/switch-binary"
	"github.com/davecgh/go-spew/spew"
)

var networkKey = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

func main() {
	client, err := gozw.NewDefaultClient("/tmp/data.db", "/dev/ttyACM0", 115200, networkKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	spew.Dump(client.Controller)

	for _, node := range client.Nodes() {
		fmt.Println(node.String())
	}

	node, err := client.Node(2)
	if err != nil {
		log.Fatalf("retrieve node: %v", err)
	}

	err = node.SendCommand(&switchbinary.Get{})
	if err != nil {
		log.Fatalf("send command: %v", err)
	}

	time.Sleep(2 * time.Second)

}
