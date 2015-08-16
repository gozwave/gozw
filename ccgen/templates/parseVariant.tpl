{{$variant := (index .Variant 0)}}{{if $variant.StopAtMarker}}
{
  markerIndex := i
  for ; markerIndex < len(payload) && payload[markerIndex] != {{$variant.MarkerValue}}; markerIndex++ {}
  val.{{ToGoName .Name}} = payload[i:markerIndex]
}
{{else}}
val.{{ToGoName .Name}} = payload[i:]
{{end}}
