package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveManager_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	am, err := NewArchiveManager(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc123"}
	entry, err := am.Save("myapp/prod", secrets, "initial archive")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	if entry.Path != "myapp/prod" {
		t.Errorf("expected path myapp/prod, got %s", entry.Path)
	}
	if entry.Note != "initial archive" {
		t.Errorf("expected note 'initial archive', got %s", entry.Note)
	}
	if entry.Secrets["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t")
	}
}

func TestArchiveManager_LoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	am, _ := NewArchiveManager(dir)

	original := map[string]string{"TOKEN": "tok_abc"}
	saved, err := am.Save("service/token", original, "")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	file := am.fileName(saved.Path, saved.ArchivedAt)
	loaded, err := am.Load(file)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Secrets["TOKEN"] != "tok_abc" {
		t.Errorf("expected TOKEN=tok_abc after round-trip")
	}
	if !loaded.ArchivedAt.Equal(saved.ArchivedAt) {
		t.Errorf("timestamp mismatch after round-trip")
	}
}

func TestArchiveManager_LoadMissing(t *testing.T) {
	dir := t.TempDir()
	am, _ := NewArchiveManager(dir)

	_, err := am.Load(filepath.Join(dir, "nonexistent.json"))
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestArchiveManager_EmptyPathErrors(t *testing.T) {
	dir := t.TempDir()
	am, _ := NewArchiveManager(dir)

	_, err := am.Save("", map[string]string{"K": "V"}, "")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestArchiveManager_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	am, _ := NewArchiveManager(dir)

	entry, err := am.Save("perm/test", map[string]string{"X": "1"}, "")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	file := am.fileName(entry.Path, entry.ArchivedAt)
	info, err := os.Stat(file)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file perm 0600, got %o", perm)
	}
}
