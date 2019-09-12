{{if eq .Type "VARIANT"}}
    {{template "unmarshal-variant" .}}
  {{else if eq .Type "STRUCT_BYTE"}}{{$name := ToGoName .Name}}
    if len(payload) <= i {
      return errors.New("slice index out of bounds")
    }

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
        cmd.{{$name}}.{{ToGoName .FlagName}} = payload[i] & {{.FlagMask}} == {{.FlagMask}}
      {{end}}
    {{end}}
    i += 1
  {{else if eq .Type "ARRAY"}}
    if len(payload) <= i {
      return errors.New("slice index out of bounds")
    }

    {{if (index .ArrayAttrib 0).IsAscii}}
      cmd.{{ToGoName .Name}} = string(payload[i:i+{{(index .ArrayAttrib 0).Length}}])
    {{else}}
      cmd.{{ToGoName .Name}} = payload[i:i+{{(index .ArrayAttrib 0).Length}}]
    {{end}}
    i += {{(index .ArrayAttrib 0).Length}}
  {{else if eq .Type "BITMASK"}}
    if len(payload) <= i {
      return errors.New("slice index out of bounds")
    }

    cmd.{{ToGoName .Name}} = payload[i:]
  {{else if eq .Type "DWORD"}}
    if len(payload) <= i {
      return errors.New("slice index out of bounds")
    }

    cmd.{{ToGoName .Name}} = binary.BigEndian.Uint32(payload[i:i+4])
    i += 4
  {{else if eq .Type "BIT_24"}}
    if len(payload) <= i {
      return errors.New("slice index out of bounds")
    }

    cmd.{{ToGoName .Name}} = binary.BigEndian.Uint32(payload[i:i+3])
    i += 3
  {{else if eq .Type "WORD"}}
    if len(payload) <= i {
      return errors.New("slice index out of bounds")
    }

    cmd.{{ToGoName .Name}} = binary.BigEndian.Uint16(payload[i:i+2])
    i += 2
  {{else if eq .Type "MARKER"}}
    if len(payload) <= i {
      return errors.New("slice index out of bounds")
    }

    i += 1 // skipping MARKER
    if len(payload) <= i {
      return nil
    }
  {{else}}
    if len(payload) <= i {
      return errors.New("slice index out of bounds")
    }

    {{if .IsNotReserved}}
      cmd.{{ToGoName .Name}} = payload[i]
      i++
    {{end}}
  {{end}}