package main

import (
	"errors"

	"github.com/spf13/pflag"
)

// flag is used as a common flag representation where flags can be specified that are re-used across commands.
// See pflag.StringP for details about these attributes
type flag struct {
	Name  string
	Short string
	Value any
	Usage string
}

// Apply to pflag.FlagSet basic on type of flag.Value or panic if unknown type
func (f *flag) Apply(flagSet *pflag.FlagSet) {
	switch f.Value.(type) {
	case string:
		flagSet.StringP(f.Name, f.Short, f.Value.(string), f.Usage)
	case []string:
		flagSet.StringSliceP(f.Name, f.Short, f.Value.([]string), f.Usage)
	case bool:
		flagSet.BoolP(f.Name, f.Short, f.Value.(bool), f.Usage)
	default:
		panic("unknown type")
	}
}

var outputFlag = flag{
	Name:  "output",
	Short: "o",
	Value: "diagram.%s", // first %s is substituted by the default extension
	Usage: "Optionally set the location of the output file",
}

var globsFlag = flag{
	Name:  "globs",
	Short: "g",
	Value: []string{},
	Usage: "glob patterns to match (e.g. '**/*.json', '*.json', 'file.json')",
}

var baseURIFlag = flag{
	Name:  "base-uri",
	Short: "",
	Value: "",
	Usage: "when provided the base-uri will be used to resolve json schemas using the $id property. If the schema does not start with 'http' or 'https' it's assumed to be a 'file://' reference",
}

var allowOverwriteFlag = flag{
	Name:  "overwrite",
	Short: "",
	Value: false,
	Usage: "if provided allows existing files to be overwritten (i.e. regenerate)",
}

var toolFlag = flag{
	Name:  "tool",
	Short: "",
	Value: "",
	Usage: "path to the tool to render the 'd2' file, if empty 'which d2' is used",
}

var containerBasePathFlag = flag{
	Name:  "container-base-path",
	Short: "",
	Value: "",
	Usage: "can only be used in conjunction with a file based --base-uri, if set will put classes in containers representing directories",
}

// ErrNoGlobs is returned when no globs are provided (which is a no-op)
var ErrNoGlobs = errors.New("no globs provided")

// ErrNoOverwrite is returned when a file would be overwritten which is not allowed
var ErrNoOverwrite = errors.New("file exists but overwrite of file is not allowed")
