package vault

import (
	"strings"
	"testing"
)

// buildChainFromNames is a test helper that constructs a Chain from a slice of
// step names using the same built-in mapping as cmd/chain.go.
func buildChainFromNames(t *testing.T, names []string, prefix string) *Chain {
	t.Helper()
	chain := NewChain()
	for _, name := range names {
		switch strings.ToLower(name) {
		case "prefix":
			p := prefix
			chain.Add("prefix", func(m map[string]string) (map[string]string, error) {
				out := make(map[string]string, len(m))
				for k, v := range m {
					out[p+k] = v
				}
				return out, nil
			})
		case "upper":
			chain.Add("upper", func(m map[string]string) (map[string]string, error) {
				out := make(map[string]string, len(m))
				for k, v := range m {
					out[k] = strings.ToUpper(v)
				}
				return out, nil
			})
		default:
			t.Fatalf("unknown step name: %q", name)
		}
	}
	return chain
}

func TestBuildChainFromNames_PrefixThenUpper(t *testing.T) {
	chain := buildChainFromNames(t, []string{"prefix", "upper"}, "SVC_")
	result, err := chain.Run(map[string]string{"host": "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["SVC_host"] != "LOCALHOST" {
		t.Errorf("expected LOCALHOST, got %q", result["SVC_host"])
	}
}

func TestBuildChainFromNames_StepsPreserved(t *testing.T) {
	chain := buildChainFromNames(t, []string{"prefix", "upper"}, "X_")
	steps := chain.Steps()
	if len(steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(steps))
	}
	if steps[0] != "prefix" || steps[1] != "upper" {
		t.Errorf("unexpected step order: %v", steps)
	}
}
