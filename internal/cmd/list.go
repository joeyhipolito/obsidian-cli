package cmd

import (
	"fmt"

	"github.com/joeyhipolito/obsidian-cli/internal/output"
)

// ListOutput represents the JSON output format for the list command.
type ListOutput struct {
	Directory string   `json:"directory"`
	Notes     []string `json:"notes"`
}

// ListCmd lists notes in the vault or a subdirectory.
func ListCmd(vaultPath, dir string, jsonOutput bool) error {
	// TODO: implement — walk vault directory, list .md files
	if jsonOutput {
		return output.JSON(ListOutput{
			Directory: dir,
			Notes:     []string{},
		})
	}

	fmt.Printf("obsidian list: not yet implemented (dir: %s)\n", dir)
	return nil
}
