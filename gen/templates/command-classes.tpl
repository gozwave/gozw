// THIS FILE IS AUTO-GENERATED
// DO NOT MODIFY

package cc

const (
  {{range .CommandClasses}}
  {{.GetConstName}} CommandClassID = {{.Key}}{{end}}
)

func (c CommandClassID) String() string {
  switch c {
    {{range .CommandClasses}}
    {{if eq .Version 1}}
    case {{.GetConstName}}:
      return "{{.Help}}"
    {{end}}
    {{end}}
    default:
      return fmt.Sprintf("Unknown (0x%X)", byte(c))
  }
}
