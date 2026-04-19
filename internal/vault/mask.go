package vault

import (
	"strings"
)

// MaskOptions controls how secret values are masked.
type MaskOptions struct {
	// ShowChars is the number of trailing characters to reveal.
	ShowChars int
	// MaskChar is the character used for masking.
	MaskChar string
}

var defaultMaskOptions = MaskOptions{
	ShowChars: 4,
	MaskChar:  "*",
}

// MaskSecrets returns a copy of secrets with values partially masked.
// Keys matching sensitiveKeys are fully masked regardless of ShowChars.
func MaskSecrets(secrets map[string]string, sensitiveKeys []string, opts *MaskOptions) map[string]string {
	if opts == nil {
		opts = &defaultMaskOptions
	}
	if opts.MaskChar == "" {
		opts.MaskChar = "*"
	}

	sensitiveSet := make(map[string]struct{}, len(sensitiveKeys))
	for _, k := range sensitiveKeys {
		sensitiveSet[strings.ToUpper(k)] = struct{}{}
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		_, isSensitive := sensitiveSet[strings.ToUpper(k)]
		if isSensitive || len(v) <= opts.ShowChars {
			result[k] = strings.Repeat(opts.MaskChar, 8)
		} else {
			visible := v[len(v)-opts.ShowChars:]
			result[k] = strings.Repeat(opts.MaskChar, 8) + visible
		}
	}
	return result
}
