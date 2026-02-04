package cli

import (
	"fmt"

	pri "ffvi_editor/models/pr"
)

// handleEditCommand loads a save, modifies it, and saves it
func (c *CLI) handleEditCommand(file string, charID, level, hp, mp int, output string) error {
	// Load the save file
	save, err := c.LoadSaveFile(file)
	if err != nil {
		return fmt.Errorf("failed to load save file: %w", err)
	}

	// Determine output path
	outputPath := output
	if outputPath == "" {
		outputPath = file
	}

	// If character ID is specified, modify that character
	if charID >= 0 {
		character := pri.GetCharacterByID(charID)
		if character == nil {
			return fmt.Errorf("character with ID %d not found", charID)
		}

		// Apply modifications if values are >= 0
		if level >= 0 {
			character.Level = level
			fmt.Printf("Set %s level to %d\n", character.Name, level)
		}
		if hp >= 0 {
			character.HP.Current = hp
			character.HP.Max = hp
			fmt.Printf("Set %s HP to %d\n", character.Name, hp)
		}
		if mp >= 0 {
			character.MP.Current = mp
			character.MP.Max = mp
			fmt.Printf("Set %s MP to %d\n", character.Name, mp)
		}
	}

	// Save the modified file
	if err := c.SaveSaveFile(save, outputPath); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("Successfully edited save file: %s\n", outputPath)
	return nil
}
