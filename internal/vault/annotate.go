package vault

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Annotation holds metadata attached to a secret key.
type Annotation struct {
	Key   string
	Value string
}

// AnnotatedSecret pairs a secret value with its annotations.
type AnnotatedSecret struct {
	Value       string
	Annotations []Annotation
}

// AnnotateSecrets attaches annotations to secrets matching the given keys.
// annotations is a map of secretKey -> []Annotation.
// Keys not present in secrets are ignored unless requirePresent is true.
func AnnotateSecrets(
	secrets map[string]string,
	annotations map[string][]Annotation,
	requirePresent bool,
) (map[string]AnnotatedSecret, error) {
	if secrets == nil {
		return nil, errors.New("annotate: secrets map must not be nil")
	}

	result := make(map[string]AnnotatedSecret, len(secrets))
	for k, v := range secrets {
		result[k] = AnnotatedSecret{Value: v}
	}

	for key, anns := range annotations {
		if _, ok := secrets[key]; !ok {
			if requirePresent {
				return nil, fmt.Errorf("annotate: key %q not found in secrets", key)
			}
			continue
		}
		entry := result[key]
		entry.Annotations = append(entry.Annotations, anns...)
		result[key] = entry
	}

	return result, nil
}

// GetAnnotation returns the value of a named annotation for a given secret key.
// Returns an empty string and false if not found.
func GetAnnotation(annotated map[string]AnnotatedSecret, secretKey, annotationKey string) (string, bool) {
	entry, ok := annotated[secretKey]
	if !ok {
		return "", false
	}
	for _, a := range entry.Annotations {
		if a.Key == annotationKey {
			return a.Value, true
		}
	}
	return "", false
}

// FilterByAnnotation returns only those secrets that carry an annotation
// whose key equals annotationKey and whose value contains matchValue.
func FilterByAnnotation(annotated map[string]AnnotatedSecret, annotationKey, matchValue string) map[string]AnnotatedSecret {
	out := make(map[string]AnnotatedSecret)
	for k, entry := range annotated {
		for _, a := range entry.Annotations {
			if a.Key == annotationKey && strings.Contains(a.Value, matchValue) {
				out[k] = entry
				break
			}
		}
	}
	return out
}

// AnnotationSummary returns a sorted slice of "key=value" strings for all
// annotations on a given secret key.
func AnnotationSummary(annotated map[string]AnnotatedSecret, secretKey string) []string {
	entry, ok := annotated[secretKey]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(entry.Annotations))
	for _, a := range entry.Annotations {
		out = append(out, fmt.Sprintf("%s=%s", a.Key, a.Value))
	}
	sort.Strings(out)
	return out
}
