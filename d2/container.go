package d2

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Emptyless/jsonschema-transform/internal/domain"
)

// ContainerParser derives the nesting of containers for a particular domain.Source
// e.g. if 'some,nested,value' is returned the eventual class is rendered inside a three layer deep Container
type ContainerParser interface {
	Containers(source domain.Source) []string
}

// DirContainerParser uses a base directory to parse containers by stripping the RootPath from the domain.Source
// of the domain.Class and using the directory representation as containers
type DirContainerParser struct {
	// RootPath to remove from the domain.Source to determine directory containers
	RootPath string
}

// Containers separated by the filepath.Separator in the filesystem relative to some RootPath
// if the source contains the file:// protocol, it's removed automatically
func (d DirContainerParser) Containers(source domain.Source) []string {
	path := source.Path()
	dir := filepath.Dir(path)
	file := regexp.MustCompile("(file:/?/?)/")

	if match := file.FindStringSubmatch(dir); len(match) > 0 {
		dir = strings.TrimPrefix(dir, match[len(match)-1])
	}

	dir = strings.TrimPrefix(dir, d.RootPath)
	dir = strings.TrimPrefix(dir, "/")
	dir = strings.TrimSuffix(dir, "/")

	return strings.Split(dir, string(filepath.Separator))
}

// Containers is a helper type to perform methods on a group of Container
type Containers []*Container

// Contains a Container with Container.Name
// returns true and the Container if found or false otherwise
func (c Containers) Contains(name string) (*Container, bool) {
	for _, container := range c {
		if container.Name == name {
			return container, true
		}
	}

	return nil, false
}

// Container to put domain.Class in as a grouping
type Container struct {
	// Name of the container
	Name string

	// Classes in the container
	Classes []*domain.Class

	// Containers nested in this Container
	Containers Containers
}

// Add a domain.Class to the Container by parsing its container names using a ContainerParser
func (c *Container) Add(class *domain.Class, parser ContainerParser) {
	parts := parser.Containers(class.Source)

	container := c
	for _, part := range parts {
		match, ok := container.Containers.Contains(part)
		if !ok {
			match = &Container{Name: part}
			container.Containers = append(container.Containers, match)
		}

		container = match // replace container for next part
	}

	container.Classes = append(container.Classes, class)
}

// Render a Container using RenderContainer. This avoids the circular dependency problem my adding
// the RenderContainer as a function to the ContainerTemplate
func (c *Container) Render() string {
	return RenderContainer(c)
}
