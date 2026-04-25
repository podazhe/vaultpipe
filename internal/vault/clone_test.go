package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneSecrets_NoConflicts(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}

	res, err := CloneSecrets(src, dst, CloneOptions{})
	require.NoError(t, err)
	assert.Len(t, res.Cloned, 2)
	assert.Empty(t, res.Skipped)
	assert.Equal(t, "1", dst["A"])
	assert.Equal(t, "2", dst["B"])
}

func TestCloneSecrets_SkipsExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	res, err := CloneSecrets(src, dst, CloneOptions{Overwrite: false})
	require.NoError(t, err)
	assert.Len(t, res.Skipped, 1)
	assert.Equal(t, "old", dst["A"])
}

func TestCloneSecrets_OverwriteExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	res, err := CloneSecrets(src, dst, CloneOptions{Overwrite: true})
	require.NoError(t, err)
	assert.Len(t, res.Overwrote, 1)
	assert.Equal(t, "new", dst["A"])
}

func TestCloneSecrets_WithPrefix(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	dst := map[string]string{}

	res, err := CloneSecrets(src, dst, CloneOptions{Prefix: "PROD_"})
	require.NoError(t, err)
	assert.Len(t, res.Cloned, 1)
	assert.Equal(t, "val", dst["PROD_KEY"])
	_, original := dst["KEY"]
	assert.False(t, original)
}

func TestCloneSecrets_KeyFilter(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}

	res, err := CloneSecrets(src, dst, CloneOptions{KeyFilter: []string{"A", "C"}})
	require.NoError(t, err)
	assert.Len(t, res.Cloned, 2)
	_, hasB := dst["B"]
	assert.False(t, hasB)
}

func TestCloneSecrets_DryRun(t *testing.T) {
	src := map[string]string{"X": "secret"}
	dst := map[string]string{}

	res, err := CloneSecrets(src, dst, CloneOptions{DryRun: true})
	require.NoError(t, err)
	assert.Len(t, res.Cloned, 1)
	assert.Empty(t, dst, "dry-run must not mutate dst")
}

func TestCloneSecrets_NilSrcErrors(t *testing.T) {
	_, err := CloneSecrets(nil, map[string]string{}, CloneOptions{})
	assert.Error(t, err)
}

func TestCloneSecrets_NilDstErrors(t *testing.T) {
	_, err := CloneSecrets(map[string]string{}, nil, CloneOptions{})
	assert.Error(t, err)
}

func TestCloneResult_Summary(t *testing.T) {
	r := CloneResult{
		Cloned:    []string{"A", "B"},
		Skipped:   []string{"C"},
		Overwrote: []string{},
	}
	assert.Equal(t, "cloned=2 skipped=1 overwrote=0", r.Summary())
}
