package main

import (
	"os"
	"os/signal"

	"github.com/bjyoungblood/gozw/arcturus"
)

func main() {

	server := arcturus.NewServer()
	go server.Listen("unix", "/tmp/arc")

	defer server.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
