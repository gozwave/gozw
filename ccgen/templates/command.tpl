// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package {{.CommandClass.GetPackageName}}

// {{.Help}}
{{$version := .CommandClass.Version}}
{{$typeName := (ToGoName .Command.Name) "V" $version}}
type {{$typeName}} struct {
  {{range $_, $param := .Command.Params}}
    {{template "commandStruct.tpl" $param}}
  {{end}}
}

func Parse{{$typeName}}(payload []byte) {{$typeName}} {
  val := {{$typeName}}{}

  {{template "commandParseParams.tpl" .Command.Params}}

  return val
}