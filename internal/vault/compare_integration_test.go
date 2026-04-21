package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newCompareTestServer(t *testing.T, leftData, rightData map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data map[string]interface{}
		switch r.URL.Path {
		case "/v1/secret/data/left":
			data = leftData
		case "/v1/secret/data/right":
			data = rightData
		default:
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		})
	}))
}

func TestCompare_Integration_DetectsDifferences(t *testing.T) {
	leftData := map[string]interface{}{"KEY_A": "val1", "KEY_B": "shared"}
	rightData := map[string]interface{}{"KEY_B": "shared", "KEY_C": "val3"}

	srv := newCompareTestServer(t, leftData, rightData)
	defer srv.Close()

	client := testClientWithServer(t, srv)

	left, err := ReadSecrets(context.Background(), client, "secret/left")
	if err != nil {
		t.Fatalf("read left: %v", err)
	}
	right, err := ReadSecrets(context.Background(), client, "secret/right")
	if err != nil {
		t.Fatalf("read right: %v", err)
	}

	result := CompareSecrets(left, right)
	if !result.HasDifferences() {
		t.Error("expected differences between left and right")
	}
	if len(result.OnlyInLeft) != 1 || result.OnlyInLeft[0] != "KEY_A" {
		t.Errorf("expected OnlyInLeft=[KEY_A], got %v", result.OnlyInLeft)
	}
	if len(result.OnlyInRight) != 1 || result.OnlyInRight[0] != "KEY_C" {
		t.Errorf("expected OnlyInRight=[KEY_C], got %v", result.OnlyInRight)
	}
	if len(result.Identical) != 1 || result.Identical[0] != "KEY_B" {
		t.Errorf("expected Identical=[KEY_B], got %v", result.Identical)
	}
}
