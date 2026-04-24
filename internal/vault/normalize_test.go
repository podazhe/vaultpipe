package vault

import (
	"testing"
)

func TestNormalizeSecrets_UppercaseKeys(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost", "api_key": "abc"}
	res := NormalizeSecrets(secrets, NormalizeOptions{UppercaseKeys: true})
	if _, ok := res.Normalized["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in normalized map")
	}
	if len(res.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(res.Changes))
	}
}

func TestNormalizeSecrets_TrimValues(t *testing.T) {
	secrets := map[string]string{"key": "  value  "}
	res := NormalizeSecrets(secrets, NormalizeOptions{TrimValues: true})
	if got := res.Normalized["key"]; got != "value" {
		t.Errorf("expected 'value', got %q", got)
	}
	if len(res.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(res.Changes))
	}
}

func TestNormalizeSecrets_ReplaceHyphens(t *testing.T) {
	secrets := map[string]string{"my-secret-key": "val"}
	res := NormalizeSecrets(secrets, NormalizeOptions{ReplaceHyphens: true})
	if _, ok := res.Normalized["my_secret_key"]; !ok {
		t.Error("expected my_secret_key after hyphen replacement")
	}
}

func TestNormalizeSecrets_StripNonPrint(t *testing.T) {
	secrets := map[string]string{"key": "val\x00ue"}
	res := NormalizeSecrets(secrets, NormalizeOptions{StripNonPrint: true})
	if got := res.Normalized["key"]; got != "value" {
		t.Errorf("expected 'value' after stripping non-printable, got %q", got)
	}
}

func TestNormalizeSecrets_NoChanges(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	res := NormalizeSecrets(secrets, NormalizeOptions{UppercaseKeys: true, TrimValues: true})
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(res.Changes))
	}
}

func TestNormalizeSecrets_CombinedOptions(t *testing.T) {
	secrets := map[string]string{"my-key": "  hello  "}
	res := NormalizeSecrets(secrets, NormalizeOptions{
		UppercaseKeys:  true,
		TrimValues:     true,
		ReplaceHyphens: true,
	})
	if got, ok := res.Normalized["MY_KEY"]; !ok || got != "hello" {
		t.Errorf("expected MY_KEY=hello, got %q=%q", "MY_KEY", got)
	}
}

func TestNormalizeSecrets_EmptyMap(t *testing.T) {
	res := NormalizeSecrets(map[string]string{}, NormalizeOptions{UppercaseKeys: true})
	if len(res.Normalized) != 0 {
		t.Error("expected empty normalized map")
	}
	if len(res.Changes) != 0 {
		t.Error("expected no changes for empty input")
	}
}
