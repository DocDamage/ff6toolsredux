package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"ffvi_editor/scripting"
)

// handleScriptCommand runs a script on a save file
// Loads a save file, executes a Lua script with save data bindings, and saves the result
func (c *CLI) handleScriptCommand(file, scriptFile, output string) error {
	// Load the save file
	save, err := c.LoadSaveFile(file)
	if err != nil {
		return fmt.Errorf("failed to load save file: %w", err)
	}

	// Read the Lua script from file
	code, err := os.ReadFile(scriptFile)
	if err != nil {
		return fmt.Errorf("failed to read script file: %w", err)
	}

	// Determine output path (defaults to input file)
	outputPath := output
	if outputPath == "" {
		outputPath = file
	}

	// Execute the Lua script with save context
	fmt.Printf("Running script: %s\n", scriptFile)
	res, err := scripting.RunSnippetWithSave(context.Background(), string(code), save)
	if err != nil {
		return fmt.Errorf("script execution failed: %w", err)
	}

	// Print script output if a table was returned
	if res != nil {
		fmt.Println("Script returned:")
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(res)
	}

	// Save the modified file
	if err := c.SaveSaveFile(save, outputPath); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("Script executed successfully. Output saved to: %s\n", outputPath)
	return nil
}
