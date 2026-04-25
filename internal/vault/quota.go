package vault

import (
	"errors"
	"fmt"
)

// QuotaOptions configures how secret quota enforcement behaves.
type QuotaOptions struct {
	MaxKeys      int
	MaxValueLen  int
	MaxTotalSize int // bytes (sum of all value lengths)
	DryRun       bool
}

// QuotaViolation describes a single quota breach.
type QuotaViolation struct {
	Rule    string
	Detail  string
}

func (v QuotaViolation) Error() string {
	return fmt.Sprintf("quota violation [%s]: %s", v.Rule, v.Detail)
}

// QuotaResult holds the outcome of a quota check.
type QuotaResult struct {
	Violations []QuotaViolation
	TotalKeys  int
	TotalSize  int
}

func (r QuotaResult) OK() bool { return len(r.Violations) == 0 }

// EnforceQuota checks secrets against the provided QuotaOptions.
// In dry-run mode it collects all violations without returning an error.
// In normal mode it returns the first violation as an error.
func EnforceQuota(secrets map[string]string, opts QuotaOptions) (QuotaResult, error) {
	result := QuotaResult{
		TotalKeys: len(secrets),
	}

	for _, v := range secrets {
		result.TotalSize += len(v)
	}

	if opts.MaxKeys > 0 && result.TotalKeys > opts.MaxKeys {
		result.Violations = append(result.Violations, QuotaViolation{
			Rule:   "max_keys",
			Detail: fmt.Sprintf("have %d keys, limit is %d", result.TotalKeys, opts.MaxKeys),
		})
	}

	if opts.MaxTotalSize > 0 && result.TotalSize > opts.MaxTotalSize {
		result.Violations = append(result.Violations, QuotaViolation{
			Rule:   "max_total_size",
			Detail: fmt.Sprintf("total size %d bytes exceeds limit %d", result.TotalSize, opts.MaxTotalSize),
		})
	}

	if opts.MaxValueLen > 0 {
		for k, v := range secrets {
			if len(v) > opts.MaxValueLen {
				result.Violations = append(result.Violations, QuotaViolation{
					Rule:   "max_value_len",
					Detail: fmt.Sprintf("key %q value length %d exceeds limit %d", k, len(v), opts.MaxValueLen),
				})
			}
		}
	}

	if !opts.DryRun && len(result.Violations) > 0 {
		return result, errors.New(result.Violations[0].Error())
	}

	return result, nil
}
