package vault

import (
	"testing"
)

func TestResolveSecrets_BasicVault(t *testing.T) {
	secrets := map[string]string{"db-password": "s3cr3t", "api-key": "abc"}
	resolved, err := ResolveSecrets(secrets, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resolved) != 2 {
		t.Fatalf("expected 2 resolved secrets, got %d", len(resolved))
	}
	for _, r := range resolved {
		if r.Source != "vault" {
			t.Errorf("expected source vault, got %s", r.Source)
		}
	}
}

func TestResolveSecrets_OverrideTakesPrecedence(t *testing.T) {
	secrets := map[string]string{"TOKEN": "original"}
	opts := ResolveOptions{Overrides: map[string]string{"TOKEN": "overridden"}}
	resolved, err := ResolveSecrets(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved[0].Value != "overridden" {
		t.Errorf("expected overridden, got %s", resolved[0].Value)
	}
	if resolved[0].Source != "override" {
		t.Errorf("expected source override, got %s", resolved[0].Source)
	}
}

func TestResolveSecrets_OverrideOnly(t *testing.T) {
	opts := ResolveOptions{Overrides: map[string]string{"EXTRA": "val"}}
	resolved, err := ResolveSecrets(nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resolved) != 1 || resolved[0].Key != "EXTRA" {
		t.Errorf("expected EXTRA key, got %+v", resolved)
	}
}

func TestResolveSecrets_NormaliseKey(t *testing.T) {
	secrets := map[string]string{"my-service/db-pass": "pw"}
	opts := ResolveOptions{BasePath: "my-service"}
	resolved, err := ResolveSecrets(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved[0].Key != "DB_PASS" {
		t.Errorf("expected DB_PASS, got %s", resolved[0].Key)
	}
}

func TestResolveSecrets_Empty(t *testing.T) {
	_, err := ResolveSecrets(nil, ResolveOptions{})
	if err == nil {
		t.Error("expected error for empty input")
	}
}
