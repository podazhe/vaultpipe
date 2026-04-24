package vault

import (
	"errors"
	"fmt"
	"sort"
)

// LabeledSecrets holds a map of secrets annotated with arbitrary string labels.
type LabeledSecrets struct {
	secrets map[string]string
	labels  map[string]map[string]string // key -> label -> value
}

// NewLabeledSecrets creates a new LabeledSecrets instance from an existing secret map.
func NewLabeledSecrets(secrets map[string]string) *LabeledSecrets {
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	return &LabeledSecrets{
		secrets: copy,
		labels:  make(map[string]map[string]string),
	}
}

// AddLabel attaches a label key/value pair to a secret key.
// Returns an error if the secret key does not exist.
func (ls *LabeledSecrets) AddLabel(secretKey, labelKey, labelValue string) error {
	if secretKey == "" {
		return errors.New("secret key must not be empty")
	}
	if labelKey == "" {
		return errors.New("label key must not be empty")
	}
	if _, ok := ls.secrets[secretKey]; !ok {
		return fmt.Errorf("secret key %q not found", secretKey)
	}
	if ls.labels[secretKey] == nil {
		ls.labels[secretKey] = make(map[string]string)
	}
	ls.labels[secretKey][labelKey] = labelValue
	return nil
}

// GetLabels returns the labels associated with a secret key.
func (ls *LabeledSecrets) GetLabels(secretKey string) map[string]string {
	result := make(map[string]string)
	for k, v := range ls.labels[secretKey] {
		result[k] = v
	}
	return result
}

// FilterByLabel returns a map of secrets whose labels contain the given key/value pair.
func (ls *LabeledSecrets) FilterByLabel(labelKey, labelValue string) map[string]string {
	result := make(map[string]string)
	for secretKey, secretVal := range ls.secrets {
		if lbls, ok := ls.labels[secretKey]; ok {
			if v, found := lbls[labelKey]; found && v == labelValue {
				result[secretKey] = secretVal
			}
		}
	}
	return result
}

// ListLabeled returns all secret keys that have at least one label, sorted.
func (ls *LabeledSecrets) ListLabeled() []string {
	keys := make([]string, 0, len(ls.labels))
	for k := range ls.labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// RemoveLabel removes a specific label from a secret key.
// Returns an error if the label does not exist.
func (ls *LabeledSecrets) RemoveLabel(secretKey, labelKey string) error {
	lbls, ok := ls.labels[secretKey]
	if !ok {
		return fmt.Errorf("no labels found for secret key %q", secretKey)
	}
	if _, exists := lbls[labelKey]; !exists {
		return fmt.Errorf("label %q not found on secret key %q", labelKey, secretKey)
	}
	delete(lbls, labelKey)
	if len(lbls) == 0 {
		delete(ls.labels, secretKey)
	}
	return nil
}
