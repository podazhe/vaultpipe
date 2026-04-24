package vault

import (
	"testing"
)

func baseUnsortedSecrets() map[string]string {
	return map[string]string{
		"ZEBRA":   "last",
		"APPLE":   "first",
		"MANGO":   "middle",
		"apricot": "lower",
	}
}

func TestSortSecrets_AscByKey(t *testing.T) {
	pairs := SortSecrets(baseUnsortedSecrets(), SortOptions{Order: SortAsc})
	if pairs[0].Key != "APPLE" {
		t.Fatalf("expected APPLE first, got %s", pairs[0].Key)
	}
	if pairs[len(pairs)-1].Key != "apricot" {
		t.Fatalf("expected apricot last, got %s", pairs[len(pairs)-1].Key)
	}
}

func TestSortSecrets_DescByKey(t *testing.T) {
	pairs := SortSecrets(baseUnsortedSecrets(), SortOptions{Order: SortDesc})
	if pairs[0].Key != "apricot" {
		t.Fatalf("expected apricot first in desc, got %s", pairs[0].Key)
	}
}

func TestSortSecrets_IgnoreCase(t *testing.T) {
	pairs := SortSecrets(baseUnsortedSecrets(), SortOptions{Order: SortAsc, IgnoreCase: true})
	// APPLE and apricot both start with 'a'; APPLE < apricot lexicographically when lowered
	if pairs[0].Key != "APPLE" && pairs[0].Key != "apricot" {
		t.Fatalf("unexpected first key: %s", pairs[0].Key)
	}
}

func TestSortSecrets_ByValue(t *testing.T) {
	pairs := SortSecrets(baseUnsortedSecrets(), SortOptions{Order: SortAsc, ByValue: true})
	// values: first, last, lower, middle
	if pairs[0].Value != "first" {
		t.Fatalf("expected 'first' value first, got %s", pairs[0].Value)
	}
	if pairs[len(pairs)-1].Value != "middle" {
		t.Fatalf("expected 'middle' value last, got %s", pairs[len(pairs)-1].Value)
	}
}

func TestSortSecrets_PrefixFilter(t *testing.T) {
	pairs := SortSecrets(baseUnsortedSecrets(), SortOptions{Order: SortAsc, Prefix: "A"})
	if len(pairs) != 1 {
		t.Fatalf("expected 1 result with prefix 'A', got %d", len(pairs))
	}
	if pairs[0].Key != "APPLE" {
		t.Fatalf("expected APPLE, got %s", pairs[0].Key)
	}
}

func TestSortSecrets_Empty(t *testing.T) {
	pairs := SortSecrets(map[string]string{}, SortOptions{})
	if len(pairs) != 0 {
		t.Fatalf("expected empty result, got %d items", len(pairs))
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	pairs := SortSecrets(src, SortOptions{Order: SortAsc})
	out := ToMap(pairs)
	for k, v := range src {
		if out[k] != v {
			t.Errorf("key %s: expected %s got %s", k, v, out[k])
		}
	}
}
