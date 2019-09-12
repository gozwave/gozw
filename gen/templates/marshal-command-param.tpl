{{if eq .Type "VARIANT"}}
    {{if eq (index .Variant 0).ParamOffset 255}}
      {{template "marshal-variant" .}}
    {{else}}
      if cmd.{{ToGoName .Name}} != nil && len(cmd.{{ToGoName .Name}}) > 0 {
        payload = append(payload, cmd.{{ToGoName .Name}}...)
      }
    {{end}}
  {{else if eq .Type "STRUCT_BYTE"}}
    {{$name := ToGoName .Name}}
    {
      var val byte
      {{range .BitField}}
        {{if .IsNotReserved}}
          val |= (cmd.{{$name}}.{{ToGoName .FieldName}}{{with .Shifter}}<<byte({{.}}){{end}}){{with .FieldMask}}&byte({{.}}){{end}}
        {{end}}
      {{end}}
      {{range .FieldEnum}}
        val |= (cmd.{{$name}}.{{ToGoName .FieldName}}{{with .Shifter}}<<byte({{.}}){{end}}){{with .FieldMask}}&byte({{.}}){{end}}
      {{end}}
      {{range .BitFlag}}
        {{if .IsNotReserved}}
          if cmd.{{$name}}.{{ToGoName .FlagName}} {
            val |= byte({{.FlagMask}}) // flip bits on
          } else {
            val &= ^byte({{.FlagMask}}) // flip bits off
          }
        {{end}}
      {{end}}
      payload = append(payload, val)
    }
  {{else if eq .Type "ARRAY"}}
    if paramLen := len(cmd.{{ToGoName .Name}}); paramLen > {{(index .ArrayAttrib 0).Length}} {
      return nil, errors.New("Length overflow in array parameter {{ToGoName .Name}}")
    }
    {{if (index .ArrayAttrib 0).IsAscii}}
      payload = append(payload, []byte(cmd.{{ToGoName .Name}})...)
    {{else}}
      payload = append(payload, cmd.{{ToGoName .Name}}...)
    {{end}}
  {{else if eq .Type "BITMASK"}}
    payload = append(payload, cmd.{{ToGoName .Name}}...)
  {{else if eq .Type "DWORD"}}
    {
      buf := make([]byte, 4)
      binary.BigEndian.PutUint32(buf, cmd.{{ToGoName .Name}})
      payload = append(payload, buf...)
    }
  {{else if eq .Type "BIT_24"}}
    {
      buf := make([]byte, 4)
      binary.BigEndian.PutUint32(buf, cmd.{{ToGoName .Name}})
      if buf[0] != 0 {
        return nil, errors.New("BIT_24 value overflow")
      }
      payload = append(payload, buf[1:4]...)
    }
  {{else if eq .Type "WORD"}}
    {
      buf := make([]byte, 2)
      binary.BigEndian.PutUint16(buf, cmd.{{ToGoName .Name}})
      payload = append(payload, buf...)
    }
  {{else if eq .Type "MARKER"}}
    payload = append(payload, {{(index .Const 0).FlagMask}}) // marker
  {{else}}
    {{if .IsNotReserved}}
      payload = append(payload, cmd.{{ToGoName .Name}})
    {{end}}
  {{end}}