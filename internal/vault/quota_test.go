package vault

import (
	"strings"
	"testing"
)

func TestEnforceQuota_UnderLimits(t *testing.T) {
	secrets := map[string]string{"A": "foo", "B": "bar"}
	result, err := EnforceQuota(secrets, QuotaOptions{MaxKeys: 5, MaxValueLen: 10, MaxTotalSize: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.OK() {
		t.Fatalf("expected no violations, got %v", result.Violations)
	}
	if result.TotalKeys != 2 {
		t.Errorf("expected TotalKeys=2, got %d", result.TotalKeys)
	}
}

func TestEnforceQuota_ExceedsMaxKeys(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	_, err := EnforceQuota(secrets, QuotaOptions{MaxKeys: 2})
	if err == nil {
		t.Fatal("expected error for max_keys violation")
	}
	if !strings.Contains(err.Error(), "max_keys") {
		t.Errorf("error should mention max_keys, got: %v", err)
	}
}

func TestEnforceQuota_ExceedsMaxValueLen(t *testing.T) {
	secrets := map[string]string{"KEY": "this-is-a-very-long-value"}
	_, err := EnforceQuota(secrets, QuotaOptions{MaxValueLen: 5})
	if err == nil {
		t.Fatal("expected error for max_value_len violation")
	}
	if !strings.Contains(err.Error(), "max_value_len") {
		t.Errorf("error should mention max_value_len, got: %v", err)
	}
}

func TestEnforceQuota_ExceedsMaxTotalSize(t *testing.T) {
	secrets := map[string]string{"A": "hello", "B": "world"}
	_, err := EnforceQuota(secrets, QuotaOptions{MaxTotalSize: 5})
	if err == nil {
		t.Fatal("expected error for max_total_size violation")
	}
	if !strings.Contains(err.Error(), "max_total_size") {
		t.Errorf("error should mention max_total_size, got: %v", err)
	}
}

func TestEnforceQuota_DryRunCollectsAll(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	result, err := EnforceQuota(secrets, QuotaOptions{
		MaxKeys:     1,
		MaxTotalSize: 2,
		DryRun:      true,
	})
	if err != nil {
		t.Fatalf("dry-run should not return error, got: %v", err)
	}
	if len(result.Violations) < 2 {
		t.Errorf("expected at least 2 violations in dry-run, got %d", len(result.Violations))
	}
}

func TestEnforceQuota_NoLimitsSet(t *testing.T) {
	secrets := map[string]string{"X": strings.Repeat("z", 10000)}
	result, err := EnforceQuota(secrets, QuotaOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.OK() {
		t.Errorf("expected no violations when no limits set")
	}
}
