package vault

import (
	"fmt"
	"strings"
)

// ResolveOptions controls how secret paths are resolved.
type ResolveOptions struct {
	BasePath  string
	Overrides map[string]string // key -> literal value overrides
}

// ResolvedSecret holds a key and its resolved value.
type ResolvedSecret struct {
	Key    string
	Value  string
	Source string // "vault" or "override"
}

// ResolveSecrets merges vault secrets with any literal overrides.
// Overrides take precedence over vault values.
func ResolveSecrets(secrets map[string]string, opts ResolveOptions) ([]ResolvedSecret, error) {
	if secrets == nil && len(opts.Overrides) == 0 {
		return nil, fmt.Errorf("resolve: no secrets or overrides provided")
	}

	seen := make(map[string]struct{})
	var resolved []ResolvedSecret

	for k, v := range secrets {
		key := normaliseKey(k, opts.BasePath)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		source := "vault"
		if ov, ok := opts.Overrides[key]; ok {
			v = ov
			source = "override"
		}
		resolved = append(resolved, ResolvedSecret{Key: key, Value: v, Source: source})
	}

	// Add overrides that were not present in vault secrets.
	for k, v := range opts.Overrides {
		if _, exists := seen[k]; !exists {
			resolved = append(resolved, ResolvedSecret{Key: k, Value: v, Source: "override"})
		}
	}

	return resolved, nil
}

// normaliseKey strips a base-path prefix and upper-cases the key.
func normaliseKey(key, basePath string) string {
	if basePath != "" {
		key = strings.TrimPrefix(key, strings.TrimSuffix(basePath, "/")+"/")
	}
	return strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
}
