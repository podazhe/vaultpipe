package vault

import (
	"fmt"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exported secrets.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatYAML   ExportFormat = "yaml"
)

// ExportOptions controls how secrets are serialised.
type ExportOptions struct {
	Format  ExportFormat
	Export  bool   // prefix lines with `export`
	Prefix  string // optional key prefix
}

// ExportSecrets serialises a secrets map into the requested format.
func ExportSecrets(secrets map[string]string, opts ExportOptions) (string, error) {
	keys := sortedKeys(secrets)
	switch opts.Format {
	case FormatDotenv, "":
		return exportDotenv(secrets, keys, opts), nil
	case FormatJSON:
		return exportJSON(secrets, keys, opts.Prefix), nil
	case FormatYAML:
		return exportYAML(secrets, keys, opts.Prefix), nil
	default:
		return "", fmt.Errorf("unsupported export format: %s", opts.Format)
	}
}

func exportDotenv(secrets map[string]string, keys []string, opts ExportOptions) string {
	var sb strings.Builder
	for _, k := range keys {
		key := opts.Prefix + k
		val := strings.ReplaceAll(secrets[k], `"`, `\"`)
		if opts.Export {
			fmt.Fprintf(&sb, "export %s=\"%s\"\n", key, val)
		} else {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", key, val)
		}
	}
	return sb.String()
}

func exportJSON(secrets map[string]string, keys []string, prefix string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		val := strings.ReplaceAll(secrets[k], `"`, `\"`)
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  \"%s%s\": \"%s\"%s\n", prefix, k, val, comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}

func exportYAML(secrets map[string]string, keys []string, prefix string) string {
	var sb strings.Builder
	for _, k := range keys {
		val := secrets[k]
		fmt.Fprintf(&sb, "%s%s: \"%s\"\n", prefix, k, strings.ReplaceAll(val, `"`, `\"`))
	}
	return sb.String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
