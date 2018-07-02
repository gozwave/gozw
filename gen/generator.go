package gen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"

	"golang.org/x/tools/imports"
)

//go:generate go-bindata -pkg=gen templates/... data/...

type Generator struct {
	output    string
	config    Config
	zwClasses *ZwClasses
	tpl       *template.Template
}

type Config struct {
	CommandClasses map[string]map[int]bool `yaml:"CommandClasses"`
}

func NewGenerator(output string, configFile string) (*Generator, error) {
	config := Config{}

	configStr, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(configStr, &config); err != nil {
		return nil, err
	}

	gen := &Generator{
		output: output,
		config: config,
	}

	zwData, err := Asset("data/zwave-defs.xml")
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(bytes.NewBuffer(zwData))

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

	err := g.GenerateDevices(buf)
	if err != nil {
		return err
	}

	fp, err := os.Create(g.output)
	if err != nil {
		return err
	}

	defer fp.Close()

	formatted, err := goFmtAndImports(g.output, buf)
	if err != nil {
		return err
	}

	fp.Write(formatted)

	return nil
}

func (g *Generator) GenParser() error {

	if g.config.CommandClasses != nil && len(g.config.CommandClasses) > 0 {
		for i, cc := range g.zwClasses.CommandClasses {
			if _, ok := g.config.CommandClasses[cc.Name]; !ok {
				g.zwClasses.CommandClasses[i].Enabled = false
				continue
			} else {
				g.zwClasses.CommandClasses[i].Enabled = true
			}

			if g.config.CommandClasses[cc.Name] != nil && len(g.config.CommandClasses[cc.Name]) > 0 {
				if should, ok := g.config.CommandClasses[cc.Name][cc.Version]; !ok || !should {
					g.zwClasses.CommandClasses[i].Enabled = false
				} else {
					g.zwClasses.CommandClasses[i].Enabled = true
				}
			}
		}
	}

	buf := bytes.NewBuffer([]byte{})

	err := g.GenerateCommandClasses(buf)
	if err != nil {
		return err
	}

	fp, err := os.Create(g.output)
	if err != nil {
		return err
	}

	defer fp.Close()

	formatted, err := goFmtAndImports(g.output, buf)
	if err != nil {
		return err
	}

	fp.Write(formatted)

	return nil
}

func (g *Generator) GenCommandClasses() error {
	skipped := []CommandClass{}
	for _, cc := range g.zwClasses.CommandClasses {

		if g.config.CommandClasses != nil && len(g.config.CommandClasses) > 0 {
			if _, ok := g.config.CommandClasses[cc.Name]; !ok {
				continue
			}

			if g.config.CommandClasses[cc.Name] != nil && len(g.config.CommandClasses[cc.Name]) > 0 {
				if should, ok := g.config.CommandClasses[cc.Name][cc.Version]; !ok || !should {
					continue
				}
			}
		}

		if can, _ := cc.CanGenerate(); !can {
			skipped = append(skipped, cc)
			continue
		}

		dirName := path.Join(g.output, cc.GetDirName())
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

	err := g.GenerateCommand(buf, cc, cmd)
	if err != nil {
		return err
	}

	filename := path.Join(dirName, cmd.GetFileName(cc)+".gen.go")
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer fp.Close()

	// formatted, err := goFmtAndImports(filename, buf)
	// if err != nil {
	// 	return err
	// }

	fp.Write(buf.Bytes())

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
					if param.Variant[0].ParamOffset != 255 {
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

func goFmtAndImports(filename string, buf *bytes.Buffer) ([]byte, error) {
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	imported, err := imports.Process(filename, formatted, nil)
	if err != nil {
		return nil, err
	}

	return imported, nil
}
