package vault

import (
	"testing"
)

// buildIndexFromMap is a test helper that populates a SecretIndex from a
// map of key->tags, all registered under a single path.
func buildIndexFromMap(path string, data map[string][]string) *SecretIndex {
	idx := NewSecretIndex()
	for key, tags := range data {
		idx.Add(path, key, tags, nil)
	}
	return idx
}

func TestBuildIndexFromMap_CountsCorrectly(t *testing.T) {
	idx := buildIndexFromMap("secret/test", map[string][]string{
		"KEY_A": {"tag1"},
		"KEY_B": {},
		"KEY_C": {"tag1", "tag2"},
	})

	if idx.Size() != 3 {
		t.Fatalf("expected 3, got %d", idx.Size())
	}
}

func TestBuildIndexFromMap_TagsPreserved(t *testing.T) {
	idx := buildIndexFromMap("secret/test", map[string][]string{
		"SENSITIVE": {"pii", "sensitive"},
	})

	e, ok := idx.Get("secret/test", "SENSITIVE")
	if !ok {
		t.Fatal("entry missing")
	}
	if len(e.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(e.Tags))
	}
}

func TestBuildIndexFromMap_IsolatesCopy(t *testing.T) {
	src := map[string][]string{"K": {"t"}}
	idx := buildIndexFromMap("p", src)
	// mutate original — index should be unaffected
	src["K"] = append(src["K"], "extra")

	e, _ := idx.Get("p", "K")
	if len(e.Tags) != 1 {
		t.Errorf("index was not isolated from source mutation: %v", e.Tags)
	}
}
