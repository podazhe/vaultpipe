package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/vaultpipe/internal/vault"
)

func newDedupeTestServer(t *testing.T, paths map[string]map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for path, data := range paths {
			if r.URL.Path == "/v1/secret/data/"+path {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"data": map[string]interface{}{"data": data},
				})
				return
			}
		}
		http.NotFound(w, r)
	}))
}

func TestDedupe_Integration_MultipleSources(t *testing.T) {
	src1 := map[string]string{"DB_HOST": "primary", "DB_PORT": "5432"}
	src2 := map[string]string{"DB_HOST": "replica", "API_KEY": "secret"}

	res, err := vault.DedupeSecrets(
		[]map[string]string{src1, src2},
		vault.DedupeOptions{CaseSensitive: true, ReportDuplicates: true},
	)
	require.NoError(t, err)
	assert.Equal(t, "primary", res.Secrets["DB_HOST"], "first source should win")
	assert.Equal(t, "5432", res.Secrets["DB_PORT"])
	assert.Equal(t, "secret", res.Secrets["API_KEY"])
	assert.Contains(t, res.Duplicates, "DB_HOST")
}
