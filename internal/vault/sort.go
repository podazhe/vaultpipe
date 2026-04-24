package vault

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering direction for secrets.
type SortOrder int

const (
	SortAsc  SortOrder = iota // alphabetical ascending
	SortDesc                  // alphabetical descending
)

// SortOptions controls how secrets are sorted.
type SortOptions struct {
	Order      SortOrder
	ByValue    bool   // sort by value instead of key
	IgnoreCase bool   // case-insensitive comparison
	Prefix     string // only sort entries whose key starts with Prefix
}

// SortSecrets returns a new slice of key-value pairs sorted according to opts.
// The original map is not modified.
func SortSecrets(secrets map[string]string, opts SortOptions) []KeyValue {
	pairs := make([]KeyValue, 0, len(secrets))
	for k, v := range secrets {
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		pairs = append(pairs, KeyValue{Key: k, Value: v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		var a, b string
		if opts.ByValue {
			a, b = pairs[i].Value, pairs[j].Value
		} else {
			a, b = pairs[i].Key, pairs[j].Key
		}
		if opts.IgnoreCase {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}
		if opts.Order == SortDesc {
			return a > b
		}
		return a < b
	})

	return pairs
}

// KeyValue holds a single secret key-value pair.
type KeyValue struct {
	Key   string
	Value string
}

// ToMap converts a sorted slice back into a plain map.
func ToMap(pairs []KeyValue) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		out[p.Key] = p.Value
	}
	return out
}
