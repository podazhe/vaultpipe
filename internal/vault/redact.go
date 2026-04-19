package vault

import "strings"

// RedactOptions controls which keys are redacted in output.
type RedactOptions struct {
	// SensitiveKeys are exact key names to redact.
	SensitiveKeys []string
	// SensitivePrefixes are key prefixes that trigger redaction.
	SensitivePrefixes []string
}

const redactedValue = "***REDACTED***"

var defaultSensitivePrefixes = []string{
	"SECRET", "PASSWORD", "PASS", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL",
}

// RedactSecrets returns a copy of secrets with sensitive values replaced.
func RedactSecrets(secrets map[string]string, opts RedactOptions) map[string]string {
	prefixes := opts.SensitivePrefixes
	if len(prefixes) == 0 {
		prefixes = defaultSensitivePrefixes
	}

	keySet := make(map[string]struct{}, len(opts.SensitiveKeys))
	for _, k := range opts.SensitiveKeys {
		keySet[strings.ToUpper(k)] = struct{}{}
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		upper := strings.ToUpper(k)
		if _, ok := keySet[upper]; ok {
			out[k] = redactedValue
			continue
		}
		if hasAnyPrefix(upper, prefixes) {
			out[k] = redactedValue
			continue
		}
		out[k] = v
	}
	return out
}
