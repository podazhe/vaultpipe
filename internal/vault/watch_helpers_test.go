package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/api"
)

// newWatchTestServer returns a test server that always responds with the given JSON body.
func newWatchTestServer(t *testing.T, body string) (*httptest.Server, *api.Client) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	client := testClientWithServer(t, server)
	return server, client
}

func TestMapsEqual_Empty(t *testing.T) {
	if !mapsEqual(nil, nil) {
		t.Error("nil maps should be equal")
	}
	if !mapsEqual(map[string]string{}, map[string]string{}) {
		t.Error("empty maps should be equal")
	}
}

func TestMapsEqual_ExtraKey(t *testing.T) {
	a := map[string]string{"x": "1", "y": "2"}
	b := map[string]string{"x": "1"}
	if mapsEqual(a, b) {
		t.Error("maps with different lengths should not be equal")
	}
}
