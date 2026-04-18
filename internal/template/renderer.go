package template

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

// Renderer renders a template string using secrets as the data source.
type Renderer struct {
	out io.Writer
}

// NewRenderer creates a Renderer that writes to out.
func NewRenderer(out io.Writer) *Renderer {
	return &Renderer{out: out}
}

// Render parses tmplSrc as a Go text/template and executes it with secrets
// as the data map. The result is written to the underlying writer.
func (r *Renderer) Render(tmplSrc string, secrets map[string]string) error {
	funcMap := template.FuncMap{
		"required": func(key string, val string) (string, error) {
			if val == "" {
				return "", fmt.Errorf("required secret %q is empty", key)
			}
			return val, nil
		},
		"secret": func(key string) string {
			return secrets[key]
		},
	}

	tmpl, err := template.New("vaultpipe").Funcs(funcMap).Parse(tmplSrc)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, secrets); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	_, err = r.out.Write(buf.Bytes())
	return err
}
