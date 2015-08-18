package zwgen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go/format"
	"os"
	"path"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

//go:generate go-bindata -pkg=zwgen templates/... data/...

type Generator struct {
	outputDir string
	zwClasses *ZwClasses
	tpl       *template.Template
}

func NewGenerator(outputDir string) (*Generator, error) {
	gen := &Generator{
		outputDir: outputDir,
	}

	gen.initTemplates()

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

	err := g.tpl.ExecuteTemplate(buf, "devices", g.zwClasses)
	if err != nil {
		return err
	}

	filename := path.Join(g.outputDir, "devices.gen.go")
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer fp.Close()

	formatted, err := goFmtAndImports(filename, buf)
	if err != nil {
		return err
	}

	fp.Write(formatted)

	return nil
}

func (g *Generator) GenParser() error {
	buf := bytes.NewBuffer([]byte{})

	err := g.tpl.ExecuteTemplate(buf, "command-classes", g.zwClasses)
	if err != nil {
		return err
	}

	filename := path.Join(g.outputDir, "command-classes.gen.go")
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer fp.Close()

	formatted, err := goFmtAndImports(filename, buf)
	if err != nil {
		return err
	}

	fp.Write(formatted)

	return nil
}

func (g *Generator) GenCommandClasses() error {
	skipped := []CommandClass{}
	for _, cc := range g.zwClasses.CommandClasses {

		if can, _ := cc.CanGenerate(); !can {
			skipped = append(skipped, cc)
			continue
		}

		dirName := path.Join(g.outputDir, cc.GetDirName())
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

	err := g.tpl.ExecuteTemplate(buf, "command", map[string]interface{}{
		"CommandClass": cc,
		"Command":      cmd,
	})
	if err != nil {
		return err
	}

	filename := path.Join(dirName, cmd.GetFileName(cc)+".gen.go")
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer fp.Close()

	formatted, err := goFmtAndImports(filename, buf)
	if err != nil {
		return err
	}

	fp.Write(formatted)

	return nil
}
func (g *Generator) initTemplates() {
	tpl := template.New("").Funcs(template.FuncMap{
		"ToGoName":    toGoName,
		"NotZeroByte": notZeroByte,
		"Trim":        strings.TrimSpace,
	})

	tpl = template.Must(tpl.New("command-classes").Parse(mustAsset("templates/command-classes.tpl")))
	tpl = template.Must(tpl.New("command-struct-fields").Parse(mustAsset("templates/command-struct-fields.tpl")))
	tpl = template.Must(tpl.New("command").Parse(mustAsset("templates/command.tpl")))
	tpl = template.Must(tpl.New("devices").Parse(mustAsset("templates/devices.tpl")))
	tpl = template.Must(tpl.New("marshal-command-params").Parse(mustAsset("templates/marshal-command-params.tpl")))
	tpl = template.Must(tpl.New("marshal-variant").Parse(mustAsset("templates/marshal-variant.tpl")))
	tpl = template.Must(tpl.New("unmarshal-command-params").Parse(mustAsset("templates/unmarshal-command-params.tpl")))
	tpl = template.Must(tpl.New("unmarshal-variant").Parse(mustAsset("templates/unmarshal-variant.tpl")))

	g.tpl = tpl
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

func mustAsset(name string) string {
	str, err := Asset(name)
	if err != nil {
		panic(err)
	}

	return string(str)
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
