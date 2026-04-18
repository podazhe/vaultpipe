package vault

import (
	"fmt"
	"strings"
)

// TransformRule defines a key renaming or prefixing rule.
type TransformRule struct {
	Prefix  string
	Renames map[string]string
	Filter  []string // if non-empty, only these keys are kept
}

// Transformer applies transformation rules to a secret map.
type Transformer struct {
	rule TransformRule
}

// NewTransformer creates a Transformer with the given rule.
func NewTransformer(rule TransformRule) *Transformer {
	return &Transformer{rule: rule}
}

// Apply returns a new map with rules applied.
func (t *Transformer) Apply(secrets map[string]string) (map[string]string, error) {
	filterSet := make(map[string]bool, len(t.rule.Filter))
	for _, k := range t.rule.Filter {
		filterSet[k] = true
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if len(filterSet) > 0 && !filterSet[k] {
			continue
		}
		newKey := k
		if renamed, ok := t.rule.Renames[k]; ok {
			if renamed == "" {
				return nil, fmt.Errorf("transform: rename target for %q is empty", k)
			}
			newKey = renamed
		}
		if t.rule.Prefix != "" {
			newKey = strings.ToUpper(t.rule.Prefix) + "_" + newKey
		}
		out[newKey] = v
	}
	return out, nil
}
