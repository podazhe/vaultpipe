package vault

import (
	"errors"
	"strings"
	"testing"
)

func TestChain_SingleStep(t *testing.T) {
	chain := NewChain().Add("upper", func(m map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[k] = strings.ToUpper(v)
		}
		return out, nil
	})

	result, err := chain.Run(map[string]string{"KEY": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "HELLO" {
		t.Errorf("expected HELLO, got %s", result["KEY"])
	}
}

func TestChain_MultipleSteps(t *testing.T) {
	chain := NewChain().
		Add("prefix", func(m map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(m))
			for k, v := range m {
				out["APP_"+k] = v
			}
			return out, nil
		}).
		Add("upper", func(m map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(m))
			for k, v := range m {
				out[k] = strings.ToUpper(v)
			}
			return out, nil
		})

	result, err := chain.Run(map[string]string{"DB": "postgres"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_DB"] != "POSTGRES" {
		t.Errorf("expected POSTGRES, got %s", result["APP_DB"])
	}
}

func TestChain_StepError(t *testing.T) {
	chain := NewChain().Add("fail", func(m map[string]string) (map[string]string, error) {
		return nil, errors.New("boom")
	})

	_, err := chain.Run(map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "fail") {
		t.Errorf("expected step name in error, got: %v", err)
	}
}

func TestChain_Steps(t *testing.T) {
	chain := NewChain().
		Add("a", func(m map[string]string) (map[string]string, error) { return m, nil }).
		Add("b", func(m map[string]string) (map[string]string, error) { return m, nil })

	steps := chain.Steps()
	if len(steps) != 2 || steps[0] != "a" || steps[1] != "b" {
		t.Errorf("unexpected steps: %v", steps)
	}
}

func TestChain_IsolatesInput(t *testing.T) {
	orig := map[string]string{"K": "v"}
	chain := NewChain().Add("mutate", func(m map[string]string) (map[string]string, error) {
		m["K"] = "changed"
		return m, nil
	})

	_, err := chain.Run(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orig["K"] != "v" {
		t.Error("original map was mutated")
	}
}
