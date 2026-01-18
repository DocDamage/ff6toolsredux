package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ffvi_editor/io/pr"
	"ffvi_editor/scripting"
)

// LoadSaveFile loads a save file from the specified path
func (c *CLI) LoadSaveFile(filepath string) (*pr.PR, error) {
	p := pr.New()
	if err := p.Load(filepath, 0); err != nil {
		return nil, fmt.Errorf("failed to load save file: %w", err)
	}
	return p, nil
}

// SaveSaveFile saves a save file to the specified path
func (c *CLI) SaveSaveFile(save *pr.PR, filepath string) error {
	if err := save.Save(0, filepath, 0); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	fmt.Printf("Successfully saved to: %s\n", filepath)
	return nil
}

// handleEdit loads a save, modifies it, and saves it
// TODO: Implement character editing in CLI (Phase 4)
// Placeholder for character modification commands (stats, equipment, abilities)
// Will integrate with character models from io/models/pr/characters.go
func (c *CLI) handleEditCommand(file string, charID, level, hp, mp int, output string) error {
	return fmt.Errorf("CLI edit command not yet implemented (Phase 4)")
}

// handleExport exports save data to JSON in various formats
// TODO: Implement export command in CLI (Phase 4)
// Placeholder for save file export to JSON, CSV, or other formats
// Will use existing export functions from UI dashboards
func (c *CLI) handleExportCommand(file, output, format string) error {
	return fmt.Errorf("CLI export command not yet implemented (Phase 4)")
}

// handleImport imports JSON data into a save file
// TODO: Implement import command in CLI (Phase 4)
// Placeholder for importing save files from various formats
// Will integrate with existing save loading infrastructure
func (c *CLI) handleImportCommand(file, input, format string, backup bool) error {
	return fmt.Errorf("CLI import command not yet implemented (Phase 4)")
}

// handleBatch performs batch operations on a save file
// TODO: Implement batch command in CLI (Phase 4)
// Placeholder for batch processing multiple saves with rules
// Will support templated modifications across save files
func (c *CLI) handleBatchCommand(file, operation, output string) error {
	return fmt.Errorf("CLI batch command not yet implemented (Phase 4)")
}

// handleScript runs a script on a save file
// TODO: Implement script command in CLI (Phase 4)
// Placeholder for Lua scripting integration via ffvi_editor/scripting package
// Will execute custom modification scripts on save files
func (c *CLI) handleScriptCommand(file, scriptFile, output string) error {
	return fmt.Errorf("CLI script command not yet implemented (Phase 4)")
}

// handleValidate validates a save file
// TODO: Implement validate command in CLI (Phase 4)
// Placeholder for save file validation and integrity checks
// Will use existing validation infrastructure from saver.go
func (c *CLI) handleValidateCommand(file string, fix bool) error {
	return fmt.Errorf("CLI validate command not yet implemented (Phase 4)")
}

// handleBackup creates a backup of a save file
// TODO: Implement backup command in CLI (Phase 4)
// Placeholder for automatic backup and restore operations
// Will integrate with cloud backup infrastructure from Phase 4
func (c *CLI) handleBackupCommand(file, output string) error {
	return fmt.Errorf("CLI backup command not yet implemented (Phase 4)")
}

// combatPackCommand exposes Combat Depth Pack helpers via CLI
func (c *CLI) combatPackCommand() error {
	fs := flag.NewFlagSet("combat-pack", flag.ExitOnError)
	mode := fs.String("mode", "help", "Mode: encounter|boss|companion|smoke|help")
	file := fs.String("file", "", "Save file path (optional, enables save manipulation)")
	zone := fs.String("zone", "", "Zone name for encounter tuner")
	encounterRate := fs.Float64("rate", 1.0, "Encounter rate multiplier (encounter mode)")
	eliteChance := fs.Float64("elite", 0.10, "Elite encounter chance (encounter mode)")
	affixes := fs.String("affixes", "", "Comma-separated affixes (boss mode)")
	profile := fs.String("profile", "", "Profile name (companion mode)")
	risk := fs.String("risk", "normal", "Risk tolerance (companion mode)")

	if err := fs.Parse(c.args[1:]); err != nil {
		return err
	}

	// Load save if provided
	var save *pr.PR
	if *file != "" {
		var err error
		save, err = c.LoadSaveFile(*file)
		if err != nil {
			return fmt.Errorf("failed to load save: %w", err)
		}
		fmt.Printf("Loaded save: %s\n", *file)
	}

	buildResult := func(res scripting.LuaResult, err error) error {
		if err != nil {
			return err
		}
		if res == nil {
			fmt.Println("OK (no table returned)")
			return nil
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(res)
	}

	switch strings.ToLower(*mode) {
	case "encounter":
		if *zone == "" {
			return fmt.Errorf("--zone is required for encounter mode")
		}
		code := scripting.BuildEncounterScript(*zone, strconv.FormatFloat(*encounterRate, 'f', -1, 64), strconv.FormatFloat(*eliteChance, 'f', -1, 64))
		fmt.Printf("Running encounter tuning for zone=%s rate=%.2f elite=%.2f\n", *zone, *encounterRate, *eliteChance)
		return buildResult(scripting.RunSnippetWithSave(context.Background(), code, save))
	case "boss":
		if *affixes == "" {
			return fmt.Errorf("--affixes is required for boss mode")
		}
		code := scripting.BuildBossScript(*affixes)
		fmt.Printf("Running boss remix with affixes=%s\n", *affixes)
		return buildResult(scripting.RunSnippetWithSave(context.Background(), code, save))
	case "companion":
		if *profile == "" {
			return fmt.Errorf("--profile is required for companion mode")
		}
		code := scripting.BuildCompanionScript(*profile, *risk)
		fmt.Printf("Running companion director profile=%s risk=%s\n", *profile, *risk)
		return buildResult(scripting.RunSnippetWithSave(context.Background(), code, save))
	case "smoke":
		fmt.Println("Running Combat Depth Pack smoke tests")
		return buildResult(scripting.RunSnippetWithSave(context.Background(), scripting.BuildSmokeScript(), save))
	default:
		fmt.Println("Combat Depth Pack CLI")
		fmt.Println("  --mode encounter   --zone 'Mt. Kolts' --rate 1.2 --elite 0.15 [--file save.json]")
		fmt.Println("  --mode boss        --affixes enraged,glass_cannon [--file save.json]")
		fmt.Println("  --mode companion   --profile aggressive --risk aggressive [--file save.json]")
		fmt.Println("  --mode smoke       [--file save.json] (runs plugin smoke tests)")
		fmt.Println("\nNote: --file enables save data manipulation via Lua bindings")
		return nil
	}
}
