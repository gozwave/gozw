{{if eq .Type "STRUCT_BYTE"}}
{{range $_, $subVal := .BitField}}{{ToPascalCase .FieldName}} byte
{{end}}
{{range $_, $subVal := .BitFlag}}{{ToPascalCase .FlagName}} bool
{{end}}
{{else}}
{{ToPascalCase .Name}} {{GetGoType .Type}}
{{end}}
