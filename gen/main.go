package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
)

func main() {

	app := cli.NewApp()
	app.Name = ""
	app.Usage = "Generate code for the Z-Wave protocol"

	before := func(c *cli.Context) error {
		outDir := c.String("output")
		if outDir == "" {
			return errors.New("Must specify output directory")
		}

		config := c.String("config")
		if config == "" {
			return errors.New("Must specify config file")
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "devices",
			Usage: "Generate device info",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output filename",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Config file",
				},
			},
			Before: before,
			Action: func(ctx *cli.Context) {
				gen, err := NewGenerator(ctx.String("output"), ctx.String("config"))
				if err != nil {

					fmt.Println(errors.Wrap(err, "new generator"))
					os.Exit(1)
				}

				err = gen.GenDevices()
				if err != nil {
					panic(err)
				}
			},
		},

		{
			Name:  "command-classes",
			Usage: "Generate command class",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output directory",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Config file",
				},
			},
			Before: before,
			Action: func(ctx *cli.Context) {
				gen, err := NewGenerator(ctx.String("output"), ctx.String("config"))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				err = gen.GenCommandClasses()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			},
		},

		{
			Name:  "parser",
			Usage: "Generate command class parser",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output directory",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Config file",
				},
			},
			Before: before,
			Action: func(ctx *cli.Context) {
				gen, err := NewGenerator(ctx.String("output"), ctx.String("config"))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				err = gen.GenParser()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
