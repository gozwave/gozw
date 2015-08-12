package ccgen

import (
	"bytes"
	"encoding/xml"
	"os"
	"strings"
	"text/template"

	"github.com/reiver/go-stringcase"
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
	err := g.tpl.ExecuteTemplate(buf, "commandClass.tpl", g.zwClasses.CommandClasses[1])
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func (g *Generator) initTemplates() error {
	funcs := template.FuncMap{
		"ToPascalCase":  stringcase.ToPascalCase,
		"ToPackageName": toPackageName,
		"GetGoType":     getGoType,
		"NotZeroByte":   notZeroByte,
		"Trim":          strings.TrimSpace,
	}

	tpl, err := template.New("").Funcs(funcs).ParseGlob("ccgen/templates/*.tpl")
	if err != nil {
		return err
	}

	g.tpl = tpl

	return nil
}
