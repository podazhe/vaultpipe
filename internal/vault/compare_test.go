package vault

import (
	"strings"
	"testing"
)

func TestCompareSecrets_OnlyInLeft(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1"}

	r := CompareSecrets(left, right)
	if len(r.OnlyInLeft) != 1 || r.OnlyInLeft[0] != "B" {
		t.Errorf("expected OnlyInLeft=[B], got %v", r.OnlyInLeft)
	}
	if len(r.OnlyInRight) != 0 {
		t.Errorf("expected no OnlyInRight, got %v", r.OnlyInRight)
	}
}

func TestCompareSecrets_OnlyInRight(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "C": "3"}

	r := CompareSecrets(left, right)
	if len(r.OnlyInRight) != 1 || r.OnlyInRight[0] != "C" {
		t.Errorf("expected OnlyInRight=[C], got %v", r.OnlyInRight)
	}
}

func TestCompareSecrets_Different(t *testing.T) {
	left := map[string]string{"A": "old"}
	right := map[string]string{"A": "new"}

	r := CompareSecrets(left, right)
	if len(r.Different) != 1 || r.Different[0] != "A" {
		t.Errorf("expected Different=[A], got %v", r.Different)
	}
	if len(r.Identical) != 0 {
		t.Errorf("expected no Identical, got %v", r.Identical)
	}
}

func TestCompareSecrets_Identical(t *testing.T) {
	secrets := map[string]string{"X": "val", "Y": "other"}
	r := CompareSecrets(secrets, secrets)
	if len(r.Identical) != 2 {
		t.Errorf("expected 2 identical, got %d", len(r.Identical))
	}
	if r.HasDifferences() {
		t.Error("expected no differences")
	}
}

func TestCompareSecrets_Summary(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"B": "changed", "C": "3"}

	r := CompareSecrets(left, right)
	summary := r.Summary()
	if !strings.Contains(summary, "Only in left:") {
		t.Error("summary missing 'Only in left:'")
	}
	if !r.HasDifferences() {
		t.Error("expected differences")
	}
}

func TestCompareSecrets_BothEmpty(t *testing.T) {
	r := CompareSecrets(map[string]string{}, map[string]string{})
	if r.HasDifferences() {
		t.Error("empty maps should have no differences")
	}
}
