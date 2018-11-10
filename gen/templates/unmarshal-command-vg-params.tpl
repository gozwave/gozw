{{ range $_, $param := .Params}}
    {{if eq $param.Type "VARIANT"}}
        {{template "unmarshal-variant" .}}
    {{else if eq $param.Type "STRUCT_BYTE"}}{{$name := ToGoName $param.Name}}
        if len(payload) <= i {
            return errors.New("slice index out of bounds")
        }

        var {{ToGoNameLower $param.Name}} {{$name}}
        {{range $_, $bf := $param.BitField}}
            {{if $bf.IsNotReserved}}
                {{ToGoNameLower $param.Name}}.{{ToGoName $bf.FieldName}} = (payload[i]{{with $bf.FieldMask}}&{{.}}{{end}}){{with $bf.Shifter}}>>{{.}}{{end}}
            {{end}}
        {{end}}
        {{range $_, $fe := $param.FieldEnum}}
            {{ToGoNameLower $param.Name}}.{{ToGoName $fe.FieldName}} = (payload[i]{{with $fe.FieldMask}}&{{.}}{{end}}){{with $fe.Shifter}}>>{{.}}{{end}}
        {{end}}
        {{range $_, $bf := $param.BitFlag}}
            {{if $bf.IsNotReserved}}
                {{ToGoNameLower $param.Name}}.{{ToGoName $bf.FlagName}} = payload[i] & {{$bf.FlagMask}} == {{$bf.FlagMask}}
            {{end}}
        {{end}}
        i += 1
    {{else if eq $param.Type "ARRAY"}}
        if len(payload) <= i {
            return errors.New("slice index out of bounds")
        }

        {{if (index $param.ArrayAttrib 0).IsAscii}}
            {{ToGoNameLower $param.Name}} := string(payload[i:i+{{(index $param.ArrayAttrib 0).Length}}])
        {{else}}
            {{ToGoNameLower $param.Name}} := payload[i:i+{{(index $param.ArrayAttrib 0).Length}}]
        {{end}}
        i += {{(index $param.ArrayAttrib 0).Length}}
    {{else if eq $param.Type "BITMASK"}}
        if len(payload) <= i {
            return errors.New("slice index out of bounds")
        }

        {{ToGoNameLower $param.Name}} := payload[i:]
    {{else if eq $param.Type "DWORD"}}
        if len(payload) <= i {
            return errors.New("slice index out of bounds")
        }

        {{ToGoNameLower $param.Name}} := binary.BigEndian.Uint32(payload[i:i+4])
        i += 4
    {{else if eq $param.Type "BIT_24"}}
        if len(payload) <= i {
            return errors.New("slice index out of bounds")
        }

        {{ToGoNameLower $param.Name}} := binary.BigEndian.Uint32(payload[i:i+3])
        i += 3
    {{else if eq $param.Type "WORD"}}
        if len(payload) <= i {
            return errors.New("slice index out of bounds")
        }

        {{ToGoNameLower $param.Name}} := binary.BigEndian.Uint16(payload[i:i+2])
        i += 2
    {{else if eq $param.Type "MARKER"}}
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
        
        {{if $param.IsNotReserved}}
            {{ToGoNameLower $param.Name}} := payload[i]
            i++
        {{end}}
    {{end}}
{{end}}