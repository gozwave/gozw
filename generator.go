package ccgen

import (
	"encoding/xml"
	"os"

	"github.com/aymerick/raymond"
)

type Generator struct {
	zwClasses *ZwClasses
}

func NewGenerator() (*Generator, error) {
	gen := &Generator{}

	fp, err := os.Open("ccgen/zwave-defs.xml")
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(fp)

	zw := ZwClasses{}
	err = decoder.Decode(&zw)
	if err != nil {
		return nil, err
	}

	gen.zwClasses = &zw

	return gen, nil
}

func (g *Generator) GenDevices() (string, error) {
	devices, err := raymond.Parse(devicesTemplate)
	if err != nil {
		return "", err
	}

	return devices.Exec(g.zwClasses)
}

func (g *Generator) GenCommandClasses() (string, error) {
	devices, err := raymond.Parse(commandClassTemplate)
	if err != nil {
		return "", err
	}

	return devices.Exec(g.zwClasses.CommandClasses[1])
}
