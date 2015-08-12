{{if .IsNotReserved}}
  {{if eq .Type "STRUCT_BYTE"}}
    {{range $_, $subVal := .BitField}}
      {{if .IsNotReserved}}
        {{ToGoName .FieldName}} byte
      {{end}}
    {{end}}
    {{range $_, $subVal := .BitFlag}}
      {{if .IsNotReserved}}
        {{ToGoName .FlagName}} bool
      {{end}}
    {{end}}
  {{else}}
    {{ToGoName .Name}} {{.GetGoType}}
  {{end}}
{{end}}
