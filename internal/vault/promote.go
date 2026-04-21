package vault

import (
	"fmt"
	"sort"
)

// PromoteOptions configures how secrets are promoted between environments.
type PromoteOptions struct {
	Overwrite bool
	DryRun    bool
	Keys      []string // if non-empty, only promote these keys
}

// PromoteResult describes the outcome of a promotion operation.
type PromoteResult struct {
	Promoted  []string
	Skipped   []string
	Overwrite []string
}

// Summary returns a human-readable summary of the promotion result.
func (r PromoteResult) Summary() string {
	return fmt.Sprintf(
		"promoted=%d skipped=%d overwritten=%d",
		len(r.Promoted), len(r.Skipped), len(r.Overwrite),
	)
}

// PromoteSecrets copies secrets from src into dst according to opts.
// Keys already present in dst are skipped unless Overwrite is true.
func PromoteSecrets(src, dst map[string]string, opts PromoteOptions) (map[string]string, PromoteResult) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var result PromoteResult

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		if _, exists := out[k]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if _, exists := out[k]; exists && opts.Overwrite {
			result.Overwrite = append(result.Overwrite, k)
		}
		if !opts.DryRun {
			out[k] = v
		}
		result.Promoted = append(result.Promoted, k)
	}

	return out, result
}
