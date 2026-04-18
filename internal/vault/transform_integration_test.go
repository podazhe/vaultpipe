package vault_test

import (
	"testing"

	"github.com/yourorg/vaultpipe/internal/vault"
)

// TestTransformer_RoundTrip simulates a realistic pipeline:
// fetch secrets, filter to relevant keys, prefix for a service namespace.
func TestTransformer_RoundTrip(t *testing.T) {
	input := map[string]string{
		"db_password": "s3cr3t",
		"db_user":     "admin",
		"api_key":     "key123",
		"debug":       "true",
	}

	tr := vault.NewTransformer(vault.TransformRule{
		Prefix:  "myapp",
		Filter:  []string{"db_password", "db_user", "api_key"},
		Renames: map[string]string{"api_key": "API_TOKEN"},
	})

	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := map[string]string{
		"MYAPP_db_password": "s3cr3t",
		"MYAPP_db_user":     "admin",
		"MYAPP_API_TOKEN":   "key123",
	}
	for k, want := range cases {
		if got := out[k]; got != want {
			t.Errorf("key %s: want %q got %q", k, want, got)
		}
	}
	if _, ok := out["debug"]; ok {
		t.Error("debug should have been filtered out")
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}
