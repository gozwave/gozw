{{ $keys := .OrderedKeys}}
{{ $command := .}}
{{ range $key := $keys}}
  {{ $type := $command.GetFieldType $key}}
  {{- if eq $type "param"}}
    {{ $param := $command.GetParam $key }}
    {{template "marshal-command-param" $param}}
  {{- else if eq $type "vg"}}
    {{ $vg := $command.GetVg $key }}
    {{template "marshal-command-vg-params" $vg}}
  {{- end}}
{{end}}