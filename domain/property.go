package domain

import "github.com/kaptinlin/jsonschema"

// Property of a Class
type Property struct {
	// Schema from which the property is parsed
	Schema *jsonschema.Schema

	// Parent Class the Property belongs to
	Parent *Class

	// Name of the Property
	Name string

	// Type of the Property
	Type string

	// Docstring of the Property
	Docstring string
}
