{{$variant := (index .Variant 0)}}{{if $variant.MarkerDelimited}}
if len(payload) <= i {
	return fmt.Errorf("slice index out of bounds (.{{ToGoName .Name}}) %d<=%d", len(payload), i)
}

{
	fieldStart := i
	for ; i < len(payload) && payload[i] != {{$variant.MarkerValue}}; i++ {}
	cmd.{{ToGoName .Name}} = payload[fieldStart:i]
}
{{else}}
{{if ne $variant.ParamOffset 255}}
if len(payload) <= i {
	return fmt.Errorf("slice index out of bounds (.{{ToGoName .Name}}) %d<=%d", len(payload), i)
}

{
	length := (payload[{{$variant.ParamOffset}}+2]{{with $variant.SizeOffset}}>>{{.}}{{end}}){{with $variant.SizeMask}}&{{.}}{{end}}
	cmd.{{ToGoName .Name}} = payload[i:i+int(length)]
	i += int(length)
}
{{else if ne $variant.RemainingBytes 0}}
if len(payload) <= i {
	return fmt.Errorf("slice index out of bounds (.{{ToGoName .Name}}) %d<=%d", len(payload), i)
}

cmd.{{ToGoName .Name}} = payload[i:len(payload)-{{$variant.RemainingBytes}}]
i += len(cmd.{{ToGoName .Name}})
{{else}}
if len(payload) <= i {
	return nil
}

cmd.{{ToGoName .Name}} = payload[i:]
{{end}}
{{end}}
