package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joeyhipolito/obsidian-cli/internal/config"
)

// ConfigureCmd runs an interactive configuration setup.
// Prompts for Gemini API key and vault path, writes ~/.obsidian/config.
func ConfigureCmd() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Obsidian CLI Configuration")
	fmt.Println("==========================")
	fmt.Println()

	// Check for existing config
	if config.Exists() {
		fmt.Printf("Existing configuration found at %s\n", config.Path())
		fmt.Print("Overwrite? [y/N] ")
		reply, _ := reader.ReadString('\n')
		reply = strings.TrimSpace(reply)
		if !strings.EqualFold(reply, "y") {
			fmt.Println("Configuration cancelled.")
			return nil
		}
		fmt.Println()
	}

	// Load existing config for defaults
	existing, _ := config.Load()

	// Prompt for Gemini API key
	fmt.Println("Get your Gemini API key from:")
	fmt.Println("https://aistudio.google.com/api-keys")
	fmt.Println()
	if existing.GeminiAPIKey != "" {
		masked := existing.GeminiAPIKey[:4] + "..." + existing.GeminiAPIKey[len(existing.GeminiAPIKey)-4:]
		fmt.Printf("Gemini API Key [%s]: ", masked)
	} else {
		fmt.Print("Gemini API Key: ")
	}
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		apiKey = existing.GeminiAPIKey
	}
	if apiKey == "" {
		return fmt.Errorf("Gemini API key is required")
	}

	// Prompt for vault path
	fmt.Println()
	fmt.Println("Path to your Obsidian vault:")
	if existing.VaultPath != "" {
		fmt.Printf("Vault path [%s]: ", existing.VaultPath)
	} else {
		fmt.Print("Vault path: ")
	}
	vaultPath, _ := reader.ReadString('\n')
	vaultPath = strings.TrimSpace(vaultPath)
	if vaultPath == "" {
		vaultPath = existing.VaultPath
	}
	if vaultPath == "" {
		return fmt.Errorf("vault path is required")
	}

	// Expand ~ in vault path
	if strings.HasPrefix(vaultPath, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			vaultPath = home + vaultPath[1:]
		}
	}

	// Validate vault path exists
	info, err := os.Stat(vaultPath)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("vault path does not exist or is not a directory: %s", vaultPath)
	}

	// Save configuration
	cfg := &config.Config{
		GeminiAPIKey: apiKey,
		VaultPath:    vaultPath,
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println()
	fmt.Printf("Configuration saved to %s\n", config.Path())
	fmt.Println()
	fmt.Println("Test your setup:")
	fmt.Println("  obsidian list")
	fmt.Println("  obsidian search \"test\"")
	fmt.Println()
	fmt.Println("Troubleshoot:")
	fmt.Println("  obsidian doctor")

	return nil
}

// ConfigureShowCmd prints the current configuration (with API key masked).
func ConfigureShowCmd(jsonOutput bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !config.Exists() {
		fmt.Println("No configuration file found.")
		fmt.Println("Run 'obsidian configure' to set up.")
		return nil
	}

	// Mask API key for display
	maskedKey := ""
	if cfg.GeminiAPIKey != "" {
		if len(cfg.GeminiAPIKey) > 8 {
			maskedKey = cfg.GeminiAPIKey[:4] + "..." + cfg.GeminiAPIKey[len(cfg.GeminiAPIKey)-4:]
		} else {
			maskedKey = "****"
		}
	}

	if jsonOutput {
		output := map[string]string{
			"config_path":   config.Path(),
			"gemini_apikey": maskedKey,
			"vault_path":    cfg.VaultPath,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(output)
	}

	fmt.Printf("Config file: %s\n", config.Path())
	fmt.Printf("Gemini API key: %s\n", maskedKey)
	fmt.Printf("Vault path: %s\n", cfg.VaultPath)
	return nil
}
