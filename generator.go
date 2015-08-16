package ccgen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

type Generator struct {
	zwClasses *ZwClasses
	tpl       *template.Template
}

func NewGenerator() (*Generator, error) {
	gen := &Generator{}

	if err := gen.initTemplates(); err != nil {
		return nil, err
	}

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

	err = gen.fixVariants()
	if err != nil {
		return nil, err
	}

	return gen, nil
}

func (g *Generator) GenDevices() error {
	buf := bytes.NewBuffer([]byte{})
	err := g.tpl.ExecuteTemplate(buf, "devices.tpl", g.zwClasses)
	if err != nil {
		return err
	}

	filename := "zwave/command-class/devices.gen.go"
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer fp.Close()

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	imported, err := imports.Process(filename, formatted, nil)
	if err != nil {
		return err
	}

	fp.Write(imported)

	return nil
}

func (g *Generator) GenParser() error {
	buf := bytes.NewBuffer([]byte{})
	err := g.tpl.ExecuteTemplate(buf, "command-classes.tpl", g.zwClasses)
	if err != nil {
		return err
	}

	filename := "zwave/command-class/command-classes.gen.go"
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer fp.Close()

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	imported, err := imports.Process(filename, formatted, nil)
	if err != nil {
		return err
	}

	fp.Write(imported)

	return nil
}

func (g *Generator) GenCommandClasses() error {
	skipped := []CommandClass{}
	for _, cc := range g.zwClasses.CommandClasses {

		if can, _ := cc.CanGenerate(); !can {
			skipped = append(skipped, cc)
			continue
		}

		dirName := "zwave/command-class/" + cc.GetDirName()
		err := os.Mkdir(dirName, 0775)
		if err != nil && !strings.HasSuffix(err.Error(), "file exists") {
			return err
		}

		for _, cmd := range cc.Commands {
			err := g.generateCommand(dirName, cc, cmd)
			if err != nil {
				return err
			}
		}

	}

	if len(skipped) > 0 {
		fmt.Println("Skipped generation for the following command classes:")
		for _, cc := range skipped {
			_, reason := cc.CanGenerate()
			fmt.Printf("  - %s (%s)\n", cc.Name, reason)
		}
	}

	return nil
}

func (g *Generator) generateCommand(dirName string, cc CommandClass, cmd Command) error {
	buf := bytes.NewBuffer([]byte{})

	err := g.tpl.ExecuteTemplate(buf, "command.tpl", map[string]interface{}{
		"CommandClass": cc,
		"Command":      cmd,
	})
	if err != nil {
		return err
	}

	filename := dirName + "/" + cmd.GetFileName(cc) + ".gen.go"
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer fp.Close()

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(string(buf.Bytes()))
		fmt.Println(cc.Name)
		return err
	}

	imported, err := imports.Process(filename, formatted, nil)
	if err != nil {
		fmt.Println(cc.Name)
		return err
	}

	fp.Write(imported)

	return nil
}

func (g *Generator) initTemplates() error {
	funcs := template.FuncMap{
		"ToGoName":    toGoName,
		"NotZeroByte": notZeroByte,
		"Trim":        strings.TrimSpace,
	}

	tpl, err := template.New("").Funcs(funcs).ParseGlob("ccgen/templates/*.tpl")
	if err != nil {
		return err
	}

	g.tpl = tpl

	return nil
}

func (g *Generator) fixVariants() error {
	for _, cc := range g.zwClasses.CommandClasses {
		if !cc.CanGen() {
			continue
		}

		for _, cmd := range cc.Commands {
			for i, param := range cmd.Params {

				if param.Type == "VARIANT" {
					if param.Variant[0].ParamOffset != byte(255) {
						continue
					}

					if len(cmd.Params) > i+1 && cmd.Params[i+1].Type == "MARKER" {
						// This command is followed by marker, where it stops
						param.Variant[0].MarkerDelimited = true
						param.Variant[0].MarkerValue = cmd.Params[i+1].Const[0].FlagMask
					} else if len(cmd.Params) == i+1 {
						// If this is the last param, we don't need to do anything special
						param.Variant[0].MarkerDelimited = false
					} else {
						// This variant is followed by more bytes (only applies to a few commands
						// like Security / Message Encapsulation and Firmware Update Md)

						var remainingBytes uint8
						for j := i + 1; j < len(cmd.Params); j++ {
							paramLength, err := cmd.Params[j].GetEncodedByteLength()
							if err != nil {
								return err
							}

							remainingBytes += paramLength
						}

						param.Variant[0].RemainingBytes = remainingBytes
					}
				}

			}
		}
	}

	return nil
}
