{{if eq .Type "STRUCT_BYTE"}}
{{range $_, $subVal := .BitField}}{{ToGoName .FieldName}} byte
{{end}}
{{range $_, $subVal := .BitFlag}}{{ToGoName .FlagName}} bool
{{end}}
{{else}}
{{ToGoName .Name}} {{.GetGoType}}
{{end}}
