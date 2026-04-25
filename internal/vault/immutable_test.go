package vault

import (
	"testing"
)

func baseImmutableSecrets() *ImmutableSecrets {
	return NewImmutableSecrets(map[string]string{
		"DB_PASS": "s3cr3t",
		"API_KEY": "abc123",
		"DEBUG":   "true",
	})
}

func TestImmutableSecrets_GetAndSet(t *testing.T) {
	im := baseImmutableSecrets()
	if err := im.Set("DEBUG", "false"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := im.Get("DEBUG")
	if !ok || v != "false" {
		t.Errorf("expected DEBUG=false, got %q (found=%v)", v, ok)
	}
}

func TestImmutableSecrets_LockPreventsSet(t *testing.T) {
	im := baseImmutableSecrets()
	im.Lock("DB_PASS")
	if err := im.Set("DB_PASS", "newpass"); err == nil {
		t.Error("expected error when setting locked key, got nil")
	}
	v, _ := im.Get("DB_PASS")
	if v != "s3cr3t" {
		t.Errorf("locked value should be unchanged, got %q", v)
	}
}

func TestImmutableSecrets_LockPreventsDelete(t *testing.T) {
	im := baseImmutableSecrets()
	im.Lock("API_KEY")
	if err := im.Delete("API_KEY"); err == nil {
		t.Error("expected error when deleting locked key, got nil")
	}
}

func TestImmutableSecrets_DeleteUnlocked(t *testing.T) {
	im := baseImmutableSecrets()
	if err := im.Delete("DEBUG"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := im.Get("DEBUG"); ok {
		t.Error("expected DEBUG to be deleted")
	}
}

func TestImmutableSecrets_IsLocked(t *testing.T) {
	im := baseImmutableSecrets()
	im.Lock("DB_PASS")
	if !im.IsLocked("DB_PASS") {
		t.Error("expected DB_PASS to be locked")
	}
	if im.IsLocked("API_KEY") {
		t.Error("expected API_KEY to be unlocked")
	}
}

func TestImmutableSecrets_SnapshotIsolated(t *testing.T) {
	im := baseImmutableSecrets()
	snap := im.Snapshot()
	snap["DB_PASS"] = "mutated"
	v, _ := im.Get("DB_PASS")
	if v == "mutated" {
		t.Error("snapshot mutation should not affect original")
	}
}

func TestImmutableSecrets_LockedKeys(t *testing.T) {
	im := baseImmutableSecrets()
	im.Lock("DB_PASS", "API_KEY")
	keys := im.LockedKeys()
	if len(keys) != 2 {
		t.Errorf("expected 2 locked keys, got %d", len(keys))
	}
}

func TestImmutableSecrets_InputIsolated(t *testing.T) {
	orig := map[string]string{"TOKEN": "xyz"}
	im := NewImmutableSecrets(orig)
	orig["TOKEN"] = "changed"
	v, _ := im.Get("TOKEN")
	if v != "xyz" {
		t.Errorf("expected original value xyz, got %q", v)
	}
}
