package vault

import (
	"testing"
)

func TestFlattenSecrets_Simple(t *testing.T) {
	nested := map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
	}
	got, err := FlattenSecrets(nested, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %q", got["db_host"])
	}
	if got["db_port"] != "5432" {
		t.Errorf("expected db_port=5432, got %q", got["db_port"])
	}
}

func TestFlattenSecrets_DeepNesting(t *testing.T) {
	nested := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	got, err := FlattenSecrets(nested, FlattenOptions{Separator: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["a.b.c"] != "deep" {
		t.Errorf("expected a.b.c=deep, got %q", got["a.b.c"])
	}
}

func TestFlattenSecrets_Prefix(t *testing.T) {
	nested := map[string]any{"key": "val"}
	got, err := FlattenSecrets(nested, FlattenOptions{Prefix: "APP"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["APP_key"] != "val" {
		t.Errorf("expected APP_key=val, got %v", got)
	}
}

func TestFlattenSecrets_UpperCase(t *testing.T) {
	nested := map[string]any{"db": map[string]any{"pass": "secret"}}
	got, err := FlattenSecrets(nested, FlattenOptions{UpperCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %v", got)
	}
}

func TestFlattenSecrets_NilValue(t *testing.T) {
	nested := map[string]any{"empty": nil}
	got, err := FlattenSecrets(nested, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["empty"] != "" {
		t.Errorf("expected empty string for nil value, got %q", got["empty"])
	}
}

func TestFlattenSecrets_NilInput(t *testing.T) {
	_, err := FlattenSecrets(nil, FlattenOptions{})
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestFlattenSecrets_NonStringScalar(t *testing.T) {
	nested := map[string]any{"count": 42}
	got, err := FlattenSecrets(nested, FlattenOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["count"] != "42" {
		t.Errorf("expected count=42, got %q", got["count"])
	}
}
