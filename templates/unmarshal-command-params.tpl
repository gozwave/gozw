{{with .}}i := 0{{end}}{{/* {{with}} here ensures we don't print this if there are no params */}}
{{range .}}
  if len(payload) <= i {
    return errors.New("slice index out of bounds")
  }
  {{if eq .Type "VARIANT"}}
    {{if eq (index .Variant 0).ParamOffset 255}}
      {{template "unmarshal-variant.tpl" .}}
    {{else}}
      cmd.{{ToGoName .Name}} = payload[i:i+{{(index .Variant 0).ParamOffset}}]
      i += {{(index .Variant 0).ParamOffset}}
    {{end}}
  {{else if eq .Type "STRUCT_BYTE"}}{{$name := ToGoName .Name}}
    {{range .BitField}}
      {{if .IsNotReserved}}
        cmd.{{$name}}.{{ToGoName .FieldName}} = (payload[i]{{with .FieldMask}}&{{.}}{{end}}){{with .Shifter}}>>{{.}}{{end}}
      {{end}}
    {{end}}
    {{range .FieldEnum}}
      cmd.{{$name}}.{{ToGoName .FieldName}} = (payload[i]{{with .FieldMask}}&{{.}}{{end}}){{with .Shifter}}>>{{.}}{{end}}
    {{end}}
    {{range .BitFlag}}
      {{if .IsNotReserved}}
        if payload[i] & {{.FlagMask}} == {{.FlagMask}} {
          cmd.{{$name}}.{{ToGoName .FlagName}} = true
        } else {
          cmd.{{$name}}.{{ToGoName .FlagName}} = false
        }
      {{end}}
    {{end}}
    i += 1
  {{else if eq .Type "ARRAY"}}
    {{if (index .ArrayAttrib 0).IsAscii}}
      cmd.{{ToGoName .Name}} = string(payload[i:i+{{(index .ArrayAttrib 0).Length}}])
    {{else}}
      cmd.{{ToGoName .Name}} = payload[i:i+{{(index .ArrayAttrib 0).Length}}]
    {{end}}
    i += {{(index .ArrayAttrib 0).Length}}
  {{else if eq .Type "DWORD"}}
    cmd.{{ToGoName .Name}} = binary.BigEndian.Uint32(payload[i:i+4])
    i += 4
  {{else if eq .Type "BIT_24"}}
    cmd.{{ToGoName .Name}} = binary.BigEndian.Uint32(payload[i:i+3])
    i += 3
  {{else if eq .Type "WORD"}}
    cmd.{{ToGoName .Name}} = binary.BigEndian.Uint16(payload[i:i+2])
    i += 2
  {{else if eq .Type "MARKER"}}
    i += 1 // skipping MARKER
  {{else}}
    {{if .IsNotReserved}}
      cmd.{{ToGoName .Name}} = payload[i]
      i++
    {{end}}
  {{end}}
{{end}}
