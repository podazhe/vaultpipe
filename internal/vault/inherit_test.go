package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInheritSecrets_ParentFillsGaps(t *testing.T) {
	parent := map[string]string{"DB_HOST": "db.prod", "DB_PORT": "5432"}
	child := map[string]string{"APP_NAME": "myapp"}

	res, err := InheritSecrets(parent, child, InheritOptions{})
	require.NoError(t, err)
	assert.Equal(t, "db.prod", res.Secrets["DB_HOST"])
	assert.Equal(t, "5432", res.Secrets["DB_PORT"])
	assert.Equal(t, "myapp", res.Secrets["APP_NAME"])
	assert.ElementsMatch(t, []string{"DB_HOST", "DB_PORT"}, res.Inherited)
	assert.Empty(t, res.Overridden)
}

func TestInheritSecrets_ParentWinsOnConflictByDefault(t *testing.T) {
	parent := map[string]string{"LOG_LEVEL": "warn"}
	child := map[string]string{"LOG_LEVEL": "debug"}

	res, err := InheritSecrets(parent, child, InheritOptions{})
	require.NoError(t, err)
	assert.Equal(t, "warn", res.Secrets["LOG_LEVEL"])
	assert.Contains(t, res.Overridden, "LOG_LEVEL")
	assert.Empty(t, res.Inherited)
}

func TestInheritSecrets_ChildWinsWhenOverrideSet(t *testing.T) {
	parent := map[string]string{"LOG_LEVEL": "warn"}
	child := map[string]string{"LOG_LEVEL": "debug"}

	res, err := InheritSecrets(parent, child, InheritOptions{OverrideWithChild: true})
	require.NoError(t, err)
	assert.Equal(t, "debug", res.Secrets["LOG_LEVEL"])
	assert.Contains(t, res.Overridden, "LOG_LEVEL")
}

func TestInheritSecrets_KeyFilterLimitsInheritance(t *testing.T) {
	parent := map[string]string{"DB_HOST": "db.prod", "SECRET_KEY": "s3cr3t"}
	child := map[string]string{}

	res, err := InheritSecrets(parent, child, InheritOptions{Keys: []string{"DB_HOST"}})
	require.NoError(t, err)
	assert.Equal(t, "db.prod", res.Secrets["DB_HOST"])
	assert.NotContains(t, res.Secrets, "SECRET_KEY")
	assert.Equal(t, []string{"DB_HOST"}, res.Inherited)
}

func TestInheritSecrets_NilParentErrors(t *testing.T) {
	_, err := InheritSecrets(nil, map[string]string{}, InheritOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parent")
}

func TestInheritSecrets_NilChildErrors(t *testing.T) {
	_, err := InheritSecrets(map[string]string{}, nil, InheritOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "child")
}

func TestInheritSecrets_EmptyParentNoChange(t *testing.T) {
	child := map[string]string{"APP": "v1"}
	res, err := InheritSecrets(map[string]string{}, child, InheritOptions{})
	require.NoError(t, err)
	assert.Equal(t, child, res.Secrets)
	assert.Empty(t, res.Inherited)
	assert.Empty(t, res.Overridden)
}
