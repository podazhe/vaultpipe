package vault

import (
	"testing"
)

func baseAnnotations() map[string][]Annotation {
	return map[string][]Annotation{
		"DB_PASSWORD": {
			{Key: "sensitivity", Value: "high"},
			{Key: "owner", Value: "dba-team"},
		},
		"API_KEY": {
			{Key: "sensitivity", Value: "high"},
		},
	}
}

func TestAnnotateSecrets_Basic(t *testing.T) {
	secrets := map[string]string{"DB_PASSWORD": "s3cr3t", "API_KEY": "abc123", "HOST": "localhost"}
	result, err := AnnotateSecrets(secrets, baseAnnotations(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if len(result["DB_PASSWORD"].Annotations) != 2 {
		t.Errorf("expected 2 annotations on DB_PASSWORD, got %d", len(result["DB_PASSWORD"].Annotations))
	}
	if len(result["HOST"].Annotations) != 0 {
		t.Errorf("expected no annotations on HOST")
	}
}

func TestAnnotateSecrets_RequirePresent_Errors(t *testing.T) {
	secrets := map[string]string{"HOST": "localhost"}
	_, err := AnnotateSecrets(secrets, baseAnnotations(), true)
	if err == nil {
		t.Fatal("expected error for missing key with requirePresent=true")
	}
}

func TestAnnotateSecrets_NilSecrets_Errors(t *testing.T) {
	_, err := AnnotateSecrets(nil, baseAnnotations(), false)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestGetAnnotation_Found(t *testing.T) {
	secrets := map[string]string{"DB_PASSWORD": "s3cr3t", "API_KEY": "abc123"}
	annotated, _ := AnnotateSecrets(secrets, baseAnnotations(), false)
	v, ok := GetAnnotation(annotated, "DB_PASSWORD", "sensitivity")
	if !ok || v != "high" {
		t.Errorf("expected sensitivity=high, got %q ok=%v", v, ok)
	}
}

func TestGetAnnotation_Missing(t *testing.T) {
	secrets := map[string]string{"DB_PASSWORD": "s3cr3t"}
	annotated, _ := AnnotateSecrets(secrets, baseAnnotations(), false)
	_, ok := GetAnnotation(annotated, "DB_PASSWORD", "nonexistent")
	if ok {
		t.Error("expected false for missing annotation key")
	}
}

func TestFilterByAnnotation_SensitivityHigh(t *testing.T) {
	secrets := map[string]string{"DB_PASSWORD": "s3cr3t", "API_KEY": "abc123", "HOST": "localhost"}
	annotated, _ := AnnotateSecrets(secrets, baseAnnotations(), false)
	filtered := FilterByAnnotation(annotated, "sensitivity", "high")
	if len(filtered) != 2 {
		t.Errorf("expected 2 high-sensitivity secrets, got %d", len(filtered))
	}
	if _, ok := filtered["HOST"]; ok {
		t.Error("HOST should not appear in filtered result")
	}
}

func TestAnnotationSummary_Sorted(t *testing.T) {
	secrets := map[string]string{"DB_PASSWORD": "s3cr3t"}
	annotated, _ := AnnotateSecrets(secrets, baseAnnotations(), false)
	summary := AnnotationSummary(annotated, "DB_PASSWORD")
	if len(summary) != 2 {
		t.Fatalf("expected 2 summary entries, got %d", len(summary))
	}
	// sorted: owner before sensitivity
	if summary[0] != "owner=dba-team" {
		t.Errorf("expected first entry owner=dba-team, got %q", summary[0])
	}
	if summary[1] != "sensitivity=high" {
		t.Errorf("expected second entry sensitivity=high, got %q", summary[1])
	}
}

func TestAnnotationSummary_MissingKey(t *testing.T) {
	annotated := map[string]AnnotatedSecret{}
	if s := AnnotationSummary(annotated, "NOPE"); s != nil {
		t.Errorf("expected nil summary for missing key, got %v", s)
	}
}
