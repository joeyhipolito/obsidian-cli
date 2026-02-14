// Package config handles reading and writing the Obsidian CLI configuration file.
// Configuration is stored in ~/.obsidian/config in INI-style format.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ConfigDir is the directory name for Obsidian CLI configuration.
	ConfigDir = ".obsidian"
	// ConfigFile is the configuration file name.
	ConfigFile = "config"
)

// Config represents the Obsidian CLI configuration.
type Config struct {
	GeminiAPIKey string
	VaultPath    string
	WebsitePath  string
}

// Path returns the full path to the config file (~/.obsidian/config).
func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ConfigDir, ConfigFile)
}

// Dir returns the full path to the config directory (~/.obsidian/).
func Dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ConfigDir)
}

// Load reads the configuration from ~/.obsidian/config.
// Returns an empty Config (not an error) if the file doesn't exist.
func Load() (*Config, error) {
	cfg := &Config{}
	path := Path()
	if path == "" {
		return cfg, nil
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "gemini_apikey":
			cfg.GeminiAPIKey = value
		case "vault_path":
			cfg.VaultPath = value
		case "website_path":
			cfg.WebsitePath = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return cfg, nil
}

// Save writes the configuration to ~/.obsidian/config with proper permissions.
func Save(cfg *Config) error {
	dir := Dir()
	if dir == "" {
		return fmt.Errorf("cannot determine home directory")
	}

	// Create config directory with 700 permissions
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path := Path()

	// Build config content
	var b strings.Builder
	b.WriteString("# Obsidian CLI Configuration\n")
	b.WriteString("# Created by: obsidian configure\n")
	b.WriteString("\n")
	b.WriteString("# Your Gemini API Key\n")
	b.WriteString("# Get from: https://aistudio.google.com/api-keys\n")
	fmt.Fprintf(&b, "gemini_apikey=%s\n", cfg.GeminiAPIKey)
	b.WriteString("\n")
	b.WriteString("# Path to your Obsidian vault\n")
	fmt.Fprintf(&b, "vault_path=%s\n", cfg.VaultPath)
	if cfg.WebsitePath != "" {
		b.WriteString("\n")
		b.WriteString("# Path to your website project (for obsidian sync)\n")
		fmt.Fprintf(&b, "website_path=%s\n", cfg.WebsitePath)
	}

	// Write file with 600 permissions
	if err := os.WriteFile(path, []byte(b.String()), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Exists returns true if the config file exists.
func Exists() bool {
	path := Path()
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

// Permissions returns the file permissions of the config file, or an error.
func Permissions() (os.FileMode, error) {
	path := Path()
	if path == "" {
		return 0, fmt.Errorf("cannot determine config path")
	}
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}

// ResolveAPIKey returns the Gemini API key using config priority:
// config file > environment variable.
func ResolveAPIKey() string {
	cfg, err := Load()
	if err == nil && cfg.GeminiAPIKey != "" {
		return cfg.GeminiAPIKey
	}
	return os.Getenv("GEMINI_API_KEY")
}

// ResolveVaultPath returns the vault path from config or environment.
func ResolveVaultPath() string {
	cfg, err := Load()
	if err == nil && cfg.VaultPath != "" {
		return cfg.VaultPath
	}
	return os.Getenv("OBSIDIAN_VAULT_PATH")
}

// ResolveWebsitePath returns the website path from config or environment.
func ResolveWebsitePath() string {
	cfg, err := Load()
	if err == nil && cfg.WebsitePath != "" {
		return cfg.WebsitePath
	}
	return os.Getenv("OBSIDIAN_WEBSITE_PATH")
}
