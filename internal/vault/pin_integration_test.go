package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newPinTestServer(t *testing.T, path string, version int, data map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		envelope := map[string]interface{}{
			"data": map[string]interface{}{
				"data":     data,
				"metadata": map[string]interface{}{"version": float64(version)},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(envelope)
	}))
}

func TestPinManager_Integration_PinFromVault(t *testing.T) {
	expectedData := map[string]string{"DB_PASS": "s3cr3t"}
	srv := newPinTestServer(t, "secret/data/db", 5, expectedData)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	secrets, err := client.ReadSecrets("secret/db")
	if err != nil {
		t.Fatalf("ReadSecrets: %v", err)
	}

	pm := NewPinManager()
	if err := pm.Pin("secret/db", 5, secrets); err != nil {
		t.Fatalf("Pin: %v", err)
	}

	p, err := pm.Get("secret/db")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if p.Version != 5 {
		t.Errorf("expected version 5, got %d", p.Version)
	}
	if p.Data["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %s", p.Data["DB_PASS"])
	}
}
