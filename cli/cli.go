package cli

import (
	"flag"
	"fmt"
	"os"
)

// CLI represents the command-line interface
type CLI struct {
	args []string
}

// NewCLI creates a new CLI instance
func NewCLI(args []string) *CLI {
	return &CLI{args: args}
}

// Run executes the CLI
func (c *CLI) Run() error {
	if len(c.args) == 0 {
		return c.showHelp()
	}

	command := c.args[0]

	switch command {
	case "edit":
		return c.editCommand()
	case "export":
		return c.exportCommand()
	case "import":
		return c.importCommand()
	case "batch":
		return c.batchCommand()
	case "script":
		return c.scriptCommand()
	case "validate":
		return c.validateCommand()
	case "backup":
		return c.backupCommand()
	case "combat-pack":
		return c.combatPackCommand()
	case "help", "-h", "--help":
		return c.showHelp()
	case "version", "-v", "--version":
		return c.showVersion()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		return c.showHelp()
	}
}

// editCommand edits a save file
func (c *CLI) editCommand() error {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)
	file := fs.String("file", "", "Save file path (required)")
	charID := fs.Int("char", -1, "Character ID (0-15)")
	level := fs.Int("level", -1, "Set character level")
	hp := fs.Int("hp", -1, "Set character HP")
	mp := fs.Int("mp", -1, "Set character MP")
	output := fs.String("output", "", "Output file path (defaults to input)")

	if err := fs.Parse(c.args[1:]); err != nil {
		return err
	}

	if *file == "" {
		return fmt.Errorf("--file is required")
	}

	return c.handleEditCommand(*file, *charID, *level, *hp, *mp, *output)
}

// exportCommand exports save data to JSON
func (c *CLI) exportCommand() error {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	file := fs.String("file", "", "Save file path (required)")
	output := fs.String("output", "", "Output JSON file (required)")
	format := fs.String("format", "full", "Export format: full, characters, inventory, party, magic, espers")

	if err := fs.Parse(c.args[1:]); err != nil {
		return err
	}

	if *file == "" || *output == "" {
		return fmt.Errorf("--file and --output are required")
	}

	return c.handleExportCommand(*file, *output, *format)
}

// importCommand imports JSON data into save file
func (c *CLI) importCommand() error {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	file := fs.String("file", "", "Save file path (required)")
	input := fs.String("input", "", "Input JSON file (required)")
	format := fs.String("format", "full", "Import format: full, characters, inventory, party, magic, espers")
	backup := fs.Bool("backup", true, "Create backup before import")

	if err := fs.Parse(c.args[1:]); err != nil {
		return err
	}

	if *file == "" || *input == "" {
		return fmt.Errorf("--file and --input are required")
	}

	return c.handleImportCommand(*file, *input, *format, *backup)
}

// batchCommand performs batch operations
func (c *CLI) batchCommand() error {
	fs := flag.NewFlagSet("batch", flag.ExitOnError)
	file := fs.String("file", "", "Save file path (required)")
	operation := fs.String("op", "", "Operation: max-stats, max-items, max-magic, max-all")
	output := fs.String("output", "", "Output file path (defaults to input)")

	if err := fs.Parse(c.args[1:]); err != nil {
		return err
	}

	if *file == "" || *operation == "" {
		return fmt.Errorf("--file and --op are required")
	}

	return c.handleBatchCommand(*file, *operation, *output)
}

// scriptCommand runs a Lua script
func (c *CLI) scriptCommand() error {
	fs := flag.NewFlagSet("script", flag.ExitOnError)
	file := fs.String("file", "", "Save file path (required)")
	script := fs.String("script", "", "Lua script file (required)")
	output := fs.String("output", "", "Output file path (defaults to input)")

	if err := fs.Parse(c.args[1:]); err != nil {
		return err
	}

	if *file == "" || *script == "" {
		return fmt.Errorf("--file and --script are required")
	}

	return c.handleScriptCommand(*file, *script, *output)
}

// validateCommand validates a save file
func (c *CLI) validateCommand() error {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	file := fs.String("file", "", "Save file path (required)")
	fix := fs.Bool("fix", false, "Attempt to fix issues automatically")

	if err := fs.Parse(c.args[1:]); err != nil {
		return err
	}

	if *file == "" {
		return fmt.Errorf("--file is required")
	}

	return c.handleValidateCommand(*file, *fix)
}

// backupCommand creates a backup of a save file
func (c *CLI) backupCommand() error {
	fs := flag.NewFlagSet("backup", flag.ExitOnError)
	file := fs.String("file", "", "Save file path (required)")
	output := fs.String("output", "", "Backup output path (optional)")

	if err := fs.Parse(c.args[1:]); err != nil {
		return err
	}

	if *file == "" {
		return fmt.Errorf("--file is required")
	}

	return c.handleBackupCommand(*file, *output)
}

// showHelp displays CLI help
func (c *CLI) showHelp() error {
	help := `
Final Fantasy VI Save Editor - CLI

USAGE:
    ffvi_editor [command] [options]

COMMANDS:
    edit       Edit a save file directly
    export     Export save data to JSON
    import     Import JSON data into save file
    batch      Perform batch operations
	script     Run a Lua script on a save file
	validate   Validate save file integrity
	backup     Create a backup of a save file
	combat-pack Run Combat Depth Pack helpers (Encounter/Boss/Companion/Smoke)
    help       Show this help message
    version    Show version information

EXAMPLES:
    # Set character 0 to level 99
    ffvi_editor edit --file save.json --char 0 --level 99

    # Export full save to JSON
    ffvi_editor export --file save.json --output export.json --format full

    # Import characters from JSON
    ffvi_editor import --file save.json --input characters.json --format characters

    # Maximize all stats
    ffvi_editor batch --file save.json --op max-all

    # Run custom script
    ffvi_editor script --file save.json --script custom.lua

    # Validate save file
    ffvi_editor validate --file save.json --fix

For more information, visit: https://github.com/username/ffvi-save-editor
`
	fmt.Println(help)
	return nil
}

// showVersion displays version information
func (c *CLI) showVersion() error {
	fmt.Println("Final Fantasy VI Save Editor v4.0.0")
	fmt.Println("Build: Phase 4 - Advanced Integration")
	return nil
}
