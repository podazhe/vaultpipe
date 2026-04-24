package vault

import (
	"testing"
)

func TestSecretIndex_AddAndGet(t *testing.T) {
	idx := NewSecretIndex()
	idx.Add("secret/app", "DB_PASS", []string{"sensitive"}, map[string]string{"env": "prod"})

	e, ok := idx.Get("secret/app", "DB_PASS")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Key != "DB_PASS" || e.Path != "secret/app" {
		t.Errorf("unexpected entry: %+v", e)
	}
	if len(e.Tags) != 1 || e.Tags[0] != "sensitive" {
		t.Errorf("unexpected tags: %v", e.Tags)
	}
	if e.Labels["env"] != "prod" {
		t.Errorf("unexpected labels: %v", e.Labels)
	}
}

func TestSecretIndex_Remove(t *testing.T) {
	idx := NewSecretIndex()
	idx.Add("secret/app", "API_KEY", nil, nil)
	idx.Remove("secret/app", "API_KEY")

	_, ok := idx.Get("secret/app", "API_KEY")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestSecretIndex_ListByPath(t *testing.T) {
	idx := NewSecretIndex()
	idx.Add("secret/app", "Z_KEY", nil, nil)
	idx.Add("secret/app", "A_KEY", nil, nil)
	idx.Add("secret/other", "X_KEY", nil, nil)

	list := idx.ListByPath("secret/app")
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
	if list[0].Key != "A_KEY" || list[1].Key != "Z_KEY" {
		t.Errorf("unexpected order: %v %v", list[0].Key, list[1].Key)
	}
}

func TestSecretIndex_Search(t *testing.T) {
	idx := NewSecretIndex()
	idx.Add("secret/app", "DB_HOST", nil, nil)
	idx.Add("secret/app", "DB_PASS", nil, nil)
	idx.Add("secret/app", "API_KEY", nil, nil)

	results := idx.Search("DB_")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSecretIndex_Size(t *testing.T) {
	idx := NewSecretIndex()
	if idx.Size() != 0 {
		t.Fatal("expected empty index")
	}
	idx.Add("p", "k1", nil, nil)
	idx.Add("p", "k2", nil, nil)
	if idx.Size() != 2 {
		t.Fatalf("expected size 2, got %d", idx.Size())
	}
}

func TestSecretIndex_GetMissing(t *testing.T) {
	idx := NewSecretIndex()
	_, ok := idx.Get("secret/app", "MISSING")
	if ok {
		t.Fatal("expected miss")
	}
}
