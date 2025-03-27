package parse

import "github.com/kaptinlin/jsonschema"

// Reference between two jsonschema.Schema's
type Reference struct {
	// ReferenceType can be oneOf, $ref, allOf, etc
	Type ReferenceType

	// From which jsonschema.Schema the relation started, i.e. with $ref
	From       *jsonschema.Schema
	FromParent *jsonschema.Schema

	// To which jsonschema.Schema the relation points (which could be a jsonschema.Schema or $anchor)
	To       *jsonschema.Schema
	ToParent *jsonschema.Schema
}

type ReferenceType string

// Ref between schemas using $ref
const Ref ReferenceType = "$ref"

// OneOf between schemas using oneOf
const OneOf ReferenceType = "oneOf"

// AllOf between schemas using oneOf
const AllOf ReferenceType = "allOf"

// AnyOf between schemas using oneOf
const AnyOf ReferenceType = "anyOf"
