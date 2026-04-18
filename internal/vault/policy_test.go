package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newPolicyTestServer(t *testing.T, path string, caps []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/capabilities-self" {
			http.NotFound(w, r)
			return
		}
		resp := map[string]interface{}{
			path:           caps,
			"capabilities": caps,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": resp, "request_id": "test"})
	}))
}

func newPolicyChecker(t *testing.T, serverURL string) *PolicyChecker {
	t.Helper()
	cfg := vaultapi.DefaultConfig()
	cfg.Address = serverURL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken("test-token")
	return NewPolicyChecker(client)
}

func TestPolicyChecker_HasReadCapability(t *testing.T) {
	path := "secret/data/myapp"
	srv := newPolicyTestServer(t, path, []string{"read", "list"})
	defer srv.Close()

	pc := newPolicyChecker(t, srv.URL)
	ok, err := pc.HasCapability(context.Background(), path, CapRead)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected read capability to be present")
	}
}

func TestPolicyChecker_DenyOverrides(t *testing.T) {
	path := "secret/data/forbidden"
	srv := newPolicyTestServer(t, path, []string{"deny"})
	defer srv.Close()

	pc := newPolicyChecker(t, srv.URL)
	ok, err := pc.HasCapability(context.Background(), path, CapRead)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected deny to block read capability")
	}
}

func TestPolicyChecker_CheckPath_ReturnsList(t *testing.T) {
	path := "secret/data/multi"
	srv := newPolicyTestServer(t, path, []string{"read", "list", "create"})
	defer srv.Close()

	pc := newPolicyChecker(t, srv.URL)
	caps, err := pc.CheckPath(context.Background(), path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(caps) != 3 {
		t.Errorf("expected 3 capabilities, got %d", len(caps))
	}
}
