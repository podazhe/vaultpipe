package vault

import (
	"testing"
)

func TestPromoteSecrets_NoConflicts(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"C": "3"}

	out, res := PromoteSecrets(src, dst, PromoteOptions{})

	if out["A"] != "1" || out["B"] != "2" || out["C"] != "3" {
		t.Fatalf("unexpected output: %v", out)
	}
	if len(res.Promoted) != 2 || len(res.Skipped) != 0 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestPromoteSecrets_SkipsExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	out, res := PromoteSecrets(src, dst, PromoteOptions{Overwrite: false})

	if out["A"] != "old" {
		t.Fatalf("expected old value, got %s", out["A"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Fatalf("expected skip, got %+v", res)
	}
}

func TestPromoteSecrets_OverwriteExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	out, res := PromoteSecrets(src, dst, PromoteOptions{Overwrite: true})

	if out["A"] != "new" {
		t.Fatalf("expected new value, got %s", out["A"])
	}
	if len(res.Overwrite) != 1 {
		t.Fatalf("expected overwrite record, got %+v", res)
	}
}

func TestPromoteSecrets_DryRun(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}

	out, res := PromoteSecrets(src, dst, PromoteOptions{DryRun: true})

	if _, ok := out["A"]; ok {
		t.Fatal("dry-run should not write to output")
	}
	if len(res.Promoted) != 1 {
		t.Fatalf("expected promoted entry in result, got %+v", res)
	}
}

func TestPromoteSecrets_KeyFilter(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}

	out, res := PromoteSecrets(src, dst, PromoteOptions{Keys: []string{"A", "C"}})

	if _, ok := out["B"]; ok {
		t.Fatal("B should not have been promoted")
	}
	if len(res.Promoted) != 2 {
		t.Fatalf("expected 2 promoted, got %+v", res)
	}
}

func TestPromoteResult_Summary(t *testing.T) {
	r := PromoteResult{
		Promoted:  []string{"A", "B"},
		Skipped:   []string{"C"},
		Overwrite: []string{},
	}
	got := r.Summary()
	want := "promoted=2 skipped=1 overwritten=0"
	if got != want {
		t.Fatalf("want %q got %q", want, got)
	}
}
