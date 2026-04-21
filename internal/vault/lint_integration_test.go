package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vaultpipe/vaultpipe/internal/vault"
)

func newLintTestServer(data map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		envelope := map[string]interface{}{
			"data": map[string]interface{}{
				"data": data,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(envelope)
	}))
}

func TestLint_Integration_DetectsViolations(t *testing.T) {
	server := newLintTestServer(map[string]string{
		"bad_key": "",
		"GOOD_KEY": "value",
	})
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("client: %v", err)
	}

	secrets, err := vault.ReadSecrets(client, "secret/data/test")
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	results := vault.LintSecrets(secrets, vault.LintOptions{AllowEmpty: false})

	rulesSeen := map[string]bool{}
	for _, r := range results {
		rulesSeen[r.Rule] = true
	}

	if !rulesSeen["no-lowercase-key"] {
		t.Error("expected no-lowercase-key violation for bad_key")
	}
	if !rulesSeen["no-empty-value"] {
		t.Error("expected no-empty-value violation for bad_key")
	}
}

func TestLint_Integration_CleanSecrets(t *testing.T) {
	server := newLintTestServer(map[string]string{
		"API_KEY":  "abc123",
		"DB_PASS":  "secure",
	})
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("client: %v", err)
	}

	secrets, err := vault.ReadSecrets(client, "secret/data/clean")
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	results := vault.LintSecrets(secrets, vault.LintOptions{AllowEmpty: false})
	if len(results) != 0 {
		t.Errorf("expected 0 violations, got %d", len(results))
	}
}
