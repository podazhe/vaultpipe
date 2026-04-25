package vault

import (
	"errors"
	"fmt"
)

// PatchOp represents a single patch operation on a secret map.
type PatchOp struct {
	Op    string // "set", "delete", "rename"
	Key   string
	Value string // used by "set" and "rename" (new key name for rename)
}

// PatchOptions controls behaviour of PatchSecrets.
type PatchOptions struct {
	DryRun bool
}

// PatchResult holds the outcome of a patch operation.
type PatchResult struct {
	Applied []string
	Skipped []string
	Final   map[string]string
}

// PatchSecrets applies a list of patch operations to a copy of src.
// Supported ops: "set", "delete", "rename".
func PatchSecrets(src map[string]string, ops []PatchOp, opts PatchOptions) (PatchResult, error) {
	if src == nil {
		return PatchResult{}, errors.New("patch: source map must not be nil")
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	var result PatchResult

	for _, op := range ops {
		if op.Key == "" {
			return PatchResult{}, errors.New("patch: op key must not be empty")
		}
		switch op.Op {
		case "set":
			if !opts.DryRun {
				out[op.Key] = op.Value
			}
			result.Applied = append(result.Applied, fmt.Sprintf("set:%s", op.Key))
		case "delete":
			if _, exists := out[op.Key]; !exists {
				result.Skipped = append(result.Skipped, fmt.Sprintf("delete:%s", op.Key))
				continue
			}
			if !opts.DryRun {
				delete(out, op.Key)
			}
			result.Applied = append(result.Applied, fmt.Sprintf("delete:%s", op.Key))
		case "rename":
			if op.Value == "" {
				return PatchResult{}, fmt.Errorf("patch: rename op for key %q requires a non-empty target", op.Key)
			}
			val, exists := out[op.Key]
			if !exists {
				result.Skipped = append(result.Skipped, fmt.Sprintf("rename:%s", op.Key))
				continue
			}
			if !opts.DryRun {
				out[op.Value] = val
				delete(out, op.Key)
			}
			result.Applied = append(result.Applied, fmt.Sprintf("rename:%s->%s", op.Key, op.Value))
		default:
			return PatchResult{}, fmt.Errorf("patch: unknown op %q", op.Op)
		}
	}

	result.Final = out
	return result, nil
}
