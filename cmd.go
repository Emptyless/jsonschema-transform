package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jsonschema-transform",
	Short: "Welcome to the JSON Schema Transformer CLI",
	Long: `Welcome to the JSON Schema Transformer CLI

The goal of this transformer is to generate diagrams from JSON schema files. This helps to visualize the overall schema`,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		verboseFlag, _ := cmd.Flags().GetCount("verbose")
		quietFlag, _ := cmd.Flags().GetCount("quiet")

		// We default to only showing warnings and errors, but it can be increased to info and debug, and
		// decreased to error and fatal
		logLevel := logrus.WarnLevel + logrus.Level(verboseFlag) - logrus.Level(quietFlag)
		logrus.SetLevel(logLevel)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().CountP("verbose", "v", "Increase the verbosity of the output by one level, -v shows informational logs and -vv will output debug information.")
	rootCmd.PersistentFlags().CountP("quiet", "q", "Decrease the verbosity of the output by one level, -v hides warning logs and -vv will suppress non-fatal errors")
}
