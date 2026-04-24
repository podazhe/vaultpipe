package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newChainTestServer(secrets map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/secret/data/app" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"data": secrets},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestChain_Integration_PrefixAndRedact(t *testing.T) {
	raw := map[string]string{
		"db_password": "s3cr3t",
		"api_key":     "abc123",
		"region":      "us-east-1",
	}

	chain := NewChain().
		Add("prefix", func(m map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(m))
			for k, v := range m {
				out["APP_"+k] = v
			}
			return out, nil
		}).
		Add("redact", func(m map[string]string) (map[string]string, error) {
			return RedactSecrets(m, RedactOptions{}), nil
		})

	result, err := chain.Run(raw)
	if err != nil {
		t.Fatalf("chain failed: %v", err)
	}

	if result["APP_db_password"] != "[REDACTED]" {
		t.Errorf("expected redacted password, got %q", result["APP_db_password"])
	}
	if result["APP_region"] != "us-east-1" {
		t.Errorf("expected region to pass through, got %q", result["APP_region"])
	}
}
