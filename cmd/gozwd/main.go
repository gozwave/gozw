package main

import
// "bufio"

// "io"
// "io/ioutil"
// "os"

"github.com/bjyoungblood/gozw/portal"

func main() {

	config, err := portal.LoadConfigFromYaml("./config.yaml")
	if err != nil {
		panic(err)
	}

	server, err := portal.NewPortalServer(config)
	if err != nil {
		panic(err)
	}

	server.Start()

}
