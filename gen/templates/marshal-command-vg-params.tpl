for _, vg := range cmd.{{ ToGoName .Name}} {
    {{ range $_, $param := .Params}}
        {{if eq $param.Type "VARIANT"}}
            {{if eq (index $param.Variant 0).ParamOffset 255}}
                {{template "marshal-variant" $param}}
            {{else}}
                if vg.{{ToGoName $param.Name}} != nil && len(vg.{{ToGoName $param.Name}}) > 0 {
                    payload = append(payload, vg.{{ToGoName $param.Name}}...)
                }
            {{end}}
        {{else if eq $param.Type "STRUCT_BYTE"}}
            {{$name := ToGoName $param.Name}}
            {
            var val byte
            {{range $_, $bf := $param.BitField}}
                {{if $bf.IsNotReserved}}
                val |= (vg.{{$name}}.{{ToGoName $bf.FieldName}}{{with $bf.Shifter}}<<byte({{.}}){{end}}){{with $bf.FieldMask}}&byte({{.}}){{end}}
                {{end}}
            {{end}}
            {{range $_, $fe := $param.FieldEnum}}
                val |= (vg.{{$name}}.{{ToGoName $fe.FieldName}}{{with $fe.Shifter}}<<byte({{.}}){{end}}){{with $fe.FieldMask}}&byte({{.}}){{end}}
            {{end}}
            {{range $_, $bf := $param.BitFlag}}
                {{if $bf.IsNotReserved}}
                if vg.{{$name}}.{{ToGoName $bf.FlagName}} {
                    val |= byte({{$bf.FlagMask}}) // flip bits on
                } else {
                    val &= ^byte({{$bf.FlagMask}}) // flip bits off
                }
                {{end}}
            {{end}}
            payload = append(payload, val)
            }
        {{else if eq $param.Type "ARRAY"}}
            if paramLen := len(vg.{{ToGoName $param.Name}}); paramLen > {{(index $param.ArrayAttrib 0).Length}} {
            return nil, errors.New("Length overflow in array parameter {{ToGoName $param.Name}}")
            }
            {{if (index $param.ArrayAttrib 0).IsAscii}}
            payload = append(payload, []byte(vg.{{ToGoName $param.Name}})...)
            {{else}}
            payload = append(payload, vg.{{ToGoName $param.Name}}...)
            {{end}}
        {{else if eq $param.Type "BITMASK"}}
            payload = append(payload, vg.{{ToGoName $param.Name}}...)
        {{else if eq $param.Type "DWORD"}}
            {
            buf := make([]byte, 4)
            binary.BigEndian.PutUint32(buf, vg.{{ToGoName $param.Name}})
            payload = append(payload, buf...)
            }
        {{else if eq $param.Type "BIT_24"}}
            {
            buf := make([]byte, 4)
            binary.BigEndian.PutUint32(buf, vg.{{ToGoName $param.Name}})
            if buf[0] != 0 {
                return nil, errors.New("BIT_24 value overflow")
            }
            payload = append(payload, buf[1:4]...)
            }
        {{else if eq $param.Type "WORD"}}
            {
            buf := make([]byte, 2)
            binary.BigEndian.PutUint16(buf, vg.{{ToGoName $param.Name}})
            payload = append(payload, buf...)
            }
        {{else if eq $param.Type "BYTE"}}
            payload = append(payload, vg.{{ToGoName $param.Name}})
        {{else if eq $param.Type "MARKER"}}
            payload = append(payload, {{(index $param.Const 0).FlagMask}}) // marker
        {{else}}
            {{if $param.IsNotReserved}}
            payload = append(payload, vg.{{ToGoName $param.Name}})
            {{end}}
        {{end}}
    {{end}}
}