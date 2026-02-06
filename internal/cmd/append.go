package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joeyhipolito/obsidian-cli/internal/output"
	"github.com/joeyhipolito/obsidian-cli/internal/vault"
)

// AppendOutput represents the JSON output format for the append command.
type AppendOutput struct {
	Path     string `json:"path"`
	Appended string `json:"appended"`
}

// AppendCmd appends text to an existing note.
// If text is empty, reads from stdin.
func AppendCmd(vaultPath, notePath, text string, jsonOutput bool) error {
	if text == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
		text = strings.TrimRight(string(data), "\n")
		if text == "" {
			return fmt.Errorf("no text provided\n\nUsage: obsidian append <path> <text>\n       echo 'text' | obsidian append <path>")
		}
	}

	if err := vault.AppendToNote(vaultPath, notePath, text); err != nil {
		return err
	}

	if jsonOutput {
		return output.JSON(AppendOutput{
			Path:     notePath,
			Appended: text,
		})
	}

	fmt.Printf("Appended to %s\n", notePath)
	return nil
}
