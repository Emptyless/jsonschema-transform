package domain

import "github.com/kaptinlin/jsonschema"

// Class representation of parsed source files
type Class struct {
	Source
	// Schema from which the Class is parsed
	Schema *jsonschema.Schema

	// Name of the Class
	Name string

	// Docstring of the Class
	Docstring string

	// Properties of the Class
	Properties []*Property
}
