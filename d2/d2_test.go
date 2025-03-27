package d2

import (
	"os"
	"testing"

	"github.com/Emptyless/jsonschema-transform/domain"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

func TestD2(t *testing.T) {
	// Arrange
	parser := TestParser{
		ClassData: []*domain.Class{
			{
				Source:    domain.FileSource{FilePath: "pet.json"},
				Name:      "Pet",
				Docstring: "",
				Properties: []*domain.Property{
					{Name: "id", Type: "string"},
					{Name: "name", Type: "string"},
				},
			},
		},
	}

	// Act
	b, err := D2(&parser, &Config{Format: SVG})
	_ = os.WriteFile("out.svg", b, 0o644)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, b)
}

type TestParser struct {
	ClassData    []*domain.Class
	ClassesError error

	RelationsData  []*domain.Relation
	RelationsError error
}

func (p *TestParser) Classes() ([]*domain.Class, error) {
	if err := p.ClassesError; err != nil {
		return nil, err
	}

	return p.ClassData, nil
}

func (p *TestParser) Relations() ([]*domain.Relation, error) {
	if err := p.RelationsError; err != nil {
		return nil, err
	}

	return p.RelationsData, nil
}
