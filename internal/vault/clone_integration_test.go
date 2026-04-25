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

func newCloneTestServer(t *testing.T, secrets map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{"data": map[string]interface{}{"data": secrets}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data) //nolint:errcheck
	}))
}

func TestClone_Integration_PrefixedCopy(t *testing.T) {
	src := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "s3cr3t",
		"API_KEY": "key123",
	}
	dst := map[string]string{}

	res, err := vault.CloneSecrets(src, dst, vault.CloneOptions{
		Prefix:    "STAGING_",
		KeyFilter: []string{"DB_HOST", "DB_PASS"},
	})
	require.NoError(t, err)
	assert.Len(t, res.Cloned, 2)
	assert.Equal(t, "localhost", dst["STAGING_DB_HOST"])
	assert.Equal(t, "s3cr3t", dst["STAGING_DB_PASS"])
	_, hasAPI := dst["STAGING_API_KEY"]
	assert.False(t, hasAPI, "filtered key must not appear in dst")
}

func TestClone_Integration_DryRunPreservesDst(t *testing.T) {
	src := map[string]string{"TOKEN": "abc", "SECRET": "xyz"}
	dst := map[string]string{"TOKEN": "existing"}

	res, err := vault.CloneSecrets(src, dst, vault.CloneOptions{
		Overwrite: true,
		DryRun:    true,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, res.Cloned)
	assert.Equal(t, "existing", dst["TOKEN"], "dry-run must not overwrite")
	_, hasSecret := dst["SECRET"]
	assert.False(t, hasSecret, "dry-run must not insert new keys")
}
