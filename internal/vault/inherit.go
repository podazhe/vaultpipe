package vault

import "fmt"

// InheritOptions controls how secrets are inherited from a parent path.
type InheritOptions struct {
	// OverrideWithChild means child values win on conflict.
	OverrideWithChild bool
	// Keys restricts inheritance to specific keys; empty means all keys.
	Keys []string
}

// InheritResult holds the merged output and metadata about the operation.
type InheritResult struct {
	Secrets   map[string]string
	Inherited []string // keys taken from parent
	Overridden []string // keys where child won
}

// InheritSecrets merges parent secrets into child secrets according to opts.
// The child map is never mutated; a new map is returned in the result.
func InheritSecrets(parent, child map[string]string, opts InheritOptions) (InheritResult, error) {
	if parent == nil {
		return InheritResult{}, fmt.Errorf("inherit: parent secrets must not be nil")
	}
	if child == nil {
		return InheritResult{}, fmt.Errorf("inherit: child secrets must not be nil")
	}

	allowedKeys := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		allowedKeys[k] = struct{}{}
	}

	result := make(map[string]string, len(child))
	for k, v := range child {
		result[k] = v
	}

	var inherited, overridden []string

	for k, parentVal := range parent {
		if len(allowedKeys) > 0 {
			if _, ok := allowedKeys[k]; !ok {
				continue
			}
		}

		childVal, exists := result[k]
		switch {
		case !exists:
			result[k] = parentVal
			inherited = append(inherited, k)
		case exists && opts.OverrideWithChild:
			// child already has the value and wins — record it
			_ = childVal
			overridden = append(overridden, k)
		default:
			// parent wins over child
			result[k] = parentVal
			overridden = append(overridden, k)
		}
	}

	return InheritResult{
		Secrets:    result,
		Inherited:  inherited,
		Overridden: overridden,
	}, nil
}
