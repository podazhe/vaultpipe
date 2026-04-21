package vault

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two secret maps.
type CompareResult struct {
	OnlyInLeft  []string
	OnlyInRight []string
	Different   []string
	Identical   []string
}

// Summary returns a human-readable summary of the comparison.
func (r *CompareResult) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Only in left:  %d\n", len(r.OnlyInLeft))
	fmt.Fprintf(&sb, "Only in right: %d\n", len(r.OnlyInRight))
	fmt.Fprintf(&sb, "Different:     %d\n", len(r.Different))
	fmt.Fprintf(&sb, "Identical:     %d\n", len(r.Identical))
	return sb.String()
}

// HasDifferences returns true if any keys differ between the two maps.
func (r *CompareResult) HasDifferences() bool {
	return len(r.OnlyInLeft) > 0 || len(r.OnlyInRight) > 0 || len(r.Different) > 0
}

// CompareSecrets compares two secret maps and returns a CompareResult.
func CompareSecrets(left, right map[string]string) *CompareResult {
	result := &CompareResult{}

	allKeys := make(map[string]struct{})
	for k := range left {
		allKeys[k] = struct{}{}
	}
	for k := range right {
		allKeys[k] = struct{}{}
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		lv, inLeft := left[k]
		rv, inRight := right[k]
		switch {
		case inLeft && !inRight:
			result.OnlyInLeft = append(result.OnlyInLeft, k)
		case !inLeft && inRight:
			result.OnlyInRight = append(result.OnlyInRight, k)
		case lv != rv:
			result.Different = append(result.Different, k)
		default:
			result.Identical = append(result.Identical, k)
		}
	}

	return result
}
