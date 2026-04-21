package vault

import "fmt"

// CopyOptions controls behaviour of CopySecrets.
type CopyOptions struct {
	// Overwrite existing keys in the destination map.
	Overwrite bool
	// Keys is an optional allowlist; when non-empty only these keys are copied.
	Keys []string
	// DryRun reports what would change without mutating dst.
	DryRun bool
}

// CopyResult summarises the outcome of a CopySecrets call.
type CopyResult struct {
	Copied  []string
	Skipped []string
}

func (r CopyResult) Summary() string {
	return fmt.Sprintf("copied %d key(s), skipped %d key(s)", len(r.Copied), len(r.Skipped))
}

// CopySecrets copies secrets from src into dst according to opts.
// dst is not modified when DryRun is true.
func CopySecrets(dst, src map[string]string, opts CopyOptions) (CopyResult, error) {
	if src == nil {
		return CopyResult{}, fmt.Errorf("source map must not be nil")
	}
	if dst == nil {
		return CopyResult{}, fmt.Errorf("destination map must not be nil")
	}

	allowlist := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowlist[k] = true
	}

	var result CopyResult

	for k, v := range src {
		if len(allowlist) > 0 && !allowlist[k] {
			continue
		}
		if _, exists := dst[k]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if !opts.DryRun {
			dst[k] = v
		}
		result.Copied = append(result.Copied, k)
	}

	return result, nil
}
