{{if .IsNotReserved}}
  {{if eq .Type "STRUCT_BYTE"}}
  {{ToGoName .Name}} struct {
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
    {{range $_, $subVal := .FieldEnum}}
      {{if .IsNotReserved}}
        {{ToGoName .FieldName}} byte
      {{end}}
    {{end}}
  }
  {{else}}
    {{ToGoName .Name}} {{.GetGoType}}
  {{end}}
{{end}}
