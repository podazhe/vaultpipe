package env

import (
	"fmt"
	"os"
	"strings"
)

// Format represents the output format for env files.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatExport Format = "export"
)

// Writer writes secrets to an env file or stdout.
type Writer struct {
	format Format
	output *os.File
}

// NewWriter creates a Writer targeting the given file path.
// Pass an empty path to write to stdout.
func NewWriter(path string, format Format) (*Writer, error) {
	var out *os.File
	if path == "" {
		out = os.Stdout
	} else {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return nil, fmt.Errorf("opening env file: %w", err)
		}
		out = f
	}
	return &Writer{format: format, output: out}, nil
}

// Write serialises the provided key/value pairs according to the chosen format.
func (w *Writer) Write(secrets map[string]string) error {
	for k, v := range secrets {
		key := sanitiseKey(k)
		escaped := escapeValue(v)
		var line string
		switch w.format {
		case FormatExport:
			line = fmt.Sprintf("export %s=%q\n", key, escaped)
		default: // dotenv
			line = fmt.Sprintf("%s=%q\n", key, escaped)
		}
		if _, err := fmt.Fprint(w.output, line); err != nil {
			return fmt.Errorf("writing secret %s: %w", key, err)
		}
	}
	return nil
}

// Close closes the underlying file if it is not stdout.
func (w *Writer) Close() error {
	if w.output != os.Stdout {
		return w.output.Close()
	}
	return nil
}

func sanitiseKey(k string) string {
	return strings.ToUpper(strings.NewReplacer("-", "_", "/", "_", ".", "_").Replace(k))
}

func escapeValue(v string) string {
	return v
}
