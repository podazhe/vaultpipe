package vault

import (
	"testing"
)

func baseTaggedSecrets() *TaggedSecrets {
	return NewTaggedSecrets(map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"APP_TOKEN":   "tok",
	})
}

func TestTaggedSecrets_TagAndFilter(t *testing.T) {
	ts := baseTaggedSecrets()

	if err := ts.Tag("DB_PASSWORD", "database", "sensitive"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ts.Tag("API_KEY", "sensitive"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := ts.FilterByTag("sensitive")
	if len(got) != 2 {
		t.Fatalf("expected 2 sensitive secrets, got %d", len(got))
	}
	if _, ok := got["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD in sensitive filter")
	}
	if _, ok := got["API_KEY"]; !ok {
		t.Error("expected API_KEY in sensitive filter")
	}
}

func TestTaggedSecrets_FilterByMultipleTags(t *testing.T) {
	ts := baseTaggedSecrets()
	_ = ts.Tag("DB_PASSWORD", "database", "sensitive")
	_ = ts.Tag("API_KEY", "sensitive")

	got := ts.FilterByTag("database", "sensitive")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if _, ok := got["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD to match both tags")
	}
}

func TestTaggedSecrets_TagMissingKey(t *testing.T) {
	ts := baseTaggedSecrets()
	err := ts.Tag("NONEXISTENT", "foo")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestTaggedSecrets_ListTags(t *testing.T) {
	ts := baseTaggedSecrets()
	_ = ts.Tag("DB_PASSWORD", "database", "sensitive")
	_ = ts.Tag("API_KEY", "sensitive", "external")

	tags := ts.ListTags()
	expected := []string{"database", "external", "sensitive"}
	if len(tags) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, tags)
	}
	for i, tag := range expected {
		if tags[i] != tag {
			t.Errorf("expected tags[%d]=%q, got %q", i, tag, tags[i])
		}
	}
}

func TestTaggedSecrets_DuplicateTagIgnored(t *testing.T) {
	ts := baseTaggedSecrets()
	_ = ts.Tag("API_KEY", "sensitive")
	_ = ts.Tag("API_KEY", "sensitive")

	if len(ts.Tags["API_KEY"]) != 1 {
		t.Errorf("expected 1 tag entry, got %d", len(ts.Tags["API_KEY"]))
	}
}
