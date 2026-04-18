package template

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderer_BasicInterpolation(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf)

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	tmplSrc := `host={{ index . "DB_HOST" }} port={{ index . "DB_PORT" }}`
	if err := r.Render(tmplSrc, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "localhost") || !strings.Contains(got, "5432") {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestRenderer_SecretHelper(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf)

	secrets := map[string]string{"API_KEY": "abc123"}
	tmplSrc := `key={{ secret "API_KEY" }}`

	if err := r.Render(tmplSrc, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "key=abc123" {
		t.Errorf("got %q, want %q", got, "key=abc123")
	}
}

func TestRenderer_RequiredMissing(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf)

	secrets := map[string]string{"PRESENT": "yes"}
	tmplSrc := `{{ required "MISSING" (secret "MISSING") }}`

	if err := r.Render(tmplSrc, secrets); err == nil {
		t.Fatal("expected error for missing required secret, got nil")
	}
}

func TestRenderer_InvalidTemplate(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf)

	if err := r.Render(`{{ .Unclosed`, nil); err == nil {
		t.Fatal("expected parse error, got nil")
	}
}
