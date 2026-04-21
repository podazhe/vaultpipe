package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newPromoteTestServer(t *testing.T, srcData, dstData map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/secret/data/src":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"data": srcData},
			})
		case "/v1/secret/data/dst":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"data": dstData},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestPromote_Integration_NoConflict(t *testing.T) {
	src := map[string]string{"TOKEN": "abc", "SECRET": "xyz"}
	dst := map[string]string{"OTHER": "val"}

	out, res := PromoteSecrets(src, dst, PromoteOptions{})

	if len(res.Promoted) != 2 {
		t.Fatalf("expected 2 promoted keys, got %d", len(res.Promoted))
	}
	if out["TOKEN"] != "abc" || out["SECRET"] != "xyz" || out["OTHER"] != "val" {
		t.Fatalf("unexpected merged output: %v", out)
	}
}

func TestPromote_Integration_DryRunPreservesDestination(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old", "B": "keep"}

	out, res := PromoteSecrets(src, dst, PromoteOptions{DryRun: true, Overwrite: true})

	// dry-run: destination must be unchanged
	if out["A"] != "old" {
		t.Fatalf("dry-run must not mutate destination, got A=%s", out["A"])
	}
	if len(res.Promoted) != 1 {
		t.Fatalf("expected 1 in promoted list for reporting, got %+v", res)
	}
}
