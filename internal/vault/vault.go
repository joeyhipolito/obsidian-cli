package vault

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// NoteInfo contains metadata about a note file.
type NoteInfo struct {
	Path    string `json:"path"`    // Relative path within vault
	Name    string `json:"name"`    // Filename without extension
	ModTime int64  `json:"mod_time"` // Unix timestamp of last modification
	Size    int64  `json:"size"`    // File size in bytes
}

// ReadNote reads and parses a note from the vault.
// notePath is relative to vaultPath (e.g., "daily/2026-02-07.md").
func ReadNote(vaultPath, notePath string) (*Note, error) {
	fullPath := resolvePath(vaultPath, notePath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read note: %w", err)
	}

	note := ParseNote(string(data))
	return note, nil
}

// WriteNote creates a new note file. Returns an error if the file already exists.
func WriteNote(vaultPath, notePath, content string) error {
	fullPath := resolvePath(vaultPath, notePath)

	// Check if file already exists
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("note already exists: %s", notePath)
	}

	// Create parent directories
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("cannot write note: %w", err)
	}

	return nil
}

// AppendToNote appends text to an existing note.
func AppendToNote(vaultPath, notePath, text string) error {
	fullPath := resolvePath(vaultPath, notePath)

	// Verify file exists
	if _, err := os.Stat(fullPath); err != nil {
		return fmt.Errorf("note not found: %s", notePath)
	}

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open note: %w", err)
	}
	defer f.Close()

	// Ensure text starts on a new line
	info, _ := f.Stat()
	if info.Size() > 0 {
		// Read last byte to check if file ends with newline
		existing, _ := os.ReadFile(fullPath)
		if len(existing) > 0 && existing[len(existing)-1] != '\n' {
			text = "\n" + text
		}
	}

	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	if _, err := f.WriteString(text); err != nil {
		return fmt.Errorf("cannot append to note: %w", err)
	}

	return nil
}

// ListNotes lists all .md files in a vault directory.
// dir is relative to vaultPath; empty string lists the entire vault.
func ListNotes(vaultPath, dir string) ([]NoteInfo, error) {
	searchPath := vaultPath
	if dir != "" {
		searchPath = filepath.Join(vaultPath, dir)
	}

	// Verify directory exists
	info, err := os.Stat(searchPath)
	if err != nil {
		return nil, fmt.Errorf("directory not found: %s", dir)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", dir)
	}

	var notes []NoteInfo
	err = filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible entries
		}

		// Skip hidden directories (like .obsidian, .git)
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		// Only include .md files
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		relPath, _ := filepath.Rel(vaultPath, path)

		info, err := d.Info()
		if err != nil {
			return nil // skip files we can't stat
		}

		name := strings.TrimSuffix(d.Name(), ".md")
		notes = append(notes, NoteInfo{
			Path:    relPath,
			Name:    name,
			ModTime: info.ModTime().Unix(),
			Size:    info.Size(),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("cannot list notes: %w", err)
	}

	return notes, nil
}

// resolvePath joins vault path with note path, adding .md extension if needed.
func resolvePath(vaultPath, notePath string) string {
	if !strings.HasSuffix(notePath, ".md") {
		notePath += ".md"
	}
	return filepath.Join(vaultPath, notePath)
}
