package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newIndexTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/secret/data/index-test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"data":{"HOST":"localhost","PORT":"5432","PASS":"s3cr3t"}}}`))
	})
	return httptest.NewServer(mux)
}

func TestSecretIndex_Integration_BuildFromSecrets(t *testing.T) {
	srv := newIndexTestServer(t)
	defer srv.Close()

	secrets := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"PASS": "s3cr3t",
	}

	idx := NewSecretIndex()
	for k, _ := range secrets {
		tags := []string{}
		if k == "PASS" {
			tags = []string{"sensitive"}
		}
		idx.Add("secret/index-test", k, tags, map[string]string{"source": "vault"})
	}

	if idx.Size() != 3 {
		t.Fatalf("expected 3 entries, got %d", idx.Size())
	}

	e, ok := idx.Get("secret/index-test", "PASS")
	if !ok {
		t.Fatal("PASS not indexed")
	}
	if len(e.Tags) == 0 || e.Tags[0] != "sensitive" {
		t.Errorf("expected sensitive tag on PASS")
	}

	results := idx.Search("HOST")
	if len(results) != 1 || results[0].Key != "HOST" {
		t.Errorf("search for HOST failed: %v", results)
	}
}
