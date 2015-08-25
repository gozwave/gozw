{{$variant := (index .Variant 0)}}{{if $variant.MarkerDelimited}}
{
  fieldStart := i
  for ; i < len(payload) && payload[i] != {{$variant.MarkerValue}}; i++ {}
  cmd.{{ToGoName .Name}} = payload[fieldStart:i]
}
{{else}}
{{if ne $variant.ParamOffset 255}}
{
  length := (payload[{{$variant.ParamOffset}}+2]{{with $variant.SizeOffset}}>>{{.}}{{end}}){{with $variant.SizeMask}}&{{.}}{{end}}
  cmd.{{ToGoName .Name}} = payload[i:i+int(length)]
  i += int(length)
}
{{else if ne $variant.RemainingBytes 0}}
cmd.{{ToGoName .Name}} = payload[i:len(payload)-{{$variant.RemainingBytes}}]
i += len(cmd.{{ToGoName .Name}})
{{else}}
cmd.{{ToGoName .Name}} = payload[i:]
{{end}}
{{end}}
