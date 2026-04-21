package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopySecrets_NoConflicts(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}

	res, err := CopySecrets(dst, src, CopyOptions{})
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"A", "B"}, res.Copied)
	assert.Empty(t, res.Skipped)
	assert.Equal(t, "1", dst["A"])
	assert.Equal(t, "2", dst["B"])
}

func TestCopySecrets_SkipsExisting(t *testing.T) {
	src := map[string]string{"A": "new", "B": "2"}
	dst := map[string]string{"A": "old"}

	res, err := CopySecrets(dst, src, CopyOptions{Overwrite: false})
	require.NoError(t, err)
	assert.Contains(t, res.Skipped, "A")
	assert.Contains(t, res.Copied, "B")
	assert.Equal(t, "old", dst["A"])
}

func TestCopySecrets_OverwriteExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	res, err := CopySecrets(dst, src, CopyOptions{Overwrite: true})
	require.NoError(t, err)
	assert.Contains(t, res.Copied, "A")
	assert.Empty(t, res.Skipped)
	assert.Equal(t, "new", dst["A"])
}

func TestCopySecrets_KeyFilter(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}

	res, err := CopySecrets(dst, src, CopyOptions{Keys: []string{"A", "C"}})
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"A", "C"}, res.Copied)
	_, hasB := dst["B"]
	assert.False(t, hasB)
}

func TestCopySecrets_DryRun(t *testing.T) {
	src := map[string]string{"X": "secret"}
	dst := map[string]string{}

	res, err := CopySecrets(dst, src, CopyOptions{DryRun: true})
	require.NoError(t, err)
	assert.Contains(t, res.Copied, "X")
	assert.Empty(t, dst, "dst must not be mutated during dry-run")
}

func TestCopySecrets_NilSrc(t *testing.T) {
	_, err := CopySecrets(map[string]string{}, nil, CopyOptions{})
	assert.Error(t, err)
}

func TestCopySecrets_Summary(t *testing.T) {
	r := CopyResult{Copied: []string{"A", "B"}, Skipped: []string{"C"}}
	assert.Equal(t, "copied 2 key(s), skipped 1 key(s)", r.Summary())
}
