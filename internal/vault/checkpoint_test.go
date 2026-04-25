package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckpointManager_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	m, err := NewCheckpointManager(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{"DB_PASS": "secret123", "API_KEY": "abc"}
	if err := m.Save("v1", secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}
	cp, err := m.Load("v1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cp.Name != "v1" {
		t.Errorf("name: got %q, want %q", cp.Name, "v1")
	}
	if cp.Secrets["DB_PASS"] != "secret123" {
		t.Errorf("secret mismatch")
	}
	if cp.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestCheckpointManager_LoadMissing(t *testing.T) {
	dir := t.TempDir()
	m, _ := NewCheckpointManager(dir)
	_, err := m.Load("ghost")
	if err == nil {
		t.Fatal("expected error for missing checkpoint")
	}
}

func TestCheckpointManager_Delete(t *testing.T) {
	dir := t.TempDir()
	m, _ := NewCheckpointManager(dir)
	_ = m.Save("tmp", map[string]string{"K": "V"})
	if err := m.Delete("tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	path := filepath.Join(dir, "tmp.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestCheckpointManager_DeleteMissing(t *testing.T) {
	dir := t.TempDir()
	m, _ := NewCheckpointManager(dir)
	if err := m.Delete("nope"); err == nil {
		t.Error("expected error deleting non-existent checkpoint")
	}
}

func TestCheckpointManager_EmptyName(t *testing.T) {
	dir := t.TempDir()
	m, _ := NewCheckpointManager(dir)
	if err := m.Save("", map[string]string{}); err == nil {
		t.Error("expected error for empty name on Save")
	}
	if _, err := m.Load(""); err == nil {
		t.Error("expected error for empty name on Load")
	}
}

func TestCheckpointManager_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	m, _ := NewCheckpointManager(dir)
	_ = m.Save("secure", map[string]string{"TOKEN": "xyz"})
	path := filepath.Join(dir, "secure.json")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("permissions: got %o, want 0600", perm)
	}
}
