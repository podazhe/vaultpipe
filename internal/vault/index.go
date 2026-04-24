package vault

import (
	"fmt"
	"sort"
	"strings"
)

// IndexEntry holds metadata about a secret key in the index.
type IndexEntry struct {
	Path   string
	Key    string
	Tags   []string
	Labels map[string]string
}

// SecretIndex provides fast lookup and listing of secret keys across paths.
type SecretIndex struct {
	entries map[string]*IndexEntry // composite key: path+":"+key
}

// NewSecretIndex creates an empty SecretIndex.
func NewSecretIndex() *SecretIndex {
	return &SecretIndex{
		entries: make(map[string]*IndexEntry),
	}
}

// Add registers a secret key under a path with optional tags and labels.
func (idx *SecretIndex) Add(path, key string, tags []string, labels map[string]string) {
	composite := compositeKey(path, key)
	idx.entries[composite] = &IndexEntry{
		Path:   path,
		Key:    key,
		Tags:   append([]string{}, tags...),
		Labels: copyLabels(labels),
	}
}

// Remove deletes a key from the index.
func (idx *SecretIndex) Remove(path, key string) {
	delete(idx.entries, compositeKey(path, key))
}

// Get retrieves an entry by path and key.
func (idx *SecretIndex) Get(path, key string) (*IndexEntry, bool) {
	e, ok := idx.entries[compositeKey(path, key)]
	return e, ok
}

// ListByPath returns all entries under a given path, sorted by key.
func (idx *SecretIndex) ListByPath(path string) []*IndexEntry {
	var out []*IndexEntry
	for _, e := range idx.entries {
		if e.Path == path {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// Search returns entries whose key contains the given substring.
func (idx *SecretIndex) Search(substr string) []*IndexEntry {
	var out []*IndexEntry
	for _, e := range idx.entries {
		if strings.Contains(e.Key, substr) {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Key < out[j].Key
	})
	return out
}

// Size returns the total number of indexed entries.
func (idx *SecretIndex) Size() int { return len(idx.entries) }

func compositeKey(path, key string) string {
	return fmt.Sprintf("%s:%s", path, key)
}

func copyLabels(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
