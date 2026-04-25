package vault

import (
	"testing"
)

func basePatchSecrets() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
	}
}

func TestPatchSecrets_SetNewKey(t *testing.T) {
	ops := []PatchOp{{Op: "set", Key: "NEW_KEY", Value: "newval"}}
	res, err := PatchSecrets(basePatchSecrets(), ops, PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Final["NEW_KEY"] != "newval" {
		t.Errorf("expected NEW_KEY=newval, got %q", res.Final["NEW_KEY"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied op, got %d", len(res.Applied))
	}
}

func TestPatchSecrets_DeleteExisting(t *testing.T) {
	ops := []PatchOp{{Op: "delete", Key: "API_KEY"}}
	res, err := PatchSecrets(basePatchSecrets(), ops, PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Final["API_KEY"]; ok {
		t.Error("expected API_KEY to be deleted")
	}
}

func TestPatchSecrets_DeleteMissing(t *testing.T) {
	ops := []PatchOp{{Op: "delete", Key: "NONEXISTENT"}}
	res, err := PatchSecrets(basePatchSecrets(), ops, PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped op, got %d", len(res.Skipped))
	}
}

func TestPatchSecrets_Rename(t *testing.T) {
	ops := []PatchOp{{Op: "rename", Key: "DB_HOST", Value: "DATABASE_HOST"}}
	res, err := PatchSecrets(basePatchSecrets(), ops, PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Final["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", res.Final["DATABASE_HOST"])
	}
	if _, ok := res.Final["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed after rename")
	}
}

func TestPatchSecrets_DryRun(t *testing.T) {
	ops := []PatchOp{
		{Op: "set", Key: "DB_HOST", Value: "changed"},
		{Op: "delete", Key: "API_KEY"},
	}
	res, err := PatchSecrets(basePatchSecrets(), ops, PatchOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Final["DB_HOST"] != "localhost" {
		t.Error("dry run should not modify DB_HOST")
	}
	if _, ok := res.Final["API_KEY"]; !ok {
		t.Error("dry run should not delete API_KEY")
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied (dry) ops, got %d", len(res.Applied))
	}
}

func TestPatchSecrets_UnknownOp(t *testing.T) {
	ops := []PatchOp{{Op: "upsert", Key: "X"}}
	_, err := PatchSecrets(basePatchSecrets(), ops, PatchOptions{})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestPatchSecrets_NilSource(t *testing.T) {
	_, err := PatchSecrets(nil, nil, PatchOptions{})
	if err == nil {
		t.Error("expected error for nil source")
	}
}
