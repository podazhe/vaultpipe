package vault

import (
	"fmt"
	"strings"
)

// TruncateOptions controls how secret values are truncated.
type TruncateOptions struct {
	// MaxLen is the maximum allowed value length. Values longer than this are truncated.
	MaxLen int
	// Suffix is appended to truncated values to indicate truncation (e.g. "...").
	Suffix string
	// Keys restricts truncation to specific keys. If empty, all keys are considered.
	Keys []string
	// DryRun reports which keys would be truncated without modifying values.
	DryRun bool
}

// TruncateResult holds the outcome of a truncation operation.
type TruncateResult struct {
	// Truncated contains the keys whose values were (or would be) truncated.
	Truncated []string
	// Summary is a human-readable description of the operation.
	Summary string
}

// TruncateSecrets truncates secret values that exceed MaxLen.
// Returns the modified map and a result describing what changed.
func TruncateSecrets(secrets map[string]string, opts TruncateOptions) (map[string]string, TruncateResult, error) {
	if opts.MaxLen <= 0 {
		return nil, TruncateResult{}, fmt.Errorf("truncate: MaxLen must be greater than zero")
	}

	suffix := opts.Suffix
	if len(suffix) >= opts.MaxLen {
		return nil, TruncateResult{}, fmt.Errorf("truncate: suffix length %d must be less than MaxLen %d", len(suffix), opts.MaxLen)
	}

	targetSet := buildTruncateTargetSet(opts.Keys)

	out := make(map[string]string, len(secrets))
	var truncated []string

	for k, v := range secrets {
		if len(targetSet) > 0 && !targetSet[k] {
			out[k] = v
			continue
		}
		if len(v) > opts.MaxLen {
			truncated = append(truncated, k)
			if !opts.DryRun {
				out[k] = v[:opts.MaxLen-len(suffix)] + suffix
			} else {
				out[k] = v
			}
		} else {
			out[k] = v
		}
	}

	action := "truncated"
	if opts.DryRun {
		action = "would truncate"
	}

	result := TruncateResult{
		Truncated: truncated,
		Summary:   fmt.Sprintf("%s %d key(s): %s", action, len(truncated), strings.Join(truncated, ", ")),
	}

	return out, result, nil
}

func buildTruncateTargetSet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
