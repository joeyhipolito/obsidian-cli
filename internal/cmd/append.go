package cmd

import (
	"fmt"

	"github.com/joeyhipolito/obsidian-cli/internal/output"
)

// AppendOutput represents the JSON output format for the append command.
type AppendOutput struct {
	Path    string `json:"path"`
	Appended string `json:"appended"`
}

// AppendCmd appends text to an existing note.
func AppendCmd(vaultPath, notePath, text string, jsonOutput bool) error {
	// TODO: implement — append text to note file, read from stdin if text is empty
	if jsonOutput {
		return output.JSON(AppendOutput{
			Path:    notePath,
			Appended: text,
		})
	}

	fmt.Printf("obsidian append: not yet implemented (path: %s)\n", notePath)
	return nil
}
