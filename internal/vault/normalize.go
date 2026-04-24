package vault

import (
	"strings"
	"unicode"
)

// NormalizeOptions controls how secret keys and values are normalized.
type NormalizeOptions struct {
	UppercaseKeys   bool
	TrimValues      bool
	ReplaceHyphens  bool // replace hyphens with underscores in keys
	StripNonPrint   bool // strip non-printable characters from values
}

// NormalizeResult holds the outcome of a normalization pass.
type NormalizeResult struct {
	Normalized map[string]string
	Changes    []NormalizeChange
}

// NormalizeChange records a single key or value transformation.
type NormalizeChange struct {
	Key      string
	OldKey   string // non-empty when the key itself changed
	OldValue string // non-empty when the value changed
	NewValue string
}

// NormalizeSecrets applies the given options to a map of secrets and returns
// a new map along with a change log describing every transformation made.
func NormalizeSecrets(secrets map[string]string, opts NormalizeOptions) NormalizeResult {
	result := NormalizeResult{
		Normalized: make(map[string]string, len(secrets)),
	}

	for k, v := range secrets {
		newKey := k
		newVal := v
		change := NormalizeChange{Key: k}

		if opts.ReplaceHyphens {
			newKey = strings.ReplaceAll(newKey, "-", "_")
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}
		if opts.StripNonPrint {
			newVal = stripNonPrintable(newVal)
		}

		if newKey != k {
			change.OldKey = k
			change.Key = newKey
		}
		if newVal != v {
			change.OldValue = v
			change.NewValue = newVal
		}
		if change.OldKey != "" || change.OldValue != "" {
			result.Changes = append(result.Changes, change)
		}

		result.Normalized[newKey] = newVal
	}

	return result
}

func stripNonPrintable(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, s)
}
