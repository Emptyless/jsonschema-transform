package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestD2_CreatesD2File(t *testing.T) {
	// Arrange
	outputBuffer := new(bytes.Buffer)
	rootCmd.SetOut(outputBuffer)
	args := []string{d2Cmd.Use, "--globs", "./testdata/*.json", "--base-uri", "./", "--overwrite", "--output", "diagram.d2"}
	rootCmd.SetArgs(args)

	// Act
	err := rootCmd.Execute()

	// Assert
	require.NoError(t, err)
}

func TestD2_CreatesSvg(t *testing.T) {
	// Arrange
	outputBuffer := new(bytes.Buffer)
	rootCmd.SetOut(outputBuffer)
	args := []string{d2Cmd.Use, "--globs", "./testdata/*.json", "--base-uri", "./", "--overwrite", "--output", "diagram.svg"}
	rootCmd.SetArgs(args)

	// Act
	err := rootCmd.Execute()

	// Assert
	require.NoError(t, err)
}

func TestD2_CreatesContainerizedD2(t *testing.T) {
	// Arrange
	outputBuffer := new(bytes.Buffer)
	rootCmd.SetOut(outputBuffer)
	args := []string{d2Cmd.Use, "--globs", "./testdata/*.json", "--base-uri", "./", "--container-base-path", ".", "--overwrite", "--output", "diagram_with_containers.d2"}
	rootCmd.SetArgs(args)

	// Act
	err := rootCmd.Execute()

	// Assert
	require.NoError(t, err)
}

func TestD2_CreatesContainerizedSvg(t *testing.T) {
	// Arrange
	outputBuffer := new(bytes.Buffer)
	rootCmd.SetOut(outputBuffer)
	args := []string{d2Cmd.Use, "--globs", "./testdata/*.json", "--base-uri", "./", "--container-base-path", ".", "--overwrite", "--output", "diagram_with_containers.svg"}
	rootCmd.SetArgs(args)

	// Act
	err := rootCmd.Execute()

	// Assert
	require.NoError(t, err)
}
