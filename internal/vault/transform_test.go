package vault

import (
	"testing"
)

func TestTransformer_Prefix(t *testing.T) {
	tr := NewTransformer(TransformRule{Prefix: "app"})
	out, err := tr.Apply(map[string]string{"DB_PASS": "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if out["APP_DB_PASS"] != "secret" {
		t.Errorf("expected APP_DB_PASS, got %v", out)
	}
}

func TestTransformer_Rename(t *testing.T) {
	tr := NewTransformer(TransformRule{Renames: map[string]string{"old_key": "NEW_KEY"}})
	out, err := tr.Apply(map[string]string{"old_key": "val", "other": "x"})
	if err != nil {
		t.Fatal(err)
	}
	if out["NEW_KEY"] != "val" {
		t.Errorf("expected NEW_KEY=val, got %v", out)
	}
	if out["other"] != "x" {
		t.Errorf("expected other=x, got %v", out)
	}
}

func TestTransformer_Filter(t *testing.T) {
	tr := NewTransformer(TransformRule{Filter: []string{"keep"}})
	out, err := tr.Apply(map[string]string{"keep": "yes", "drop": "no"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["drop"]; ok {
		t.Error("drop should have been filtered out")
	}
	if out["keep"] != "yes" {
		t.Errorf("expected keep=yes, got %v", out)
	}
}

func TestTransformer_EmptyRenameTarget(t *testing.T) {
	tr := NewTransformer(TransformRule{Renames: map[string]string{"key": ""}})
	_, err := tr.Apply(map[string]string{"key": "val"})
	if err == nil {
		t.Error("expected error for empty rename target")
	}
}

func TestTransformer_PrefixAndRename(t *testing.T) {
	tr := NewTransformer(TransformRule{
		Prefix:  "svc",
		Renames: map[string]string{"token": "API_TOKEN"},
	})
	out, err := tr.Apply(map[string]string{"token": "abc"})
	if err != nil {
		t.Fatal(err)
	}
	if out["SVC_API_TOKEN"] != "abc" {
		t.Errorf("expected SVC_API_TOKEN=abc, got %v", out)
	}
}
