package parse

import "github.com/kaptinlin/jsonschema"

// Reference between two jsonschema.Schema's
type Reference struct {
	// From which jsonschema.Schema the relation started, i.e. with $ref
	From       *jsonschema.Schema
	FromParent *jsonschema.Schema

	// To which jsonschema.Schema the relation points (which could be a jsonschema.Schema or $anchor)
	To       *jsonschema.Schema
	ToParent *jsonschema.Schema
}
