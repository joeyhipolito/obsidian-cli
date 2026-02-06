package cmd

import (
	"fmt"
	"time"

	"github.com/joeyhipolito/obsidian-cli/internal/output"
	"github.com/joeyhipolito/obsidian-cli/internal/vault"
)

// CreateOutput represents the JSON output format for the create command.
type CreateOutput struct {
	Path  string `json:"path"`
	Title string `json:"title"`
}

// CreateCmd creates a new note in the vault with optional frontmatter.
func CreateCmd(vaultPath, notePath, title string, jsonOutput bool) error {
	// Build note content with frontmatter
	var content string

	if title != "" {
		fm := map[string]any{
			"title":   title,
			"created": time.Now().Format("2006-01-02"),
		}
		content = vault.FormatFrontmatter(fm) + "\n# " + title + "\n"
	}

	if err := vault.WriteNote(vaultPath, notePath, content); err != nil {
		return err
	}

	if jsonOutput {
		return output.JSON(CreateOutput{
			Path:  notePath,
			Title: title,
		})
	}

	fmt.Printf("Created %s\n", notePath)
	return nil
}
