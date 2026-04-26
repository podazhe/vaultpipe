package vault

import (
	"testing"
)

func basePruneSecrets() map[string]string {
	return map[string]string{
		"APP_KEY":     "abc123",
		"DB_PASSWORD": "",
		"LEGACY_HOST": "old.host",
		"LEGACY_PORT": "5432",
		"DEBUG_MODE":  "true",
	}
}

func TestPruneSecrets_RemoveEmpty(t *testing.T) {
	out, res, err := PruneSecrets(basePruneSecrets(), PruneOptions{RemoveEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("expected DB_PASSWORD to be pruned")
	}
	if res.Retained != 4 {
		t.Errorf("expected 4 retained, got %d", res.Retained)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "DB_PASSWORD" {
		t.Errorf("unexpected removed list: %v", res.Removed)
	}
}

func TestPruneSecrets_RemoveByPrefix(t *testing.T) {
	out, res, err := PruneSecrets(basePruneSecrets(), PruneOptions{RemovePrefixes: []string{"LEGACY_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["LEGACY_HOST"]; ok {
		t.Error("expected LEGACY_HOST to be pruned")
	}
	if _, ok := out["LEGACY_PORT"]; ok {
		t.Error("expected LEGACY_PORT to be pruned")
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestPruneSecrets_RemoveByKey(t *testing.T) {
	out, res, err := PruneSecrets(basePruneSecrets(), PruneOptions{RemoveKeys: []string{"DEBUG_MODE", "APP_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DEBUG_MODE"]; ok {
		t.Error("expected DEBUG_MODE to be pruned")
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestPruneSecrets_DryRun(t *testing.T) {
	orig := basePruneSecrets()
	out, res, err := PruneSecrets(orig, PruneOptions{RemoveEmpty: true, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASSWORD"]; !ok {
		t.Error("dry-run should not remove DB_PASSWORD from returned map")
	}
	if len(res.Removed) != 1 {
		t.Errorf("dry-run should still report 1 removal, got %d", len(res.Removed))
	}
}

func TestPruneSecrets_NilErrors(t *testing.T) {
	_, _, err := PruneSecrets(nil, PruneOptions{})
	if err == nil {
		t.Error("expected error for nil secrets")
	}
}

func TestPruneSecrets_Summary(t *testing.T) {
	_, res, _ := PruneSecrets(basePruneSecrets(), PruneOptions{RemoveEmpty: true})
	s := res.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
