package vault

import (
	"fmt"
	"regexp"
	"strings"
)

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationResult holds the outcome of a secret map validation.
type ValidationResult struct {
	Errors   []string
	Warnings []string
}

func (r *ValidationResult) OK() bool { return len(r.Errors) == 0 }

func (r *ValidationResult) Summary() string {
	var sb strings.Builder
	for _, e := range r.Errors {
		sb.WriteString("ERROR: " + e + "\n")
	}
	for _, w := range r.Warnings {
		sb.WriteString("WARN:  " + w + "\n")
	}
	return sb.String()
}

// ValidateSecrets checks keys and values in a secret map.
func ValidateSecrets(secrets map[string]string, opts ValidateOptions) *ValidationResult {
	res := &ValidationResult{}
	for k, v := range secrets {
		if !validKeyRe.MatchString(k) {
			res.Errors = append(res.Errors, fmt.Sprintf("invalid key %q: must match [A-Za-z_][A-Za-z0-9_]*", k))
		}
		if opts.NoEmpty && v == "" {
			res.Errors = append(res.Errors, fmt.Sprintf("key %q has empty value", k))
		}
		if opts.WarnLong > 0 && len(v) > opts.WarnLong {
			res.Warnings = append(res.Warnings, fmt.Sprintf("key %q value exceeds %d characters", k, opts.WarnLong))
		}
		if opts.RequiredKeys != nil {
			_ = k // checked below
		}
	}
	for _, req := range opts.RequiredKeys {
		if _, ok := secrets[req]; !ok {
			res.Errors = append(res.Errors, fmt.Sprintf("required key %q is missing", req))
		}
	}
	return res
}

// ValidateOptions controls validation behaviour.
type ValidateOptions struct {
	NoEmpty      bool
	WarnLong     int
	RequiredKeys []string
}
