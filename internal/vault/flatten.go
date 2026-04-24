package vault

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested maps are flattened into a single-level secret map.
type FlattenOptions struct {
	// Separator is placed between nested key segments. Defaults to "_".
	Separator string
	// Prefix is prepended to every resulting key.
	Prefix string
	// UpperCase converts all keys to upper case after flattening.
	UpperCase bool
}

// FlattenSecrets takes a nested map[string]any and returns a flat map[string]string
// suitable for use as environment variables or secret maps.
func FlattenSecrets(nested map[string]any, opts FlattenOptions) (map[string]string, error) {
	if nested == nil {
		return nil, fmt.Errorf("flatten: input map must not be nil")
	}
	sep := opts.Separator
	if sep == "" {
		sep = "_"
	}

	result := make(map[string]string)
	if err := flattenRecurse(nested, opts.Prefix, sep, result); err != nil {
		return nil, err
	}

	if opts.UpperCase {
		upper := make(map[string]string, len(result))
		for k, v := range result {
			upper[strings.ToUpper(k)] = v
		}
		return upper, nil
	}
	return result, nil
}

func flattenRecurse(node map[string]any, prefix, sep string, out map[string]string) error {
	keys := make([]string, 0, len(node))
	for k := range node {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := node[k]
		fullKey := k
		if prefix != "" {
			fullKey = prefix + sep + k
		}
		switch val := v.(type) {
		case map[string]any:
			if err := flattenRecurse(val, fullKey, sep, out); err != nil {
				return err
			}
		case string:
			out[fullKey] = val
		case nil:
			out[fullKey] = ""
		default:
			out[fullKey] = fmt.Sprintf("%v", val)
		}
	}
	return nil
}
