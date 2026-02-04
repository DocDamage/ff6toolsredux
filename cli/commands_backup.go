package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"ffvi_editor/io/pr"
	"ffvi_editor/scripting"
)

// handleBackupCommand creates a backup of a save file
// Creates a timestamped backup copy of the specified save file
func (c *CLI) handleBackupCommand(file, output string) error {
	// Check if source file exists
	info, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source file not found: %s", file)
		}
		return fmt.Errorf("failed to access source file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("source path is a directory, not a file: %s", file)
	}

	// Generate backup filename if output not specified
	backupPath := output
	if backupPath == "" {
		timestamp := time.Now().Format("20060102_150405")
		backupPath = fmt.Sprintf("%s.%s.backup", file, timestamp)
	}

	// Read source file
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Write backup file
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	fmt.Printf("Backup created successfully: %s -> %s\n", file, backupPath)
	return nil
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
