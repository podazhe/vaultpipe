package vault

import (
	"testing"
)

func TestDiffSecrets_Added(t *testing.T) {
	old := map[string]string{"A": "1"}
	next := map[string]string{"A": "1", "B": "2"}
	d := DiffSecrets(old, next)
	if len(d.Added) != 1 || d.Added["B"] != "2" {
		t.Fatalf("expected B added, got %v", d.Added)
	}
	if len(d.Removed) != 0 || len(d.Changed) != 0 {
		t.Fatal("unexpected removed or changed")
	}
}

func TestDiffSecrets_Removed(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1"}
	d := DiffSecrets(old, next)
	if len(d.Removed) != 1 || d.Removed["B"] != "2" {
		t.Fatalf("expected B removed, got %v", d.Removed)
	}
}

func TestDiffSecrets_Changed(t *testing.T) {
	old := map[string]string{"A": "old"}
	next := map[string]string{"A": "new"}
	d := DiffSecrets(old, next)
	if len(d.Changed) != 1 || d.Changed["A"] != "new" {
		t.Fatalf("expected A changed, got %v", d.Changed)
	}
}

func TestDiffSecrets_NoChanges(t *testing.T) {
	old := map[string]string{"A": "1"}
	next := map[string]string{"A": "1"}
	d := DiffSecrets(old, next)
	if d.HasChanges() {
		t.Fatal("expected no changes")
	}
}

func TestDiffSecrets_Summary(t *testing.T) {
	old := map[string]string{"A": "1"}
	next := map[string]string{"B": "2"}
	d := DiffSecrets(old, next)
	s := d.Summary()
	if s != "added=1 removed=1 changed=0" {
		t.Fatalf("unexpected summary: %s", s)
	}
}
