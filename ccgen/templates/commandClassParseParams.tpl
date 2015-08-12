{{with .}}i := 2{{end}}
{{range $_, $param := .}}
{{if eq .Type "VARIANT"}}
  val.{{ToGoName .Name}} = payload[i:i+{{(index .Variant 0).ParamOffset}}]
  i += {{(index .Variant 0).ParamOffset}}
{{else if eq .Type "STRUCT_BYTE"}}
  {{range $_, $subVal := .BitField}}
    val.{{ToGoName .FieldName}} = (payload[i]{{with .FieldMask}}&{{.}}{{end}}){{with .Shifter}}<<{{.}}{{end}}
  {{end}}
  {{range $_, $subVal := .BitFlag}}
    if payload[i] & {{.FlagMask}} == {{.FlagMask}} {
      val.{{ToGoName .FlagName}} = true
    } else {
      val.{{ToGoName .FlagName}} = false
    }
  {{end}}
  i += 1
{{else if eq .Type "BIT_24"}}
  val.{{ToGoName .Name}} = binary.BigEndian.Uint32(payload[i:i+3])
  i += 3
{{else}}
  val.{{ToGoName .Name}} = payload[i]
  i++
{{end}}
{{end}}
