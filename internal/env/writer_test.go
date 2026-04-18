package env

import (
	"os"
	"strings"
	"testing"
)

func TestSanitiseKey(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"my-secret", "MY_SECRET"},
		{"db/password", "DB_PASSWORD"},
		{"api.key", "API_KEY"},
	}
	for _, c := range cases {
		got := sanitiseKey(c.in)
		if got != c.want {
			t.Errorf("sanitiseKey(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestWriter_Dotenv(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	w, err := NewWriter(tmp.Name(), FormatDotenv)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	defer w.Close()

	secrets := map[string]string{"db-password": "s3cr3t"}
	if err := w.Write(secrets); err != nil {
		t.Fatalf("Write: %v", err)
	}
	w.Close()

	data, _ := os.ReadFile(tmp.Name())
	if !strings.Contains(string(data), "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in output, got: %s", data)
	}
}

func TestWriter_Export(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	w, err := NewWriter(tmp.Name(), FormatExport)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	defer w.Close()

	if err := w.Write(map[string]string{"api.key": "abc123"}); err != nil {
		t.Fatalf("Write: %v", err)
	}
	w.Close()

	data, _ := os.ReadFile(tmp.Name())
	if !strings.HasPrefix(string(data), "export ") {
		t.Errorf("expected 'export ' prefix, got: %s", data)
	}
}
