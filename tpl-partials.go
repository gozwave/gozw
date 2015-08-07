package ccgen

import "github.com/aymerick/raymond"

func init() {
	raymond.RegisterPartials(map[string]string{
		"paramInStruct": paramInStruct,
	})
}

const paramInStruct = `
{{toPascalCase Name}}: payload[{{Key}}],
`
