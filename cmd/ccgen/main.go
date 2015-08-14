package main

import "github.com/helioslabs/gozw/ccgen"

func main() {

	gen, err := ccgen.NewGenerator()
	if err != nil {
		panic(err)
	}

	err = gen.GenDevices()
	if err != nil {
		panic(err)
	}

	err = gen.GenCommandClasses()
	if err != nil {
		panic(err)
	}

	err = gen.GenParser()
	if err != nil {
		panic(err)
	}
}
