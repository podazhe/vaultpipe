package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDedupeSecrets_NoConflicts(t *testing.T) {
	src := []map[string]string{
		{"FOO": "bar", "BAZ": "qux"},
	}
	res, err := DedupeSecrets(src, DedupeOptions{})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"FOO": "bar", "BAZ": "qux"}, res.Secrets)
	assert.Empty(t, res.Duplicates)
}

func TestDedupeSecrets_PreferFirst(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	res, err := DedupeSecrets(sources, DedupeOptions{CaseSensitive: true})
	require.NoError(t, err)
	assert.Equal(t, "first", res.Secrets["KEY"])
}

func TestDedupeSecrets_PreferLast(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	res, err := DedupeSecrets(sources, DedupeOptions{CaseSensitive: true, PreferLast: true})
	require.NoError(t, err)
	assert.Equal(t, "second", res.Secrets["KEY"])
}

func TestDedupeSecrets_CaseInsensitive(t *testing.T) {
	sources := []map[string]string{
		{"foo": "lower"},
		{"FOO": "upper"},
	}
	res, err := DedupeSecrets(sources, DedupeOptions{CaseSensitive: false, ReportDuplicates: true})
	require.NoError(t, err)
	assert.Len(t, res.Secrets, 1)
	assert.Contains(t, res.Duplicates, "FOO")
}

func TestDedupeSecrets_ReportsDuplicates(t *testing.T) {
	sources := []map[string]string{
		{"A": "1", "B": "2"},
		{"A": "3", "C": "4"},
	}
	res, err := DedupeSecrets(sources, DedupeOptions{CaseSensitive: true, ReportDuplicates: true})
	require.NoError(t, err)
	assert.Contains(t, res.Duplicates, "A")
	assert.Equal(t, "4", res.Secrets["C"])
}

func TestDedupeSecrets_EmptySources(t *testing.T) {
	res, err := DedupeSecrets([]map[string]string{}, DedupeOptions{})
	require.NoError(t, err)
	assert.Empty(t, res.Secrets)
}
