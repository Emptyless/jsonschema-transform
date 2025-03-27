package parse

import (
	"maps"
	"reflect"
	"slices"

	"github.com/kaptinlin/jsonschema"
)

// Cache tracks processed jsonschema.Schema's
type Cache interface {
	// Process a schema
	Process(schema *jsonschema.Schema)
	// HasProcessed a schema
	HasProcessed(schema *jsonschema.Schema) bool
	// Schemas stored in Cache
	Schemas() []*jsonschema.Schema
}

// MapCache tracks processed jsonschema.Schema's
type MapCache struct {
	processed map[*jsonschema.Schema]struct{}
}

// Schemas stored in Cache
func (c *MapCache) Schemas() []*jsonschema.Schema {
	if c == nil || c.processed == nil {
		return nil
	}

	return slices.Collect(maps.Keys(c.processed))
}

// Process marks a jsonschema.Schema as processed
func (c *MapCache) Process(schema *jsonschema.Schema) {
	if c.processed == nil {
		c.processed = make(map[*jsonschema.Schema]struct{})
	}

	c.processed[schema] = struct{}{}
	return
}

// HasProcessed returns true iff a jsonschema.Schema is already processed
func (c *MapCache) HasProcessed(schema *jsonschema.Schema) bool {
	if c == nil || len(c.processed) == 0 {
		return false
	}

	// if pointer equivalent, return
	if _, ok := c.processed[schema]; ok {
		return true
	}

	// deep equal also return true
	for s, _ := range c.processed {
		if reflect.DeepEqual(s, schema) {
			return true
		}
	}

	return false
}
