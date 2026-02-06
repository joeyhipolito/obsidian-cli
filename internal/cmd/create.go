package cmd

import (
	"fmt"

	"github.com/joeyhipolito/obsidian-cli/internal/output"
)

// CreateOutput represents the JSON output format for the create command.
type CreateOutput struct {
	Path  string `json:"path"`
	Title string `json:"title"`
}

// CreateCmd creates a new note in the vault.
func CreateCmd(vaultPath, notePath, title string, jsonOutput bool) error {
	// TODO: implement — create note file with optional frontmatter/title
	if jsonOutput {
		return output.JSON(CreateOutput{
			Path:  notePath,
			Title: title,
		})
	}

	fmt.Printf("obsidian create: not yet implemented (path: %s)\n", notePath)
	return nil
}
