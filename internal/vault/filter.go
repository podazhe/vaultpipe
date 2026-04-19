package vault

import "strings"

// FilterOptions controls which secrets are included.
type FilterOptions struct {
	Prefixes []string
	Keys     []string
	Exclude  []string
}

// FilterSecrets returns a subset of secrets based on FilterOptions.
func FilterSecrets(secrets map[string]string, opts FilterOptions) map[string]string {
	out := make(map[string]string)

	for k, v := range secrets {
		if isExcluded(k, opts.Exclude) {
			continue
		}
		if len(opts.Keys) > 0 && !containsStr(opts.Keys, k) {
			continue
		}
		if len(opts.Prefixes) > 0 && !hasAnyPrefix(k, opts.Prefixes) {
			continue
		}
		out[k] = v
	}
	return out
}

func isExcluded(key string, exclude []string) bool {
	for _, e := range exclude {
		if key == e {
			return true
		}
	}
	return false
}

func hasAnyPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
