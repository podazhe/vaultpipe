package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockVaultServer(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Errorf("mock server encode error: %v", err)
		}
	}))
}

func TestNewClient_TokenFromConfig(t *testing.T) {
	srv := mockVaultServer(t, "", map[string]interface{}{})
	defer srv.Close()

	client, err := NewClient(Config{
		Address: srv.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_NoAuth(t *testing.T) {
	srv := mockVaultServer(t, "", map[string]interface{}{})
	defer srv.Close()

	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{Address: srv.URL})
	if err == nil {
		t.Fatal("expected error when no auth provided")
	}
}

func TestReadSecrets_KVv2(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"DB_PASSWORD": "supersecret",
				"API_KEY":     "abc123",
			},
		},
	}
	srv := mockVaultServer(t, "/v1/secret/data/myapp", payload)
	defer srv.Close()

	client, err := NewClient(Config{Address: srv.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("client creation failed: %v", err)
	}

	secrets, err := client.ReadSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("ReadSecrets failed: %v", err)
	}

	if secrets["DB_PASSWORD"] != "supersecret" {
		t.Errorf("expected DB_PASSWORD=supersecret, got %q", secrets["DB_PASSWORD"])
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
}
