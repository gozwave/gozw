package main

import (
	"github.com/bjyoungblood/gozw/common"
	"github.com/bjyoungblood/gozw/portal"
)

func main() {

	config, err := common.LoadPortalConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	server, err := portal.NewPortalServer(config)
	if err != nil {
		panic(err)
	}

	server.Start()

}
