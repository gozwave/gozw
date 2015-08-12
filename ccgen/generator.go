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

	return gen, nil
}

func (g *Generator) GenDevices() (string, error) {
	buf := bytes.NewBuffer([]byte{})
	err := g.tpl.ExecuteTemplate(buf, "devices.tpl", g.zwClasses)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func (g *Generator) GenCommandClasses() (string, error) {
	buf := bytes.NewBuffer([]byte{})

	skipped := []CommandClass{}
	for _, cc := range g.zwClasses.CommandClasses {

		if can, _ := cc.CanGenerate(); !can {
			skipped = append(skipped, cc)
			continue
		}

		// fmt.Println(cc.Name)
		buf := bytes.NewBuffer([]byte{})
		err := g.tpl.ExecuteTemplate(buf, "commandClass.tpl", cc)
		if err != nil {
			panic(err)
			return "", err
		}

		dirName := "zwave/command-class/" + cc.GetPackageName()
		filename := dirName + "/" + cc.GetPackageName() + ".go"
		os.Mkdir(dirName, 0775)
		fp, err := os.Create(filename)
		if err != nil {
			panic(err)
		}

		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println(string(buf.Bytes()))
			fmt.Println(cc.Name)
			panic(err)
		}

		imported, err := imports.Process(filename, formatted, nil)
		if err != nil {
			fmt.Println(cc.Name)
			panic(err)
		}

		fp.Write(imported)
		fp.Close()
	}

	if len(skipped) > 0 {
		fmt.Println("Skipped generation for the following command classes:")
		for _, cc := range skipped {
			_, reason := cc.CanGenerate()
			fmt.Printf("  - %s\n      Reason: %s\n", cc.Name, reason)
		}
	}

	return string(buf.Bytes()), nil
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
