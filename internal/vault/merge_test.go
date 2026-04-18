package vault

import (
	"testing"
)

func TestMergeSecrets_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}

	r := MergeSecrets(MergeStrategyFirst, a, b)
	if len(r.Secrets) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Secrets))
	}
	if len(r.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %v", r.Conflicts)
	}
}

func TestMergeSecrets_ConflictFirst(t *testing.T) {
	a := map[string]string{"FOO": "from-a"}
	b := map[string]string{"FOO": "from-b"}

	r := MergeSecrets(MergeStrategyFirst, a, b)
	if r.Secrets["FOO"] != "from-a" {
		t.Fatalf("expected 'from-a', got %q", r.Secrets["FOO"])
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0] != "FOO" {
		t.Fatalf("expected conflict on FOO, got %v", r.Conflicts)
	}
}

func TestMergeSecrets_ConflictLast(t *testing.T) {
	a := map[string]string{"FOO": "from-a"}
	b := map[string]string{"FOO": "from-b"}

	r := MergeSecrets(MergeStrategyLast, a, b)
	if r.Secrets["FOO"] != "from-b" {
		t.Fatalf("expected 'from-b', got %q", r.Secrets["FOO"])
	}
}

func TestMergeSecrets_EmptySources(t *testing.T) {
	r := MergeSecrets(MergeStrategyFirst)
	if len(r.Secrets) != 0 {
		t.Fatal("expected empty result")
	}
	if r.SourceCount != 0 {
		t.Fatal("expected source count 0")
	}
}

func TestMergeSecrets_ConflictDeduplication(t *testing.T) {
	a := map[string]string{"KEY": "a"}
	b := map[string]string{"KEY": "b"}
	c := map[string]string{"KEY": "c"}

	r := MergeSecrets(MergeStrategyFirst, a, b, c)
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict entry, got %d", len(r.Conflicts))
	}
}
