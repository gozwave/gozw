{{$variant := (index .Variant 0)}}{{if $variant.StopAtMarker}}
{
  if cmd.{{ToGoName .Name}} != nil && len(cmd.{{ToGoName .Name}}) > 0 {
    payload = append(payload, cmd.{{ToGoName .Name}}...)
  }
  payload = append(payload, {{$variant.MarkerValue}})
}
{{else}}
payload = append(payload, cmd.{{ToGoName .Name}}...)
{{end}}
