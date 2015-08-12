package {{.GetPackageName}}

// {{.Help}}
{{$version := .Version}}
{{range $_, $command := .Commands}}
{{$typeName := (ToGoName $command.Name) "V" $version}}
type {{$typeName}} struct {
  {{range $_, $param := $command.Params}}{{template "commandClassStruct.tpl" $param}}{{end}}
}

func Parse{{$typeName}}(payload []byte) {{$typeName}} {
  val := {{$typeName}}{}

  {{range $_, $param := .Params}}
  {{template "commandClassParseParam.tpl" $param}}
  {{end}}

  return val
}

{{end}}
