package vault

import (
	"strings"
	"testing"
)

var exportFixture = map[string]string{
	"DB_PASS": "s3cr3t",
	"API_KEY": "abc123",
}

func TestExportSecrets_Dotenv(t *testing.T) {
	out, err := ExportSecrets(exportFixture, ExportOptions{Format: FormatDotenv})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `API_KEY="abc123"`) {
		t.Errorf("missing API_KEY line, got:\n%s", out)
	}
	if !strings.Contains(out, `DB_PASS="s3cr3t"`) {
		t.Errorf("missing DB_PASS line, got:\n%s", out)
	}
}

func TestExportSecrets_DotenvExport(t *testing.T) {
	out, err := ExportSecrets(exportFixture, ExportOptions{Format: FormatDotenv, Export: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "export API_KEY") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestExportSecrets_JSON(t *testing.T) {
	out, err := ExportSecrets(exportFixture, ExportOptions{Format: FormatJSON})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
	if !strings.Contains(out, `"API_KEY": "abc123"`) {
		t.Errorf("missing API_KEY in JSON, got:\n%s", out)
	}
}

func TestExportSecrets_YAML(t *testing.T) {
	out, err := ExportSecrets(exportFixture, ExportOptions{Format: FormatYAML})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `API_KEY: "abc123"`) {
		t.Errorf("missing API_KEY in YAML, got:\n%s", out)
	}
}

func TestExportSecrets_Prefix(t *testing.T) {
	out, err := ExportSecrets(exportFixture, ExportOptions{Format: FormatDotenv, Prefix: "APP_"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "APP_API_KEY") {
		t.Errorf("expected prefixed key, got:\n%s", out)
	}
}

func TestExportSecrets_InvalidFormat(t *testing.T) {
	_, err := ExportSecrets(exportFixture, ExportOptions{Format: "toml"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportSecrets_SortedOutput(t *testing.T) {
	out, _ := ExportSecrets(exportFixture, ExportOptions{Format: FormatDotenv})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "API_KEY") {
		t.Errorf("expected API_KEY first (sorted), got: %s", lines[0])
	}
}
