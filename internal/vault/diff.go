package vault

import "fmt"

// SecretDiff represents the changes between two secret snapshots.
type SecretDiff struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]string
}

// HasChanges returns true if there are any differences.
func (d *SecretDiff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Summary returns a human-readable summary of the diff.
func (d *SecretDiff) Summary() string {
	return fmt.Sprintf("added=%d removed=%d changed=%d",
		len(d.Added), len(d.Removed), len(d.Changed))
}

// DiffSecrets computes the difference between an old and new secret map.
func DiffSecrets(old, next map[string]string) SecretDiff {
	diff := SecretDiff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]string),
	}

	for k, v := range next {
		oldVal, exists := old[k]
		if !exists {
			diff.Added[k] = v
		} else if oldVal != v {
			diff.Changed[k] = v
		}
	}

	for k, v := range old {
		if _, exists := next[k]; !exists {
			diff.Removed[k] = v
		}
	}

	return diff
}
