package vault

import (
	"fmt"
	"strings"
)

// RenameOptions controls how secrets are renamed.
type RenameOptions struct {
	// Rules maps old key names to new key names.
	Rules map[string]string
	// Prefix replaces a leading prefix on matching keys.
	OldPrefix string
	NewPrefix string
	// FailMissing returns an error if a rename rule targets a key not present.
	FailMissing bool
	// DryRun collects what would change without mutating.
	DryRun bool
}

// RenameResult describes the outcome of a rename operation.
type RenameResult struct {
	Renamed  []string // keys that were renamed (new name)
	Skipped  []string // rule targets that were absent
	Secrets  map[string]string
}

// RenameSecrets applies rename rules to secrets, returning a new map and a
// result summary. Original map is never mutated.
func RenameSecrets(secrets map[string]string, opts RenameOptions) (RenameResult, error) {
	if secrets == nil {
		return RenameResult{}, fmt.Errorf("rename: secrets map must not be nil")
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	result := RenameResult{Secrets: out}

	// Apply explicit key→key rules.
	for oldKey, newKey := range opts.Rules {
		if newKey == "" {
			return RenameResult{}, fmt.Errorf("rename: target key for %q must not be empty", oldKey)
		}
		val, ok := out[oldKey]
		if !ok {
			if opts.FailMissing {
				return RenameResult{}, fmt.Errorf("rename: key %q not found in secrets", oldKey)
			}
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}
		if !opts.DryRun {
			delete(out, oldKey)
			out[newKey] = val
		}
		result.Renamed = append(result.Renamed, newKey)
	}

	// Apply prefix substitution.
	if opts.OldPrefix != "" {
		for k, v := range out {
			if strings.HasPrefix(k, opts.OldPrefix) {
				newKey := opts.NewPrefix + strings.TrimPrefix(k, opts.OldPrefix)
				if !opts.DryRun {
					delete(out, k)
					out[newKey] = v
				}
				result.Renamed = append(result.Renamed, newKey)
			}
		}
	}

	return result, nil
}
