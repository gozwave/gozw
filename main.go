package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/helioslabs/zwgen/zwgen"
)

func main() {

	app := cli.NewApp()
	app.Name = "zwgen"
	app.Usage = "Generate code for the Z-Wave protocol"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output directory",
		},
	}

	app.Before = func(c *cli.Context) error {
		outDir := c.GlobalString("output")
		if outDir == "" {
			return errors.New("Must specify output directory")
		}

		return nil
	}

	app.Action = func(ctx *cli.Context) {
		gen, err := zwgen.NewGenerator(ctx.GlobalString("output"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = gen.GenDevices()
		if err != nil {
			panic(err)
		}

		err = gen.GenCommandClasses()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = gen.GenParser()
		if err != nil {
			panic(err)
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
