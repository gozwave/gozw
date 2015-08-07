package ccgen

const commandClassTemplate = `package {{toPackageName Name}}

const

// {{Help}}
{{#each Commands}}
type {{toPascalCase Name}}V{{../Version}} struct {
  {{#each Params}}
  {{toPascalCase Name}} {{getGoType Type}}
  {{/each}}
}

func Parse{{toPascalCase Name}}V{{../Version}}(payload []byte) {{toPascalCase Name}}V{{../Version}} {
  val := {{toPascalCase Name}}V{{../Version}}{
    {{#each Params}}
    {{> paramInStruct}}
    {{/each}}
  }

  return val
}

{{/each}}
`
