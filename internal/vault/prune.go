package vault

import (
	"fmt"
	"strings"
)

// PruneOptions controls which secrets are removed.
type PruneOptions struct {
	// RemoveEmpty removes secrets whose value is an empty string.
	RemoveEmpty bool
	// RemovePrefixes removes secrets whose key starts with any of these prefixes.
	RemovePrefixes []string
	// RemoveKeys removes secrets whose key exactly matches one of these values.
	RemoveKeys []string
	// DryRun reports what would be pruned without modifying the map.
	DryRun bool
}

// PruneResult summarises the outcome of a prune operation.
type PruneResult struct {
	Removed []string
	Retained int
}

func (r PruneResult) Summary() string {
	return fmt.Sprintf("pruned %d key(s), retained %d", len(r.Removed), r.Retained)
}

// PruneSecrets removes unwanted entries from secrets according to opts.
// It returns a new map and a PruneResult describing what was removed.
func PruneSecrets(secrets map[string]string, opts PruneOptions) (map[string]string, PruneResult, error) {
	if secrets == nil {
		return nil, PruneResult{}, fmt.Errorf("prune: secrets map must not be nil")
	}

	removeSet := make(map[string]struct{}, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		removeSet[k] = struct{}{}
	}

	out := make(map[string]string, len(secrets))
	var removed []string

	for k, v := range secrets {
		if shouldPrune(k, v, removeSet, opts) {
			removed = append(removed, k)
			continue
		}
		out[k] = v
	}

	result := PruneResult{
		Removed:  removed,
		Retained: len(out),
	}

	if opts.DryRun {
		return secrets, result, nil
	}
	return out, result, nil
}

func shouldPrune(key, value string, removeSet map[string]struct{}, opts PruneOptions) bool {
	if opts.RemoveEmpty && value == "" {
		return true
	}
	if _, ok := removeSet[key]; ok {
		return true
	}
	for _, p := range opts.RemovePrefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
