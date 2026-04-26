package vault

import (
	"testing"
)

func TestRenameSecrets_ExplicitRule(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "value1", "KEEP": "value2"}
	res, err := RenameSecrets(secrets, RenameOptions{Rules: map[string]string{"OLD_KEY": "NEW_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", res.Secrets["NEW_KEY"])
	}
	if _, ok := res.Secrets["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if len(res.Renamed) != 1 || res.Renamed[0] != "NEW_KEY" {
		t.Errorf("unexpected renamed list: %v", res.Renamed)
	}
}

func TestRenameSecrets_MissingKeySkipped(t *testing.T) {
	secrets := map[string]string{"KEEP": "v"}
	res, err := RenameSecrets(secrets, RenameOptions{Rules: map[string]string{"MISSING": "NEW"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestRenameSecrets_FailMissing(t *testing.T) {
	secrets := map[string]string{"KEEP": "v"}
	_, err := RenameSecrets(secrets, RenameOptions{
		Rules:       map[string]string{"GHOST": "NEW"},
		FailMissing: true,
	})
	if err == nil {
		t.Fatal("expected error for missing key with FailMissing=true")
	}
}

func TestRenameSecrets_PrefixSubstitution(t *testing.T) {
	secrets := map[string]string{"APP_FOO": "1", "APP_BAR": "2", "OTHER": "3"}
	res, err := RenameSecrets(secrets, RenameOptions{OldPrefix: "APP_", NewPrefix: "SVC_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["SVC_FOO"] != "1" || res.Secrets["SVC_BAR"] != "2" {
		t.Errorf("prefix rename failed: %v", res.Secrets)
	}
	if res.Secrets["OTHER"] != "3" {
		t.Error("non-matching key should be unchanged")
	}
}

func TestRenameSecrets_DryRun(t *testing.T) {
	secrets := map[string]string{"OLD": "val"}
	res, err := RenameSecrets(secrets, RenameOptions{
		Rules:  map[string]string{"OLD": "NEW"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// DryRun: original key must still be present.
	if _, ok := res.Secrets["OLD"]; !ok {
		t.Error("DryRun should not remove OLD key")
	}
	if len(res.Renamed) != 1 {
		t.Errorf("expected 1 renamed entry, got %d", len(res.Renamed))
	}
}

func TestRenameSecrets_EmptyTargetErrors(t *testing.T) {
	secrets := map[string]string{"KEY": "v"}
	_, err := RenameSecrets(secrets, RenameOptions{Rules: map[string]string{"KEY": ""}})
	if err == nil {
		t.Fatal("expected error for empty target key")
	}
}

func TestRenameSecrets_NilSecrets(t *testing.T) {
	_, err := RenameSecrets(nil, RenameOptions{})
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}
