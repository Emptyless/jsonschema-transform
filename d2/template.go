package d2

import (
	"html/template"
	"strings"

	"github.com/Emptyless/jsonschema-transform/internal/domain"
)

var ClassTemplate = NewTemplate("Class", `
{{- $.Name }}: {
  shape: class
{{- range $property := $.Properties }}
  "{{ $property.Name }}": "{{ $property.Type }}"
{{- end }}
}
`)

var RelationTemplate = NewTemplate("Relation", `
{{- $.From.Name }} -- {{ $.To.Name }}: {{ $.Type }}`)

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
func NewTemplate(name, input string) *template.Template {
	return template.Must(template.New(name).Parse(input))
}
