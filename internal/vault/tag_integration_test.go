package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTagTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/secret/data/myapp":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"data": {
					"data": {
						"DB_PASSWORD": "hunter2",
						"API_KEY": "key-abc",
						"APP_SECRET": "topsecret"
					}
				}
			}`))
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestTag_Integration_FilterSensitive(t *testing.T) {
	srv := newTagTestServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	secrets, err := ReadSecrets(client, "secret/myapp")
	if err != nil {
		t.Fatalf("ReadSecrets: %v", err)
	}

	ts := NewTaggedSecrets(secrets)
	_ = ts.Tag("DB_PASSWORD", "sensitive", "database")
	_ = ts.Tag("API_KEY", "sensitive")
	_ = ts.Tag("APP_SECRET", "sensitive", "internal")

	sensitive := ts.FilterByTag("sensitive")
	if len(sensitive) != 3 {
		t.Errorf("expected 3 sensitive secrets, got %d", len(sensitive))
	}

	dbOnly := ts.FilterByTag("database")
	if len(dbOnly) != 1 {
		t.Errorf("expected 1 database secret, got %d", len(dbOnly))
	}
	if dbOnly["DB_PASSWORD"] != "hunter2" {
		t.Errorf("unexpected value for DB_PASSWORD: %q", dbOnly["DB_PASSWORD"])
	}

	tags := ts.ListTags()
	if len(tags) != 3 {
		t.Errorf("expected 3 distinct tags, got %d: %v", len(tags), tags)
	}
}
