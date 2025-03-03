package d2

import (
	"regexp"
	"strings"
	"text/template"

	"github.com/Emptyless/jsonschema-transform/internal/domain"
)

var ClassTemplate = NewTemplate("Class", `
"{{- $.Name }}": {
  shape: class
{{- range $property := $.Properties }}
  "{{ $property.Name }}": "{{ $property.Type }}"
{{- end }}
}`)

var RelationTemplate = NewTemplate("Relation", `
{{- $.From.Name }} -- {{ $.To.Name }}: "{{ $.Type | safe }}"`)

var ContainerTemplate = NewTemplate("Container", `
{{- if $.Name -}}
"{{ $.Name }}": {
{{- range $container := $.Containers -}}
{{ $container.Render | indent 4 }}

{{ end -}}
{{- range $class := $.Classes }}
{{ class $class | indent 4 }}
{{ end -}}
}
{{- else -}}
{{- range $container := $.Containers -}}
{{ $container.Render }}

{{ end -}}
{{- range $class := $.Classes }}
{{ class $class }}
{{ end -}}
{{- end }}`, template.FuncMap{
	"class": RenderClass,
	"indent": func(amount int, input string) string {
		lines := strings.Split(input, "\n")
		var builder strings.Builder
		for i, line := range lines {
			builder.WriteString(strings.Repeat(" ", amount) + line)
			if i != len(lines)-1 {
				builder.WriteString("\n")
			}
		}

		return builder.String()
	},
})

// RenderContainer to string output
func RenderContainer(container *Container) string {
	var builder strings.Builder
	if err := ContainerTemplate.Execute(&builder, container); err != nil {
		panic(err)
	}

	return builder.String()
}

// RenderClass to string output
func RenderClass(class *domain.Class) string {
	var builder strings.Builder
	if err := ClassTemplate.Execute(&builder, class); err != nil {
		panic(err)
	}

	return builder.String()
}

// RenderRelation to string output
func RenderRelation(relation *domain.Relation) string {
	var builder strings.Builder
	if err := RelationTemplate.Execute(&builder, relation); err != nil {
		panic(err)
	}

	return builder.String()
}

// NewTemplate from name, input
func NewTemplate(name, input string, funcs ...template.FuncMap) *template.Template {
	tpl := template.New(name)
	tpl.Funcs(template.FuncMap{
		"safe": func(inp string) string {
			re := regexp.MustCompile(`\$`)
			for _, match := range re.FindAllString(inp, -1) {
				inp = strings.Replace(inp, match, `\$`, 1)
			}

			return inp
		},
	})
	for _, fn := range funcs {
		tpl.Funcs(fn)
	}

	return template.Must(tpl.Parse(input))
}
