package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Emptyless/jsonschema-transform/d2"
	"github.com/Emptyless/jsonschema-transform/internal/parse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// d2Cmd registered to the rootCmd
var d2Cmd = &cobra.Command{
	Use:          "d2",
	Short:        "generate a d2 diagram from the json schemas",
	Long:         "generate a d2 diagram from the json schemas",
	Example:      fmt.Sprintf("%s d2", rootCmd.Use),
	SilenceUsage: true,
	RunE:         handleD2,
}

// init the d2Cmd command
func init() {
	rootCmd.AddCommand(d2Cmd)
	outputFlag.Apply(d2Cmd.Flags())
	globsFlag.Apply(d2Cmd.Flags())
	baseURIFlag.Apply(d2Cmd.Flags())
	allowOverwriteFlag.Apply(d2Cmd.Flags())
	toolFlag.Apply(d2Cmd.Flags())
	containerBasePathFlag.Apply(d2Cmd.Flags())
	depthFlag.Apply(d2Cmd.Flags())
	d2Cmd.Flags().StringP("", "", "", "additional args passed to the D2 (e.g. jsonschema-transform d2 --globs schema.json -- --layout elk")
}

// handleD2 for the d2Cmd command
func handleD2(cmd *cobra.Command, _ []string) error {
	globs := cmd.Flag(globsFlag.Name).Value.(pflag.SliceValue).GetSlice()
	if len(globs) == 0 {
		return ErrNoGlobs
	}

	containerBasePath := cmd.Flag(containerBasePathFlag.Name).Value.String()
	baseURI := cmd.Flag(baseURIFlag.Name).Value.String()
	if containerBasePath != "" && HasHTTPPrefix(baseURI) {
		return fmt.Errorf("cannot use --%s when --%s is not a file:// based URI", containerBasePathFlag.Name, baseURIFlag.Name)
	} else if containerBasePath != "" && baseURI == "" {
		return fmt.Errorf("cannot use --%s when --%s is empty", containerBasePathFlag.Name, baseURIFlag.Name)
	}

	depth, err := cmd.Flags().GetInt(depthFlag.Name)
	if err != nil {
		return err
	}

	parser := parse.NewParser(globs...).SetDepth(depth)
	if baseURI != "" && !HasHTTPPrefix(baseURI) {
		baseURIAbs, err := filepath.Abs(baseURI)
		if err != nil {
			logrus.Error("unable to determine absolute filepath", err)
			return err
		}

		// set parser with absolute baseURI
		parser.SetBaseURI("file://" + baseURIAbs)

		// resolve containerBasePath relative to baseURIAbs
		if containerBasePath != "" {
			containerBasePath = path.Join(baseURIAbs, containerBasePath)
		}
	} else if baseURI != "" && HasHTTPPrefix(baseURI) {
		parser.SetBaseURI(baseURI)
	}

	outputFile := cmd.Flag(outputFlag.Name).Value.String()
	if strings.Contains(outputFile, "%s") {
		outputFile = strings.ReplaceAll(outputFile, "%s", "d2") // replace variable type with d2 extension
	}

	format, err := d2.FormatFromFile(outputFile)
	if err != nil {
		return err
	}

	output, err := d2.D2(parser, &d2.Config{
		Format:            format,
		Tool:              cmd.Flag(toolFlag.Name).Value.String(),
		Args:              cmd.Flags().Args(),
		ContainerBasePath: containerBasePath,
	})
	if err != nil {
		return err
	}

	if _, statErr := os.Stat(outputFile); statErr == nil && cmd.Flag(allowOverwriteFlag.Name).Value.String() == "false" {
		return ErrNoOverwrite
	}

	if writeFileErr := os.WriteFile(outputFile, output, 0o644); writeFileErr != nil {
		return writeFileErr
	}

	logrus.Info("d2 diagram written to ", outputFile)

	return nil
}

// HasHTTPPrefix checks if the baseURI starts with either http or https
func HasHTTPPrefix(baseURI string) bool {
	return baseURI != "" && (strings.HasPrefix(baseURI, "http") || strings.HasPrefix(baseURI, "https"))
}
