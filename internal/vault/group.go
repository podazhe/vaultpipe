package vault

import (
	"errors"
	"fmt"
	"sort"
)

// GroupSecrets organises a flat secret map into named groups based on key prefixes.
// Each group collects secrets whose keys begin with the group's prefix, stripping
// the prefix from the key in the resulting sub-map.
//
// Example:
//
//	GroupSecrets({"DB_HOST": "localhost", "DB_PORT": "5432", "APP_NAME": "vaultpipe"}, ...)
//	// group "DB"  → {"HOST": "localhost", "PORT": "5432"}
//	// group "APP" → {"NAME": "vaultpipe"}
type GroupOptions struct {
	// Groups maps a group name to the key prefix used to select secrets.
	// The prefix is case-sensitive and must not be empty.
	Groups map[string]string

	// StripPrefix controls whether the prefix (and the separator) is removed
	// from each key inside the group. Defaults to true.
	StripPrefix bool

	// Separator is placed between prefix and the rest of the key.
	// Defaults to "_" when empty.
	Separator string

	// Ungrouped collects keys that did not match any group when set to true.
	// The resulting group is stored under the key "" (empty string).
	Ungrouped bool
}

// GroupResult holds the grouped secret maps keyed by group name.
type GroupResult map[string]map[string]string

// GroupSecrets partitions secrets into named groups defined by prefix rules.
// It returns an error if Groups is nil/empty or if any prefix is blank.
func GroupSecrets(secrets map[string]string, opts GroupOptions) (GroupResult, error) {
	if len(opts.Groups) == 0 {
		return nil, errors.New("group: at least one group must be defined")
	}

	sep := opts.Separator
	if sep == "" {
		sep = "_"
	}

	// Validate prefixes up front.
	for name, prefix := range opts.Groups {
		if prefix == "" {
			return nil, fmt.Errorf("group: prefix for group %q must not be empty", name)
		}
	}

	result := make(GroupResult, len(opts.Groups))
	for name := range opts.Groups {
		result[name] = make(map[string]string)
	}

	var ungrouped map[string]string
	if opts.Ungrouped {
		ungrouped = make(map[string]string)
	}

	for k, v := range secrets {
		matched := false
		for name, prefix := range opts.Groups {
			candidate := prefix + sep
			if len(k) >= len(candidate) && k[:len(candidate)] == candidate {
				var outKey string
				if opts.StripPrefix {
					outKey = k[len(candidate):]
				} else {
					outKey = k
				}
				if outKey != "" {
					result[name][outKey] = v
				}
				matched = true
				break
			}
		}
		if !matched && opts.Ungrouped {
			ungrouped[k] = v
		}
	}

	if opts.Ungrouped {
		result[""] = ungrouped
	}

	return result, nil
}

// GroupNames returns the sorted group names present in a GroupResult.
func GroupNames(r GroupResult) []string {
	names := make([]string, 0, len(r))
	for name := range r {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
