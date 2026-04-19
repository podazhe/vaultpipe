package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleSecrets = map[string]string{
	"APP_HOST":   "localhost",
	"APP_PORT":   "8080",
	"DB_HOST":    "db.local",
	"DB_PASS":    "secret",
	"UNRELATED":  "value",
}

func TestFilterSecrets_ByPrefix(t *testing.T) {
	out := FilterSecrets(sampleSecrets, FilterOptions{Prefixes: []string{"APP_"}})
	assert.Len(t, out, 2)
	assert.Contains(t, out, "APP_HOST")
	assert.Contains(t, out, "APP_PORT")
}

func TestFilterSecrets_ByKeys(t *testing.T) {
	out := FilterSecrets(sampleSecrets, FilterOptions{Keys: []string{"DB_HOST", "UNRELATED"}})
	assert.Len(t, out, 2)
	assert.Contains(t, out, "DB_HOST")
	assert.Contains(t, out, "UNRELATED")
}

func TestFilterSecrets_Exclude(t *testing.T) {
	out := FilterSecrets(sampleSecrets, FilterOptions{Exclude: []string{"DB_PASS", "UNRELATED"}})
	assert.NotContains(t, out, "DB_PASS")
	assert.NotContains(t, out, "UNRELATED")
	assert.Len(t, out, 3)
}

func TestFilterSecrets_PrefixAndExclude(t *testing.T) {
	out := FilterSecrets(sampleSecrets, FilterOptions{
		Prefixes: []string{"DB_"},
		Exclude:  []string{"DB_PASS"},
	})
	assert.Len(t, out, 1)
	assert.Contains(t, out, "DB_HOST")
}

func TestFilterSecrets_NoOptions(t *testing.T) {
	out := FilterSecrets(sampleSecrets, FilterOptions{})
	assert.Len(t, out, len(sampleSecrets))
}

func TestFilterSecrets_EmptyInput(t *testing.T) {
	out := FilterSecrets(map[string]string{}, FilterOptions{Prefixes: []string{"APP_"}})
	assert.Empty(t, out)
}
