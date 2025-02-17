package d2

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/Emptyless/jsonschema-transform/internal/domain"
)

// ErrParserFailure is returned when there is some failure by the parser
var ErrParserFailure = errors.New("failed to parse classes or relations")

// ErrNoClasses is returned when no Classes are parsed
var ErrNoClasses = errors.New("no classes parsed")

// Parser implementation that returns a parse.Class slice
type Parser interface {
	Classes() ([]*domain.Class, error)
	Relations() ([]*domain.Relation, error)
}

// Config used when rendering D2
type Config struct {
	// Format used when outputting diagram
	Format Format

	// Tool used to render svg (if SVG Format)
	Tool string

	// Args used by the tool
	Args []string

	// ContainerBasePath if set will wrap all domain.Class in containers based on the
	// DirContainerParser
	ContainerBasePath string
}

// D2 transform Parser with Config into .d2 or .svg output
func D2(parser Parser, cfg *Config) ([]byte, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.Format == "" {
		cfg.Format = Native
	}

	if cfg.Args == nil {
		cfg.Args = []string{}
	}

	if (cfg.Format == SVG || cfg.Format == PNG) && cfg.Tool == "" {
		output, outputErr := exec.Command("which", "d2").Output()
		if outputErr != nil {
			return nil, outputErr
		}

		cfg.Tool = strings.TrimSpace(string(output))
	}

	classes, err := parser.Classes()
	if err != nil {
		return nil, errors.Join(ErrParserFailure, err)
	} else if len(classes) == 0 {
		return nil, ErrNoClasses
	}

	relations, err := parser.Relations()
	if err != nil {
		return nil, errors.Join(ErrParserFailure, err)
	}

	buffer := new(bytes.Buffer)

	// if the ContainerBasePath is set, render containerized
	if cfg.ContainerBasePath != "" {
		containerParser := DirContainerParser{RootPath: cfg.ContainerBasePath}
		container := Container{Name: ""}

		// add classes to the container
		for _, c := range classes {
			container.Add(c, containerParser)
		}

		buffer.WriteString(container.Render())

		// update the name such that the relations can be created
		for _, r := range relations {
			from := containerParser.Containers(r.From)
			prefix := strings.Join(from, ".")
			if !strings.HasPrefix(r.From.Name, prefix) {
				r.From.Name = prefix + "." + r.From.Name
			}

			to := containerParser.Containers(r.To)
			prefix = strings.Join(to, ".")
			if !strings.HasPrefix(r.To.Name, prefix) {
				r.To.Name = prefix + "." + r.To.Name
			}
		}
	} else {
		for _, c := range classes {
			buffer.WriteString(RenderClass(c))
			buffer.WriteString("\n\n")
		}
	}

	for i, r := range relations {
		buffer.WriteString(RenderRelation(r))
		if i != len(relations)-1 {
			buffer.WriteString("\n\n")
		}
	}

	return cfg.Format.Render(buffer, cfg)
}
