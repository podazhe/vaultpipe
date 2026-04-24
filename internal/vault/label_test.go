package vault

import (
	"testing"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"TOKEN":       "tok-xyz",
	}
}

func TestLabeledSecrets_AddAndGet(t *testing.T) {
	ls := NewLabeledSecrets(baseSecrets())
	if err := ls.AddLabel("DB_PASSWORD", "env", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lbls := ls.GetLabels("DB_PASSWORD")
	if lbls["env"] != "production" {
		t.Errorf("expected label env=production, got %v", lbls)
	}
}

func TestLabeledSecrets_MissingSecretKey(t *testing.T) {
	ls := NewLabeledSecrets(baseSecrets())
	err := ls.AddLabel("DOES_NOT_EXIST", "env", "staging")
	if err == nil {
		t.Fatal("expected error for missing secret key, got nil")
	}
}

func TestLabeledSecrets_EmptyLabelKey(t *testing.T) {
	ls := NewLabeledSecrets(baseSecrets())
	err := ls.AddLabel("API_KEY", "", "value")
	if err == nil {
		t.Fatal("expected error for empty label key")
	}
}

func TestLabeledSecrets_FilterByLabel(t *testing.T) {
	ls := NewLabeledSecrets(baseSecrets())
	_ = ls.AddLabel("DB_PASSWORD", "sensitivity", "high")
	_ = ls.AddLabel("API_KEY", "sensitivity", "high")
	_ = ls.AddLabel("TOKEN", "sensitivity", "medium")

	high := ls.FilterByLabel("sensitivity", "high")
	if len(high) != 2 {
		t.Errorf("expected 2 high-sensitivity secrets, got %d", len(high))
	}
	if _, ok := high["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD in filtered results")
	}
	if _, ok := high["API_KEY"]; !ok {
		t.Error("expected API_KEY in filtered results")
	}
}

func TestLabeledSecrets_ListLabeled(t *testing.T) {
	ls := NewLabeledSecrets(baseSecrets())
	_ = ls.AddLabel("TOKEN", "owner", "team-a")
	_ = ls.AddLabel("API_KEY", "owner", "team-b")

	keys := ls.ListLabeled()
	if len(keys) != 2 {
		t.Errorf("expected 2 labeled keys, got %d", len(keys))
	}
	// should be sorted
	if keys[0] != "API_KEY" || keys[1] != "TOKEN" {
		t.Errorf("expected sorted keys [API_KEY TOKEN], got %v", keys)
	}
}

func TestLabeledSecrets_RemoveLabel(t *testing.T) {
	ls := NewLabeledSecrets(baseSecrets())
	_ = ls.AddLabel("DB_PASSWORD", "env", "production")

	if err := ls.RemoveLabel("DB_PASSWORD", "env"); err != nil {
		t.Fatalf("unexpected error on RemoveLabel: %v", err)
	}
	if len(ls.GetLabels("DB_PASSWORD")) != 0 {
		t.Error("expected no labels after removal")
	}
}

func TestLabeledSecrets_RemoveLabel_Missing(t *testing.T) {
	ls := NewLabeledSecrets(baseSecrets())
	err := ls.RemoveLabel("API_KEY", "nonexistent")
	if err == nil {
		t.Fatal("expected error when removing non-existent label")
	}
}
