package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newNormalizeTestServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		envelope := map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		}
		_ = json.NewEncoder(w).Encode(envelope)
	}))
}

func TestNormalize_Integration_UpperAndTrim(t *testing.T) {
	raw := map[string]string{
		"db-host":  "  postgres  ",
		"api-key":  "secret",
		"APP_NAME": "vaultpipe",
	}

	res := NormalizeSecrets(raw, NormalizeOptions{
		UppercaseKeys:  true,
		TrimValues:     true,
		ReplaceHyphens: true,
	})

	cases := map[string]string{
		"DB_HOST":  "postgres",
		"API_KEY":  "secret",
		"APP_NAME": "vaultpipe",
	}
	for key, want := range cases {
		if got := res.Normalized[key]; got != want {
			t.Errorf("key %s: want %q, got %q", key, want, got)
		}
	}

	// APP_NAME was already uppercase with no hyphen — should not appear in changes
	for _, ch := range res.Changes {
		if ch.Key == "APP_NAME" && ch.OldKey == "" && ch.OldValue == "" {
			t.Error("APP_NAME should not have been recorded as a change")
		}
	}
}
