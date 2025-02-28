package parse

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/Emptyless/jsonschema-transform/internal/domain"
	"github.com/kaptinlin/jsonschema"
)

// ErrUnknownSchema is returned when the schema is not resolved or found
var ErrUnknownSchema = errors.New("unknown schema")

// ClassParser tracks the transformation from jsonschema.Schema to Class
type ClassParser struct {
	*Parser

	// classes stores the last computed result
	classes []*domain.Class

	// relations stores the last computed result
	relations []*domain.Relation

	// queue of remaining jsonschema.Schema's to process
	queue []*jsonschema.Schema

	// references between two jsonschema.Schema's
	references []*Reference
}

// Classes returns the parsed Class slice that can be used by transformations
func (p *Parser) Classes() ([]*domain.Class, error) {
	if p.classParser == nil {
		p.classParser = &ClassParser{
			Parser: p,
		}
	}

	return p.classParser.Classes()
}

// Classes returns the parsed Class slice that can be used by transformations
func (p *ClassParser) Classes() ([]*domain.Class, error) {
	if p.classes != nil {
		return p.classes, nil
	}

	schemas, err := p.Schemas()
	if err != nil {
		return nil, err
	}

	p.queue = append(p.queue, schemas...)
	for len(p.queue) > 0 {
		schema := p.queue[0]
		p.queue = p.queue[1:]

		class, classErr := p.NewClass(schema)
		if classErr != nil {
			return nil, classErr
		}

		// Add a source
		source := schema.GetSchemaURI()
		if p.BaseURI != "" {
			file := regexp.MustCompile("(file:/?/?)/")

			if match := file.FindStringSubmatch(source); len(match) > 0 {
				source = strings.TrimPrefix(source, match[len(match)-1])
			}

			source = p.BaseURI + source
		}

		class.Source = domain.FileSource{FilePath: source}

		p.classes = append(p.classes, class)
	}

	if p.Parser.Depth >= 0 {
		relations, relationsErr := p.Relations()
		if relationsErr != nil {
			return nil, relationsErr
		}

		depthMap := DepthMap(schemas, p.classes, relations)
		classes := []*domain.Class{}
		for k, v := range depthMap {
			if v <= p.Parser.Depth {
				classes = append(classes, k)
			}
		}

		p.classes = classes
		p.relations = nil // reset relations to recalculate
	}

	return p.classes, nil
}

// Relations returns the parsed Reference's between various domain.Class and domain.Property
func (p *Parser) Relations() ([]*domain.Relation, error) {
	c := &ClassParser{
		Parser: p,
	}

	return c.Relations()
}

// Relations returns the parsed Reference's between various domain.Class and domain.Property
func (p *ClassParser) Relations() ([]*domain.Relation, error) {
	if p.relations != nil {
		return p.relations, nil
	}

	// first compute classes if not computed before
	if p.classes == nil {
		_, err := p.Classes()
		if err != nil {
			return nil, err
		}
	}

	for _, reference := range p.references {
		var from *domain.Class
		var to *domain.Class

		for _, class := range p.classes {
			if reference.FromParent == class.Schema {
				from = class
				continue
			}

			if reference.ToParent == class.Schema {
				to = class
				continue
			}

			if from != nil && to != nil {
				break
			}
		}

		if to == nil && p.Depth > -1 {
			continue // reference is too deep and hence filtered from result
		}

		if from == nil || to == nil {
			return nil, fmt.Errorf("one end of the relation '%s'(%s) to '%s' is missing", reference.FromParent.ID, reference.From.Ref, reference.ToParent.ID)
		}

		p.relations = append(p.relations, &domain.Relation{
			Type: "associates",
			From: from,
			To:   to,
		})

	}

	return p.relations, nil
}

// NewClass for a jsonschema.Schema
func (p *ClassParser) NewClass(schema *jsonschema.Schema) (*domain.Class, error) {
	class := domain.Class{
		Schema: schema,
	}

	if title := schema.Title; title != nil {
		class.Name = *title
	}

	if description := schema.Description; description != nil {
		class.Docstring = *description
	}

	if properties := schema.Properties; properties != nil {
		for _, name := range slices.Sorted(maps.Keys(*properties)) {
			value := (*properties)[name]

			property, propertyErr := p.NewProperty(schema, name, value)
			if propertyErr != nil {
				return nil, fmt.Errorf("failed to parse property %s for class %s: %w", name, class.Name, propertyErr)
			}

			property.Parent = &class
			class.Properties = append(class.Properties, property)
		}
	}

	return &class, nil
}

// NewProperty for a Class based on its property jsonschema.Schema
func (p *ClassParser) NewProperty(parent *jsonschema.Schema, name string, value *jsonschema.Schema) (*domain.Property, error) {
	property := domain.Property{
		Name: name,
	}

	if value.Ref != "" || value.DynamicRef != "" {
		resolvedRef, resolvedRefErr := p.PropertyRef(parent, value)
		if resolvedRefErr != nil {
			return nil, resolvedRefErr
		}
		value = resolvedRef
	}

	property.Type = first(value.Type)

	// if Type is "array" (and thus has "items", use the type of "items")
	if items := value.Items; items != nil {
		item, err := p.NewProperty(parent, name, items)
		if err != nil {
			return nil, err
		}

		item.Type = fmt.Sprintf("[]%s", item.Type)
		return item, nil
	}

	// if the value is OneOf; process the OneOf as a oneOf type
	if oneOf := value.OneOf; len(oneOf) > 0 {
		var items []*domain.Property
		for _, schema := range oneOf {
			item, err := p.NewProperty(parent, name, schema)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}

		var types []string
		for _, item := range items {
			types = append(types, item.Type)
		}

		property.Type = fmt.Sprintf("oneOf[%s]", strings.Join(types, ","))
	}

	if format := value.Format; format != nil {
		if property.Type != "" {
			property.Type += fmt.Sprintf("[%s]", *format)
		} else {
			property.Type = *format
		}
	}

	if description := value.Description; description != nil {
		property.Docstring = *description
	}

	return &property, nil
}

// PropertyRef handles properties which are defined with a $ref (possibly with an anchor '#")
func (p *ClassParser) PropertyRef(parent *jsonschema.Schema, value *jsonschema.Schema) (*jsonschema.Schema, error) {
	var resolvedRef *jsonschema.Schema
	if v := value.ResolvedRef; v != nil {
		resolvedRef = v
	} else if v = value.ResolvedDynamicRef; v != nil {
		resolvedRef = v
	} else {
		return nil, fmt.Errorf("$ref '%s' (or $dynamicRef '%s') could not be resolved: %w", value.Ref, value.DynamicRef, ErrUnknownSchema)
	}

	resolvedRefParent, getSchemaErr := p.Parser.Compiler.GetSchema(resolvedRef.GetSchemaURI())
	if getSchemaErr != nil {
		return nil, fmt.Errorf("parent of $ref '%s' failed to load: %w", resolvedRef.GetSchemaURI(), getSchemaErr)
	}

	// if not already processed, process the parent
	if !p.Cache.HasProcessed(resolvedRefParent) {
		// add resolvedRefParent to queue for processing
		p.queue = append(p.queue, resolvedRefParent)
	}

	// add reference to references
	p.references = append(p.references, &Reference{
		Type:       Ref,
		From:       value,
		FromParent: parent,
		To:         resolvedRef,
		ToParent:   resolvedRefParent,
	})

	return resolvedRef, nil
}

// first element of slice
func first[T any](input []T) T {
	if len(input) > 0 {
		return input[0]
	}

	var t T
	return t
}
