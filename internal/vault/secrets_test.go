package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func testClientWithServer(t *testing.T, mux *http.ServeMux) *Client {
	t.Helper()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	raw, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("new vault client: %v", err)
	}
	raw.SetToken("test-token")
	return &Client{logical: raw.Logical()}
}

func TestReadSecrets_KVv2Envelope(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/secret/data/myapp/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"DB_PASSWORD": "s3cr3t",
					"API_KEY":     "abc123",
				},
			},
		})
	})

	c := testClientWithServer(t, mux)
	secrets, err := c.ReadSecrets([]string{"secret/myapp/config"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected DB_PASSWORD=s3cr3t, got %q", secrets["DB_PASSWORD"])
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
}

func TestReadSecrets_MissingPath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/secret/data/missing", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	c := testClientWithServer(t, mux)
	_, err := c.ReadSecrets([]string{"secret/missing"})
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}

func TestKvv2Path(t *testing.T) {
	cases := []struct{ in, want string }{
		{"secret/myapp/cfg", "secret/data/myapp/cfg"},
		{"secret/data/myapp/cfg", "secret/data/myapp/cfg"},
		{"secret", "secret"},
	}
	for _, tc := range cases {
		if got := kvv2Path(tc.in); got != tc.want {
			t.Errorf("kvv2Path(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
