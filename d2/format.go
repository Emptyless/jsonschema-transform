package d2

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// ErrUnknownFormat is returned when the supplied format is not recognized
var ErrUnknownFormat = errors.New("unknown format")

// Format of output string supported by D2
type Format string

// SVG format
const SVG Format = "svg"

// Native D2 script format
const Native Format = "d2"

// FormatFromFile parses the given path and determines the output format used by the D2 parser
func FormatFromFile(path string) (Format, error) {
	ext := filepath.Ext(path)
	if ext == "" {
		ext = "." + path // not a filepath but just the extension, implicitly use the path for simplicity
	}

	switch ext {
	case ".svg":
		return SVG, nil
	case ".d2":
		return Native, nil
	default:
		return "", ErrUnknownFormat
	}
}

// Render Format to a byte slice
func (f Format) Render(buffer *bytes.Buffer, cfg *Config) ([]byte, error) {
	switch f {
	case Native:
		return buffer.Bytes(), nil
	case SVG:
		// create temporary diagram file
		diagramFile, createDiagramFile := os.CreateTemp("", "diagram-*.d2")
		if createDiagramFile != nil {
			return nil, createDiagramFile
		}

		// write the diagram code to the file
		if _, err := diagramFile.Write(buffer.Bytes()); err != nil {
			return nil, err
		}

		// create temporary svg file
		svgFile, createSvgFileErr := os.CreateTemp("", "diagram-*.svg")
		if createSvgFileErr != nil {
			return nil, createSvgFileErr
		} else if closeErr := svgFile.Close(); closeErr != nil {
			return nil, closeErr
		}

		output, outputErr := exec.Command(cfg.Tool, diagramFile.Name(), svgFile.Name()).Output()
		if err := new(exec.ExitError); errors.As(outputErr, &err) {
			var builder strings.Builder
			for i, line := range strings.Split(buffer.String(), "\n") {
				builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, line))
			}
			logrus.Debug("diagram.d2:\n")
			logrus.Debug(builder.String())

			return nil, errors.New(string(err.Stderr))
		} else if outputErr != nil {
			return nil, outputErr
		}

		logrus.Debug(string(output))

		svgFile, openSvgFileErr := os.Open(svgFile.Name())
		if openSvgFileErr != nil {
			return nil, openSvgFileErr
		}

		return io.ReadAll(svgFile)
	default:
		return nil, ErrUnknownFormat
	}
}
