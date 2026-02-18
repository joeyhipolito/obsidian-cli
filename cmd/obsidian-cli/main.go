// Package main implements the obsidian binary.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeyhipolito/obsidian-cli/internal/cmd"
	"github.com/joeyhipolito/obsidian-cli/internal/config"
)

const version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse command line arguments
	args := os.Args[1:]

	// Handle help and version flags
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printUsage()
		return nil
	}

	if args[0] == "--version" || args[0] == "-v" {
		fmt.Printf("obsidian version %s\n", version)
		return nil
	}

	// Parse subcommand
	subcommand := args[0]
	remainingArgs := args[1:]

	// Extract --json from anywhere in args
	jsonOutput := false
	var filteredArgs []string
	for _, arg := range remainingArgs {
		if arg == "--json" {
			jsonOutput = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	// Extract additional flags
	dryRun := false
	forceFlag := false
	applyFlag := false
	fixFlag := false
	staleDays := 30
	ingestSource := ""
	ingestTopic := ""
	ingestDomain := ""
	ingestSince := ""
	var cleanedArgs []string
	for i := 0; i < len(filteredArgs); i++ {
		switch filteredArgs[i] {
		case "--dry-run":
			dryRun = true
		case "--force":
			forceFlag = true
		case "--apply":
			applyFlag = true
		case "--fix":
			fixFlag = true
		case "--stale-days":
			if i+1 < len(filteredArgs) {
				if n, err := parseInt(filteredArgs[i+1]); err == nil {
					staleDays = n
				}
				i++
			}
		case "--source":
			if i+1 < len(filteredArgs) {
				ingestSource = filteredArgs[i+1]
				i++
			}
		case "--topic":
			if i+1 < len(filteredArgs) {
				ingestTopic = filteredArgs[i+1]
				i++
			}
		case "--domain":
			if i+1 < len(filteredArgs) {
				ingestDomain = filteredArgs[i+1]
				i++
			}
		case "--since":
			if i+1 < len(filteredArgs) {
				ingestSince = filteredArgs[i+1]
				i++
			}
		default:
			cleanedArgs = append(cleanedArgs, filteredArgs[i])
		}
	}
	filteredArgs = cleanedArgs

	// Commands that don't require config
	switch subcommand {
	case "configure":
		if len(filteredArgs) > 0 && filteredArgs[0] == "show" {
			return cmd.ConfigureShowCmd(jsonOutput)
		}
		return cmd.ConfigureCmd()
	case "doctor":
		return cmd.DoctorCmd(jsonOutput)
	case "read", "append", "create", "list", "search", "index", "sync", "enrich", "maintain", "ingest":
		// handled below after vault resolution
	default:
		return fmt.Errorf("unknown command: %s\n\nRun 'obsidian --help' for usage", subcommand)
	}

	// Resolve vault path
	vaultPath := config.ResolveVaultPath()
	if vaultPath == "" {
		return fmt.Errorf("no vault path configured\n\nRun 'obsidian configure' to set up, or set OBSIDIAN_VAULT_PATH")
	}

	// Dispatch to appropriate command handler
	switch subcommand {
	case "read":
		if len(filteredArgs) < 1 {
			return fmt.Errorf("read requires a note path\n\nUsage: obsidian read <path>")
		}
		return cmd.ReadCmd(vaultPath, filteredArgs[0], jsonOutput)

	case "append":
		return handleAppendCommand(vaultPath, filteredArgs, jsonOutput)

	case "create":
		return handleCreateCommand(vaultPath, filteredArgs, jsonOutput)

	case "list":
		dir := ""
		if len(filteredArgs) > 0 {
			dir = filteredArgs[0]
		}
		return cmd.ListCmd(vaultPath, dir, jsonOutput)

	case "search":
		return handleSearchCommand(vaultPath, filteredArgs, jsonOutput)

	case "index":
		return cmd.IndexCmd(vaultPath, jsonOutput)

	case "sync":
		websitePath := config.ResolveWebsitePath()
		if websitePath == "" {
			return fmt.Errorf("no website path configured\n\nSet website_path in ~/.obsidian/config or OBSIDIAN_WEBSITE_PATH env var")
		}
		return cmd.SyncCmd(vaultPath, websitePath, dryRun, forceFlag, jsonOutput)

	case "enrich":
		return cmd.EnrichCmd(vaultPath, applyFlag, jsonOutput)

	case "maintain":
		return cmd.MaintainCmd(vaultPath, staleDays, fixFlag, jsonOutput)

	case "ingest":
		return cmd.IngestCmd(vaultPath, cmd.IngestOptions{
			Source:     ingestSource,
			Topic:      ingestTopic,
			Domain:     ingestDomain,
			Since:      ingestSince,
			DryRun:     dryRun,
			JSONOutput: jsonOutput,
		})
	}

	return nil
}

// parseInt parses a string to int, returning an error if invalid.
func parseInt(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number: %s", s)
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

// handleAppendCommand parses and executes the append command.
func handleAppendCommand(vaultPath string, args []string, jsonOutput bool) error {
	if len(args) < 1 {
		return fmt.Errorf("append requires a note path\n\nUsage: obsidian append <path> <text>\n       echo 'text' | obsidian append <path>")
	}
	notePath := args[0]
	text := ""
	if len(args) > 1 {
		text = strings.Join(args[1:], " ")
	}
	return cmd.AppendCmd(vaultPath, notePath, text, jsonOutput)
}

// handleCreateCommand parses and executes the create command.
func handleCreateCommand(vaultPath string, args []string, jsonOutput bool) error {
	if len(args) < 1 {
		return fmt.Errorf("create requires a note path\n\nUsage: obsidian create <path> [--title <title>]")
	}
	notePath := args[0]
	title := ""
	remaining := args[1:]

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--title":
			if i+1 >= len(remaining) {
				return fmt.Errorf("--title requires an argument")
			}
			title = remaining[i+1]
			i++
		default:
			return fmt.Errorf("unknown flag: %s", remaining[i])
		}
	}

	return cmd.CreateCmd(vaultPath, notePath, title, jsonOutput)
}

// handleSearchCommand parses and executes the search command.
func handleSearchCommand(vaultPath string, args []string, jsonOutput bool) error {
	mode := ""
	var queryParts []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--mode":
			if i+1 >= len(args) {
				return fmt.Errorf("--mode requires an argument (keyword, semantic, or hybrid)")
			}
			mode = args[i+1]
			i++
		default:
			queryParts = append(queryParts, args[i])
		}
	}

	if len(queryParts) == 0 {
		return fmt.Errorf("search requires a query\n\nUsage: obsidian search <query> [--mode keyword|semantic|hybrid]")
	}

	query := strings.Join(queryParts, " ")
	return cmd.SearchCmd(vaultPath, query, mode, jsonOutput)
}

func printUsage() {
	fmt.Printf(`obsidian - Obsidian vault CLI tool (v%s)

USAGE:
    obsidian <command> [options]

COMMANDS:
    read <path>             Read a note's content
    append <path> <text>    Append text to a note
    create <path>           Create a new note
    list [dir]              List notes in vault or directory
    search <query>          Search notes (keyword + semantic)
                            --mode keyword|semantic|hybrid (default: hybrid)
    index                   Build/update the search index
    sync                    Sync website content metadata into vault
                            --dry-run  Preview without writing
                            --force    Overwrite unchanged + include unpublished
    enrich                  Suggest links, tags, detect orphan notes
                            --apply    Write suggested links to notes
    maintain                Vault health checks and reporting
                            --stale-days N  Days before note is stale (default: 30)
                            --fix           Add frontmatter to notes missing it
    ingest                  Import data from external sources into vault
                            --source scout|learnings  (required)
                            --topic <name>            Filter scout by topic
                            --domain <name>           Filter learnings by domain
                            --since <duration>        e.g. 7d, 24h, 2w
                            --dry-run                 Preview without writing
    configure               Set up API key and vault path
    configure show          Show current configuration
    doctor                  Validate installation and configuration

GLOBAL OPTIONS:
    --json              Output in JSON format
    --help, -h          Show this help
    --version, -v       Show version

CONFIGURATION:
    obsidian configure              Interactive setup
    obsidian configure show         Show current config (key masked)
    obsidian doctor                 Validate setup and troubleshoot
    Config file: ~/.obsidian/config

EXAMPLES:
    obsidian configure                              # First-time setup
    obsidian read daily/2026-02-07.md               # Read a note
    obsidian append daily/2026-02-07.md "New task"  # Append to note
    obsidian create projects/new-idea.md            # Create a note
    obsidian list daily/                            # List notes in folder
    obsidian search "project ideas"                 # Hybrid search (default)
    obsidian search "golang" --mode keyword         # Keyword-only search
    obsidian index                                  # Build search index
    obsidian sync                                   # Sync website to vault
    obsidian sync --dry-run                         # Preview sync changes
    obsidian enrich                                 # Find note connections
    obsidian enrich --apply                         # Apply suggested links
    obsidian maintain                               # Vault health report
    obsidian ingest --source scout                  # Import scout intel
    obsidian ingest --source scout --topic "ai-models" --since 7d
    obsidian ingest --source learnings              # Import orchestrator learnings
    obsidian ingest --source learnings --domain dev --since 30d
    obsidian ingest --source scout --dry-run        # Preview what would be created
    obsidian doctor                                 # Check setup

For more information, visit: https://obsidian.md
`, version)
}
