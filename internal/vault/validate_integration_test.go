package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newValidateTestServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		})
	}))
}

func TestValidate_Integration_AllValid(t *testing.T) {
	srv := newValidateTestServer(t, map[string]interface{}{"DB_HOST": "localhost", "DB_PORT": "5432"})
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	secrets, err := ReadSecrets(client, "secret/data/app")
	if err != nil {
		t.Fatalf("ReadSecrets: %v", err)
	}
	res := ValidateSecrets(secrets, ValidateOptions{
		RequiredKeys: []string{"DB_HOST", "DB_PORT"},
		NoEmpty:      true,
	})
	if !res.OK() {
		t.Fatalf("expected validation to pass: %v", res.Errors)
	}
}

func TestValidate_Integration_MissingRequired(t *testing.T) {
	srv := newValidateTestServer(t, map[string]interface{}{"DB_HOST": "localhost"})
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	secrets, err := ReadSecrets(client, "secret/data/app")
	if err != nil {
		t.Fatalf("ReadSecrets: %v", err)
	}
	res := ValidateSecrets(secrets, ValidateOptions{
		RequiredKeys: []string{"DB_HOST", "DB_PASS"},
	})
	if res.OK() {
		t.Fatal("expected validation to fail for missing DB_PASS")
	}
}
