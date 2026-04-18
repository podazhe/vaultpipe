package vault

import (
	"fmt"
	"strings"
)

// SecretMap is a flat map of key -> secret value.
type SecretMap map[string]string

// ReadSecrets reads secrets from one or more Vault paths and merges them
// into a single SecretMap. KV v2 paths are detected by the presence of
// "secret/data/" or an explicit version field in the response.
func (c *Client) ReadSecrets(paths []string) (SecretMap, error) {
	result := make(SecretMap)
	for _, p := range paths {
		p = strings.TrimPrefix(p, "/")
		secret, err := c.logical.Read(kvv2Path(p))
		if err != nil {
			return nil, fmt.Errorf("reading path %q: %w", p, err)
		}
		if secret == nil {
			return nil, fmt.Errorf("no secret found at path %q", p)
		}

		data, err := extractData(secret.Data)
		if err != nil {
			return nil, fmt.Errorf("extracting data from %q: %w", p, err)
		}

		for k, v := range data {
			result[k] = v
		}
	}
	return result, nil
}

// kvv2Path converts a logical path to a KV v2 data path when applicable.
func kvv2Path(p string) string {
	// Already a data path.
	if strings.Contains(p, "/data/") {
		return p
	}
	// Heuristic: paths starting with <mount>/  become <mount>/data/<rest>.
	parts := strings.SplitN(p, "/", 2)
	if len(parts) == 2 {
		return parts[0] + "/data/" + parts[1]
	}
	return p
}

// extractData handles both raw maps and KV v2 nested {"data": {...}} envelopes.
func extractData(raw map[string]interface{}) (map[string]string, error) {
	out := make(map[string]string)

	if nested, ok := raw["data"]; ok {
		if m, ok := nested.(map[string]interface{}); ok {
			raw = m
		}
	}

	for k, v := range raw {
		switch val := v.(type) {
		case string:
			out[k] = val
		case fmt.Stringer:
			out[k] = val.String()
		default:
			out[k] = fmt.Sprintf("%v", val)
		}
	}
	return out, nil
}
