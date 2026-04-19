package vault

import (
	"testing"
)

func TestRedactSecrets_DefaultPrefixes(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PASSWORD":     "s3cr3t",
		"TOKEN_VALUE":  "abc123",
		"APP_NAME":     "vaultpipe",
	}
	out := RedactSecrets(secrets, RedactOptions{})
	if out["DATABASE_URL"] != "postgres://localhost/db" {
		t.Errorf("expected DATABASE_URL unchanged, got %s", out["DATABASE_URL"])
	}
	if out["PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected PASSWORD redacted, got %s", out["PASSWORD"])
	}
	if out["TOKEN_VALUE"] != "***REDACTED***" {
		t.Errorf("expected TOKEN_VALUE redacted, got %s", out["TOKEN_VALUE"])
	}
	if out["APP_NAME"] != "vaultpipe" {
		t.Errorf("expected APP_NAME unchanged, got %s", out["APP_NAME"])
	}
}

func TestRedactSecrets_ExactKeys(t *testing.T) {
	secrets := map[string]string{
		"MY_CUSTOM": "hidden",
		"SAFE_KEY":  "visible",
	}
	out := RedactSecrets(secrets, RedactOptions{SensitiveKeys: []string{"MY_CUSTOM"}})
	if out["MY_CUSTOM"] != "***REDACTED***" {
		t.Errorf("expected MY_CUSTOM redacted")
	}
	if out["SAFE_KEY"] != "visible" {
		t.Errorf("expected SAFE_KEY visible")
	}
}

func TestRedactSecrets_CustomPrefixes(t *testing.T) {
	secrets := map[string]string{
		"INTERNAL_VALUE": "secret",
		"PUBLIC_DATA":    "open",
	}
	opts := RedactOptions{SensitivePrefixes: []string{"INTERNAL"}}
	out := RedactSecrets(secrets, opts)
	if out["INTERNAL_VALUE"] != "***REDACTED***" {
		t.Errorf("expected INTERNAL_VALUE redacted")
	}
	if out["PUBLIC_DATA"] != "open" {
		t.Errorf("expected PUBLIC_DATA unchanged")
	}
}

func TestRedactSecrets_Empty(t *testing.T) {
	out := RedactSecrets(map[string]string{}, RedactOptions{})
	if len(out) != 0 {
		t.Errorf("expected empty map")
	}
}
