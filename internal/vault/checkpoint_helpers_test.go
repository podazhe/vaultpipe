package vault

import (
	"testing"
)

// buildCheckpointFromMap creates and saves a checkpoint, returning the manager.
func buildCheckpointFromMap(t *testing.T, dir, name string, secrets map[string]string) *CheckpointManager {
	t.Helper()
	m, err := NewCheckpointManager(dir)
	if err != nil {
		t.Fatalf("NewCheckpointManager: %v", err)
	}
	if err := m.Save(name, secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return m
}

func TestBuildCheckpointFromMap_LoadsCorrectly(t *testing.T) {
	dir := t.TempDir()
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	m := buildCheckpointFromMap(t, dir, "test", secrets)
	cp, err := m.Load("test")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cp.Secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(cp.Secrets))
	}
}

func TestBuildCheckpointFromMap_IsolatesCopy(t *testing.T) {
	dir := t.TempDir()
	secrets := map[string]string{"KEY": "original"}
	m := buildCheckpointFromMap(t, dir, "iso", secrets)
	secrets["KEY"] = "mutated"
	cp, _ := m.Load("iso")
	if cp.Secrets["KEY"] != "original" {
		t.Errorf("checkpoint should not reflect mutation: got %q", cp.Secrets["KEY"])
	}
}

func TestBuildCheckpointFromMap_NamePreserved(t *testing.T) {
	dir := t.TempDir()
	m := buildCheckpointFromMap(t, dir, "mycp", map[string]string{})
	cp, _ := m.Load("mycp")
	if cp.Name != "mycp" {
		t.Errorf("name: got %q, want %q", cp.Name, "mycp")
	}
}
