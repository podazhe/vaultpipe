package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ArchiveEntry holds a timestamped snapshot of secrets at a given path.
type ArchiveEntry struct {
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
	ArchivedAt time.Time        `json:"archived_at"`
	Note      string            `json:"note,omitempty"`
}

// ArchiveManager stores and retrieves archived secret snapshots on disk.
type ArchiveManager struct {
	dir string
}

// NewArchiveManager creates an ArchiveManager rooted at dir.
func NewArchiveManager(dir string) (*ArchiveManager, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("archive: create directory: %w", err)
	}
	return &ArchiveManager{dir: dir}, nil
}

// Save writes secrets to a new archive entry file named by path and timestamp.
func (a *ArchiveManager) Save(path string, secrets map[string]string, note string) (*ArchiveEntry, error) {
	if path == "" {
		return nil, fmt.Errorf("archive: path must not be empty")
	}
	entry := &ArchiveEntry{
		Path:       path,
		Secrets:    copyMap(secrets),
		ArchivedAt: time.Now().UTC(),
		Note:       note,
	}
	fileName := a.fileName(path, entry.ArchivedAt)
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("archive: marshal: %w", err)
	}
	if err := os.WriteFile(fileName, data, 0600); err != nil {
		return nil, fmt.Errorf("archive: write file: %w", err)
	}
	return entry, nil
}

// Load reads an archive entry from the given file path.
func (a *ArchiveManager) Load(filePath string) (*ArchiveEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("archive: file not found: %s", filePath)
		}
		return nil, fmt.Errorf("archive: read file: %w", err)
	}
	var entry ArchiveEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("archive: unmarshal: %w", err)
	}
	return &entry, nil
}

func (a *ArchiveManager) fileName(path string, t time.Time) string {
	safe := sanitiseKey(path)
	return fmt.Sprintf("%s/%s_%d.json", a.dir, safe, t.UnixNano())
}
