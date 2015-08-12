{{if eq .Type "VARIANT"}}
  val.{{ToPascalCase .Name}} = payload[{{(index .Variant 0).ParamOffset}}:]
{{else if eq .Type "STRUCT_BYTE"}}
  {{range $_, $subVal := .BitField}}
    val.{{ToPascalCase .FieldName}} = (payload[{{.Key}}]{{with .FieldMask}}&{{.}}{{end}}){{with .Shifter}}<<{{.}}{{end}}
  {{end}}
  {{range $_, $subVal := .BitFlag}}
    if payload[{{.Key}}] & {{.FlagMask}} == {{.FlagMask}} {
      val.{{ToPascalCase .FlagName}} = true
    } else {
      val.{{ToPascalCase .FlagName}} = false
    }
  {{end}}
{{else}}
  val.{{ToPascalCase .Name}} = payload[{{.Key}}]
{{end}}
