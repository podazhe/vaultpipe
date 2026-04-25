package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func newCheckpointTestServer(t *testing.T, path string, data map[string]interface{}) *api.Client {
	t.Helper()
	client, closer := testClientWithServer(t, path, data)
	t.Cleanup(closer)
	return client
}

func TestCheckpoint_Integration_SaveAndRestore(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"DB_URL": "postgres://localhost/prod",
			"SECRET": "topsecret",
		},
	}
	client := newCheckpointTestServer(t, "/v1/secret/data/app", payload)
	_ = client // used by ReadSecrets in a real integration scenario

	dir := t.TempDir()
	m, err := NewCheckpointManager(dir)
	if err != nil {
		t.Fatalf("NewCheckpointManager: %v", err)
	}

	secrets := map[string]string{"DB_URL": "postgres://localhost/prod", "SECRET": "topsecret"}
	if err := m.Save("integration-v1", secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}

	cp, err := m.Load("integration-v1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cp.Secrets["DB_URL"] != secrets["DB_URL"] {
		t.Errorf("DB_URL mismatch: got %q", cp.Secrets["DB_URL"])
	}
	if cp.Secrets["SECRET"] != secrets["SECRET"] {
		t.Errorf("SECRET mismatch: got %q", cp.Secrets["SECRET"])
	}
}
