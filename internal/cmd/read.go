package cmd

import (
	"fmt"

	"github.com/joeyhipolito/obsidian-cli/internal/output"
)

// ReadOutput represents the JSON output format for the read command.
type ReadOutput struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// ReadCmd reads a note's content from the vault.
func ReadCmd(vaultPath, notePath string, jsonOutput bool) error {
	// TODO: implement — read note file, parse frontmatter, return content
	if jsonOutput {
		return output.JSON(ReadOutput{
			Path:    notePath,
			Content: "",
		})
	}

	fmt.Printf("obsidian read: not yet implemented (path: %s)\n", notePath)
	return nil
}
