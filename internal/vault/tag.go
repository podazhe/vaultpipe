package vault

import (
	"fmt"
	"sort"
	"strings"
)

// TaggedSecrets associates a set of string tags with each secret key.
type TaggedSecrets struct {
	Secrets map[string]string
	Tags    map[string][]string
}

// NewTaggedSecrets wraps a secrets map with an empty tag index.
func NewTaggedSecrets(secrets map[string]string) *TaggedSecrets {
	tags := make(map[string][]string, len(secrets))
	for k := range secrets {
		tags[k] = []string{}
	}
	return &TaggedSecrets{Secrets: secrets, Tags: tags}
}

// Tag adds one or more tags to the given secret key.
// Returns an error if the key does not exist in the secrets map.
func (t *TaggedSecrets) Tag(key string, tags ...string) error {
	if _, ok := t.Secrets[key]; !ok {
		return fmt.Errorf("tag: key %q not found in secrets", key)
	}
	existing := t.Tags[key]
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if !containsStr(existing, tag) {
			existing = append(existing, tag)
		}
	}
	sort.Strings(existing)
	t.Tags[key] = existing
	return nil
}

// FilterByTag returns a map of secrets whose keys carry all of the given tags.
func (t *TaggedSecrets) FilterByTag(tags ...string) map[string]string {
	out := make(map[string]string)
	for k, v := range t.Secrets {
		if hasAllTags(t.Tags[k], tags) {
			out[k] = v
		}
	}
	return out
}

// ListTags returns a deduplicated, sorted list of all tags in use.
func (t *TaggedSecrets) ListTags() []string {
	seen := map[string]struct{}{}
	for _, ts := range t.Tags {
		for _, tag := range ts {
			seen[tag] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for tag := range seen {
		out = append(out, tag)
	}
	sort.Strings(out)
	return out
}

func hasAllTags(have, want []string) bool {
	for _, w := range want {
		if !containsStr(have, w) {
			return false
		}
	}
	return true
}
