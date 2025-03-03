package parse

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kaptinlin/jsonschema"
)

// ErrInvalidGlob is returned when (one of the) glob patterns is invalid
var ErrInvalidGlob = errors.New(`invalid glob pattern`)

// ErrReadFile is returned when a file cannot be
var ErrReadFile = errors.New(`cannot read file`)

// ErrParsingSchema is returned when the jsonschema compiler could not parse the json file
var ErrParsingSchema = errors.New(`could not parse schema`)

// Parser for .json files to execute transforms
type Parser struct {
	// Globs to search for JSON files
	Globs []string

	// StrictMode will error if a JSON file contained in the Globs cannot be parsed to a jsonschema.Schema
	// default false where its logged as a warning
	StrictMode bool

	// BaseURI used for resolving schemas
	BaseURI string

	// Cache implementation to process nested object's and $defs
	Cache Cache

	// Depth of external $refs to follow (or -1 to follow all)
	// 0 implies only referenced schemas
	Depth int

	// Compiler used to load the jsonschema.Schema's
	Compiler *jsonschema.Compiler

	// classParser used for caching intermediate results
	classParser *ClassParser
}

// NewParser for glob patterns, e.g. "*", "**/*.json", ...
func NewParser(globs ...string) *Parser {
	return &Parser{Globs: globs, Depth: -1}
}

// SetBaseURI from which file:// $id's are resolved
func (p *Parser) SetBaseURI(baseURI string) *Parser {
	p.BaseURI = baseURI

	return p
}

// SetDepth to only follow $refs that are 'depth' deep
func (p *Parser) SetDepth(depth int) *Parser {
	p.Depth = depth

	return p
}

// Schemas read by the parser
func (p *Parser) Schemas() ([]*jsonschema.Schema, error) {
	if p == nil || len(p.Globs) == 0 {
		return nil, nil // no-op
	}

	if p.Cache != nil && p.Cache.Schemas() != nil {
		return p.Cache.Schemas(), nil
	}

	if p.Cache == nil {
		p.Cache = &MapCache{}
	}

	if p.Compiler == nil {
		compiler, newCompilerErr := NewCompiler(p.BaseURI)
		if newCompilerErr != nil {
			return nil, newCompilerErr
		}
		p.Compiler = compiler
	}

	var res []*jsonschema.Schema
	for _, glob := range p.Globs {
		logrus.Info("parsing glob pattern: ", glob)
		matches, err := filepath.Glob(glob)
		if err != nil {
			return nil, errors.Join(ErrInvalidGlob, err)
		}

		if len(matches) == 0 {
			logrus.Info("no matches for glob pattern: ", glob)
		}

		for _, match := range matches {
			logrus.Info("parsing file: ", match)
			if !strings.HasSuffix(match, ".json") {
				continue
			}

			schema, readSchemaErr := ReadSchema(p.Compiler, match, p.StrictMode)
			if readSchemaErr != nil {
				return nil, readSchemaErr
			}

			schema, getSchemaErr := p.Compiler.GetSchema(schema.GetSchemaURI())
			if getSchemaErr != nil {
				return nil, getSchemaErr
			}

			res = append(res, schema)
		}
	}

	return res, nil
}

// NewCompiler for baseURI. If the baseURI is an empty string "" the current working directory is used.
func NewCompiler(baseURI string) (*jsonschema.Compiler, error) {
	compiler := jsonschema.NewCompiler()
	var workingDirectory string
	if baseURI != "" {
		compiler = compiler.SetDefaultBaseURI(baseURI)
		workingDirectory = baseURI
	} else {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return nil, cwdErr
		}
		workingDirectory = cwd
	}

	// register both the file:// and the implicit scheme as a file loader
	compiler.RegisterLoader("file", NewFileLoader(workingDirectory))
	compiler.RegisterLoader("", NewFileLoader(workingDirectory))

	return compiler, nil
}

// ReadSchema from file into a jsonschema.Schema
func ReadSchema(compiler *jsonschema.Compiler, filepath string, strict bool) (*jsonschema.Schema, error) {
	contents, readFileErr := os.ReadFile(filepath)
	if readFileErr != nil && strict {
		return nil, errors.Join(ErrReadFile, readFileErr)
	} else if readFileErr != nil {
		logrus.Warnf("could not read file %s", filepath)
		logrus.Debug(readFileErr.Error())
		return nil, nil
	}

	schema, compileSchemaErr := compiler.Compile(contents)
	if compileSchemaErr != nil && strict {
		return nil, errors.Join(ErrParsingSchema, compileSchemaErr)
	} else if compileSchemaErr != nil {
		logrus.Warnf("could not compile schema %s", filepath)
		logrus.Debug(compileSchemaErr.Error())
		return nil, nil
	}

	return schema, nil
}

// Loader used by jsonschema.Compiler::Loaders
type Loader func(url string) (io.ReadCloser, error)

// NewFileLoader constructs a loader that can read from disk
func NewFileLoader(workingDirectory string) Loader {
	workingDirectory = strings.TrimPrefix(workingDirectory, "file://")

	return func(url string) (io.ReadCloser, error) {
		url = strings.TrimPrefix(url, "file://")
		url = path.Join(workingDirectory, url)
		url = path.Clean(url)

		parts := strings.SplitN(url, "#", 2)
		if len(parts) == 2 {
			url = parts[0]
		}

		return os.Open(url)
	}
}
