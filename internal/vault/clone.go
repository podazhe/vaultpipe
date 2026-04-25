package vault

import (
	"errors"
	"fmt"
)

// CloneOptions controls the behaviour of CloneSecrets.
type CloneOptions struct {
	// Prefix is prepended to every key in the cloned map.
	Prefix string
	// KeyFilter restricts cloning to keys whose names appear in the slice.
	// An empty slice means all keys are cloned.
	KeyFilter []string
	// Overwrite allows existing keys in dst to be overwritten.
	Overwrite bool
	// DryRun reports what would happen without modifying dst.
	DryRun bool
}

// CloneResult summarises the outcome of a CloneSecrets call.
type CloneResult struct {
	Cloned    []string
	Skipped   []string
	Overwrote []string
}

// Summary returns a human-readable one-liner.
func (r CloneResult) Summary() string {
	return fmt.Sprintf("cloned=%d skipped=%d overwrote=%d",
		len(r.Cloned), len(r.Skipped), len(r.Overwrote))
}

// CloneSecrets copies secrets from src into dst according to opts.
// src and dst must not be nil.
func CloneSecrets(src, dst map[string]string, opts CloneOptions) (CloneResult, error) {
	if src == nil {
		return CloneResult{}, errors.New("clone: src must not be nil")
	}
	if dst == nil {
		return CloneResult{}, errors.New("clone: dst must not be nil")
	}

	filterSet := make(map[string]struct{}, len(opts.KeyFilter))
	for _, k := range opts.KeyFilter {
		filterSet[k] = struct{}{}
	}

	var result CloneResult

	for k, v := range src {
		if len(filterSet) > 0 {
			if _, ok := filterSet[k]; !ok {
				continue
			}
		}

		dstKey := opts.Prefix + k

		if _, exists := dst[dstKey]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, dstKey)
				continue
			}
			if !opts.DryRun {
				dst[dstKey] = v
			}
			result.Overwrote = append(result.Overwrote, dstKey)
			continue
		}

		if !opts.DryRun {
			dst[dstKey] = v
		}
		result.Cloned = append(result.Cloned, dstKey)
	}

	return result, nil
}
