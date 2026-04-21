package vault

import (
	"fmt"
	"strings"
)

// LintRule represents a single linting rule applied to secrets.
type LintRule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// LintResult holds the outcome of a lint check for a single secret key.
type LintResult struct {
	Key     string
	Rule    string
	Message string
}

// LintOptions configures which rules are applied during linting.
type LintOptions struct {
	AllowEmpty    bool
	MaxValueLen   int
	ForbidPrefix  []string
	CustomRules   []LintRule
}

var defaultLintRules = []LintRule{
	{
		Name:    "no-lowercase-key",
		Message: "key contains lowercase letters; prefer UPPER_SNAKE_CASE",
		Check: func(key, _ string) bool {
			return key != strings.ToUpper(key)
		},
	},
	{
		Name:    "no-space-in-key",
		Message: "key contains spaces",
		Check: func(key, _ string) bool {
			return strings.Contains(key, " ")
		},
	},
}

// LintSecrets runs all applicable lint rules against the provided secrets map
// and returns a slice of LintResult for every violation found.
func LintSecrets(secrets map[string]string, opts LintOptions) []LintResult {
	var results []LintResult

	rules := append([]LintRule{}, defaultLintRules...)
	rules = append(rules, opts.CustomRules...)

	for key, value := range secrets {
		if !opts.AllowEmpty && value == "" {
			results = append(results, LintResult{
				Key:     key,
				Rule:    "no-empty-value",
				Message: "value is empty",
			})
		}

		if opts.MaxValueLen > 0 && len(value) > opts.MaxValueLen {
			results = append(results, LintResult{
				Key:     key,
				Rule:    "max-value-length",
				Message: fmt.Sprintf("value exceeds max length of %d", opts.MaxValueLen),
			})
		}

		for _, prefix := range opts.ForbidPrefix {
			if strings.HasPrefix(key, prefix) {
				results = append(results, LintResult{
					Key:     key,
					Rule:    "forbidden-prefix",
					Message: fmt.Sprintf("key has forbidden prefix %q", prefix),
				})
			}
		}

		for _, rule := range rules {
			if rule.Check(key, value) {
				results = append(results, LintResult{
					Key:     key,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}

	return results
}
