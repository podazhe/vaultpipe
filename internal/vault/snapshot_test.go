package vault

import (
	"os"
	"testing"
)

func TestSnapshotManager_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	m := NewSnapshotManager(dir)

	data := map[string]string{"API_KEY": "abc123", "DB_PASS": "secret"}
	if err := m.Save("myapp/config", data); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := m.Load("myapp/config")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap.Path != "myapp/config" {
		t.Errorf("expected path myapp/config, got %s", snap.Path)
	}
	if snap.Data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", snap.Data["API_KEY"])
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestSnapshotManager_LoadMissing(t *testing.T) {
	dir := t.TempDir()
	m := NewSnapshotManager(dir)
	_, err := m.Load("nonexistent/path")
	if err == nil {
		t.Fatal("expected error loading missing snapshot")
	}
}

func TestSnapshotManager_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	m := NewSnapshotManager(dir)

	_ = m.Save("app/db", map[string]string{"PASS": "old"})
	_ = m.Save("app/db", map[string]string{"PASS": "new"})

	snap, err := m.Load("app/db")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap.Data["PASS"] != "new" {
		t.Errorf("expected PASS=new, got %s", snap.Data["PASS"])
	}
}

func TestSnapshotManager_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	m := NewSnapshotManager(dir)
	_ = m.Save("sec/path", map[string]string{"K": "V"})

	info, err := os.Stat(m.filePath("sec/path"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 perms, got %v", info.Mode().Perm())
	}
}
