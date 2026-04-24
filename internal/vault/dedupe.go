package vault

import "strings"

// DedupeOptions controls how deduplication is performed.
type DedupeOptions struct {
	// CaseSensitive determines whether key comparison is case-sensitive.
	CaseSensitive bool
	// PreferLast keeps the last occurrence of a duplicate key instead of the first.
	PreferLast bool
	// ReportDuplicates collects keys that were removed during deduplication.
	ReportDuplicates bool
}

// DedupeResult holds the deduplicated secrets and any metadata.
type DedupeResult struct {
	Secrets    map[string]string
	Duplicates []string
}

// DedupeSecrets removes duplicate keys from one or more secret maps,
// merging them in order. By default the first occurrence wins.
func DedupeSecrets(sources []map[string]string, opts DedupeOptions) (DedupeResult, error) {
	result := make(map[string]string)
	seen := make(map[string]string) // normalised key -> original key
	var duplicates []string

	for _, src := range sources {
		for k, v := range src {
			norm := normaliseDedupeKey(k, opts.CaseSensitive)
			if orig, exists := seen[norm]; exists {
				if opts.ReportDuplicates {
					duplicates = append(duplicates, k)
				}
				if opts.PreferLast {
					delete(result, orig)
					result[k] = v
					seen[norm] = k
				}
				continue
			}
			seen[norm] = k
			result[k] = v
		}
	}

	return DedupeResult{Secrets: result, Duplicates: duplicates}, nil
}

func normaliseDedupeKey(k string, caseSensitive bool) string {
	if caseSensitive {
		return k
	}
	return strings.ToLower(k)
}
