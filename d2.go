package main

import (
	"fmt"
	"os"
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
	tool.Apply(d2Cmd.Flags())
}

// handleD2 for the d2Cmd command
func handleD2(cmd *cobra.Command, _ []string) error {
	globs := cmd.Flag(globsFlag.Name).Value.(pflag.SliceValue).GetSlice()
	if len(globs) == 0 {
		return ErrNoGlobs
	}

	parser := parse.NewParser(globs...)
	if baseURI := cmd.Flag(baseURIFlag.Name).Value.String(); baseURI != "" && (!strings.HasPrefix(baseURI, "http") || strings.HasPrefix(baseURI, "https")) {
		baseURIAbs, err := filepath.Abs(baseURI)
		if err != nil {
			logrus.Error("unable to determine absolute filepath", err)
			return err
		}
		parser.SetBaseURI("file://" + baseURIAbs)
	} else if baseURI != "" && (strings.HasPrefix(baseURI, "http") || strings.HasPrefix(baseURI, "https")) {
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

	output, err := d2.D2(parser, &d2.Config{Format: format, Tool: cmd.Flag(tool.Name).Value.String()})
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
