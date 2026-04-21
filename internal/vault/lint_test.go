package vault

import (
	"testing"
)

func TestLintSecrets_LowercaseKey(t *testing.T) {
	secrets := map[string]string{
		"my_key": "value",
	}
	results := LintSecrets(secrets, LintOptions{AllowEmpty: true})
	if len(results) == 0 {
		t.Fatal("expected lint violation for lowercase key")
	}
	if results[0].Rule != "no-lowercase-key" {
		t.Errorf("expected rule no-lowercase-key, got %s", results[0].Rule)
	}
}

func TestLintSecrets_EmptyValue(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY": "",
	}
	results := LintSecrets(secrets, LintOptions{AllowEmpty: false})
	found := false
	for _, r := range results {
		if r.Rule == "no-empty-value" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-empty-value violation")
	}
}

func TestLintSecrets_AllowEmpty(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY": "",
	}
	results := LintSecrets(secrets, LintOptions{AllowEmpty: true})
	for _, r := range results {
		if r.Rule == "no-empty-value" {
			t.Error("expected no-empty-value to be suppressed")
		}
	}
}

func TestLintSecrets_MaxValueLen(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY": "this-is-a-very-long-secret-value",
	}
	results := LintSecrets(secrets, LintOptions{AllowEmpty: true, MaxValueLen: 10})
	found := false
	for _, r := range results {
		if r.Rule == "max-value-length" {
			found = true
		}
	}
	if !found {
		t.Error("expected max-value-length violation")
	}
}

func TestLintSecrets_ForbiddenPrefix(t *testing.T) {
	secrets := map[string]string{
		"TEST_KEY": "val",
	}
	results := LintSecrets(secrets, LintOptions{AllowEmpty: true, ForbidPrefix: []string{"TEST_"}})
	found := false
	for _, r := range results {
		if r.Rule == "forbidden-prefix" {
			found = true
		}
	}
	if !found {
		t.Error("expected forbidden-prefix violation")
	}
}

func TestLintSecrets_CustomRule(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY": "secret123",
	}
	customRule := LintRule{
		Name:    "no-numeric-suffix",
		Message: "value ends with a digit",
		Check: func(_, v string) bool {
			if len(v) == 0 {
				return false
			}
			last := v[len(v)-1]
			return last >= '0' && last <= '9'
		},
	}
	results := LintSecrets(secrets, LintOptions{AllowEmpty: true, CustomRules: []LintRule{customRule}})
	found := false
	for _, r := range results {
		if r.Rule == "no-numeric-suffix" {
			found = true
		}
	}
	if !found {
		t.Error("expected custom rule violation")
	}
}

func TestLintSecrets_Clean(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY":    "cleanvalue",
		"OTHER_KEY": "anothervalue",
	}
	results := LintSecrets(secrets, LintOptions{AllowEmpty: true})
	if len(results) != 0 {
		t.Errorf("expected no violations, got %d", len(results))
	}
}
