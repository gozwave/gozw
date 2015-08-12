{{with .}}i := 2{{end}}
{{range $_, $param := .}}
  {{if eq .Type "VARIANT"}}
    {{if eq (index .Variant 0).ParamOffset 255}}
      val.{{ToGoName .Name}} = payload[i:]
    {{else}}
      val.{{ToGoName .Name}} = payload[i:i+{{(index .Variant 0).ParamOffset}}]
      i += {{(index .Variant 0).ParamOffset}}
    {{end}}
  {{else if eq .Type "STRUCT_BYTE"}}
    {{range $_, $subVal := .BitField}}
      {{if .IsNotReserved}}
        val.{{ToGoName .FieldName}} = (payload[i]{{with .FieldMask}}&{{.}}{{end}}){{with .Shifter}}<<{{.}}{{end}}
      {{end}}
    {{end}}
    {{range $_, $subVal := .FieldEnum}}
      val.{{ToGoName .FieldName}} = (payload[i]{{with .FieldMask}}&{{.}}{{end}}){{with .Shifter}}<<{{.}}{{end}}
    {{end}}
    {{range $_, $subVal := .BitFlag}}
      {{if .IsNotReserved}}
        if payload[i] & {{.FlagMask}} == {{.FlagMask}} {
          val.{{ToGoName .FlagName}} = true
        } else {
          val.{{ToGoName .FlagName}} = false
        }
      {{end}}
    {{end}}
    i += 1
  {{else if eq .Type "ARRAY"}}
    {{if (index .ArrayAttrib 0).IsAscii}}
      val.{{ToGoName .Name}} = string(payload[i:i+{{(index .ArrayAttrib 0).Length}}])
    {{else}}
      val.{{ToGoName .Name}} = payload[i:i+{{(index .ArrayAttrib 0).Length}}]
    {{end}}
    i += {{(index .ArrayAttrib 0).Length}}
  {{else if eq .Type "DWORD"}}
    val.{{ToGoName .Name}} = binary.BigEndian.Uint32(payload[i:i+4])
    i += 4
  {{else if eq .Type "BIT_24"}}
    val.{{ToGoName .Name}} = binary.BigEndian.Uint32(payload[i:i+3])
    i += 3
  {{else if eq .Type "WORD"}}
    val.{{ToGoName .Name}} = binary.BigEndian.Uint16(payload[i:i+2])
    i += 2
  {{else if eq .Type "MARKER"}}
    // MARKER HERE
  {{else}}
    {{if .IsNotReserved}}
      val.{{ToGoName .Name}} = payload[i]
      i++
    {{end}}
  {{end}}
{{end}}
