package vault

import (
	"testing"
)

func TestTruncateSecrets_BasicTruncation(t *testing.T) {
	secrets := map[string]string{
		"SHORT": "hi",
		"LONG":  "this-value-is-definitely-too-long",
	}
	out, result, err := TruncateSecrets(secrets, TruncateOptions{MaxLen: 10, Suffix: "..."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SHORT"] != "hi" {
		t.Errorf("expected SHORT unchanged, got %q", out["SHORT"])
	}
	if len(out["LONG"]) != 10 {
		t.Errorf("expected LONG length 10, got %d", len(out["LONG"]))
	}
	if out["LONG"] != "this-va..." {
		t.Errorf("unexpected truncated value: %q", out["LONG"])
	}
	if len(result.Truncated) != 1 || result.Truncated[0] != "LONG" {
		t.Errorf("expected [LONG] in Truncated, got %v", result.Truncated)
	}
}

func TestTruncateSecrets_NoSuffix(t *testing.T) {
	secrets := map[string]string{"KEY": "abcdefghij"}
	out, _, err := TruncateSecrets(secrets, TruncateOptions{MaxLen: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "abcde" {
		t.Errorf("expected \"abcde\", got %q", out["KEY"])
	}
}

func TestTruncateSecrets_KeyFilter(t *testing.T) {
	secrets := map[string]string{
		"A": "this-is-long-value",
		"B": "also-a-long-value",
	}
	out, result, err := TruncateSecrets(secrets, TruncateOptions{MaxLen: 8, Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out["A"]) != 8 {
		t.Errorf("expected A truncated to 8, got %d", len(out["A"]))
	}
	if out["B"] != "also-a-long-value" {
		t.Errorf("expected B untouched, got %q", out["B"])
	}
	if len(result.Truncated) != 1 {
		t.Errorf("expected 1 truncated key, got %d", len(result.Truncated))
	}
}

func TestTruncateSecrets_DryRun(t *testing.T) {
	secrets := map[string]string{"X": "longvaluethatwillbetruncated"}
	out, result, err := TruncateSecrets(secrets, TruncateOptions{MaxLen: 5, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != secrets["X"] {
		t.Errorf("dry run should not modify values, got %q", out["X"])
	}
	if len(result.Truncated) != 1 {
		t.Errorf("expected truncated key reported in dry run")
	}
}

func TestTruncateSecrets_InvalidMaxLen(t *testing.T) {
	_, _, err := TruncateSecrets(map[string]string{"K": "v"}, TruncateOptions{MaxLen: 0})
	if err == nil {
		t.Error("expected error for MaxLen=0")
	}
}

func TestTruncateSecrets_SuffixTooLong(t *testing.T) {
	_, _, err := TruncateSecrets(map[string]string{"K": "v"}, TruncateOptions{MaxLen: 3, Suffix: "..."})
	if err == nil {
		t.Error("expected error when suffix >= MaxLen")
	}
}

func TestTruncateSecrets_NothingToTruncate(t *testing.T) {
	secrets := map[string]string{"A": "short", "B": "ok"}
	_, result, err := TruncateSecrets(secrets, TruncateOptions{MaxLen: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Truncated) != 0 {
		t.Errorf("expected no truncated keys, got %v", result.Truncated)
	}
}
