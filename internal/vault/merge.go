package vault

// MergeStrategy defines how conflicts are resolved when merging secret maps.
type MergeStrategy int

const (
	// MergeStrategyFirst keeps the value from the first (higher-priority) map.
	MergeStrategyFirst MergeStrategy = iota
	// MergeStrategyLast keeps the value from the last (lower-priority) map.
	MergeStrategyLast
)

// MergeResult holds the merged secrets and metadata about the operation.
type MergeResult struct {
	Secrets    map[string]string
	Conflicts  []string // keys that appeared in more than one source
	SourceCount int
}

// MergeSecrets merges multiple secret maps according to the given strategy.
// Sources are provided in priority order: index 0 is highest priority.
func MergeSecrets(strategy MergeStrategy, sources ...map[string]string) MergeResult {
	result := MergeResult{
		Secrets:     make(map[string]string),
		SourceCount: len(sources),
	}

	seen := make(map[string]int) // key -> first source index that set it

	for i, src := range sources {
		for k, v := range src {
			if firstIdx, exists := seen[k]; exists {
				result.Conflicts = appendUnique(result.Conflicts, k)
				if strategy == MergeStrategyLast {
					result.Secrets[k] = v
				}
				_ = firstIdx
			} else {
				seen[k] = i
				result.Secrets[k] = v
			}
		}
	}

	return result
}

func appendUnique(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}
