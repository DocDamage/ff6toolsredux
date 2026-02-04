package cli

import (
	"fmt"

	"ffvi_editor/io/pr"
	pri "ffvi_editor/models/pr"
)

// handleBatchCommand performs batch operations on a save file
// Supports: max-stats, max-items, max-magic, max-all
func (c *CLI) handleBatchCommand(file, operation, output string) error {
	// Load the save file
	p := pr.New()
	if err := p.Load(file, 0); err != nil {
		return fmt.Errorf("failed to load save file: %w", err)
	}

	// Determine output path
	outputPath := output
	if outputPath == "" {
		outputPath = file
	}

	// Validate and apply operation
	switch operation {
	case "max-stats":
		if err := applyMaxStats(); err != nil {
			return fmt.Errorf("failed to apply max-stats: %w", err)
		}
		fmt.Println("Applied max-stats: All character stats set to 255")
	case "max-items":
		if err := applyMaxItems(); err != nil {
			return fmt.Errorf("failed to apply max-items: %w", err)
		}
		fmt.Println("Applied max-items: All inventory items set to max quantity")
	case "max-magic":
		if err := applyMaxMagic(); err != nil {
			return fmt.Errorf("failed to apply max-magic: %w", err)
		}
		fmt.Println("Applied max-magic: All spells learned for all characters")
	case "max-all":
		if err := applyMaxStats(); err != nil {
			return fmt.Errorf("failed to apply max-stats: %w", err)
		}
		if err := applyMaxItems(); err != nil {
			return fmt.Errorf("failed to apply max-items: %w", err)
		}
		if err := applyMaxMagic(); err != nil {
			return fmt.Errorf("failed to apply max-magic: %w", err)
		}
		fmt.Println("Applied max-all: All stats, items, and magic maximized")
	default:
		return fmt.Errorf("unknown operation: %s (valid: max-stats, max-items, max-magic, max-all)", operation)
	}

	// Save the modified file
	if err := p.Save(0, outputPath, 0); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("Successfully saved to: %s\n", outputPath)
	return nil
}

// applyMaxStats sets all character stats to maximum (255)
func applyMaxStats() error {
	for _, char := range pri.Characters {
		if char == nil || char.IsNPC {
			continue
		}
		char.Vigor = 255
		char.Stamina = 255
		char.Speed = 255
		char.Magic = 255
	}
	return nil
}

// applyMaxItems sets all inventory items to max quantity (99)
func applyMaxItems() error {
	// Get regular inventory and set all items to 99
	inv := pri.GetInventory()
	for _, row := range inv.GetRows() {
		if row.ItemID > 0 {
			row.Count = 99
		}
	}
	return nil
}

// applyMaxMagic teaches all spells to all characters
func applyMaxMagic() error {
	for _, char := range pri.Characters {
		if char == nil || char.IsNPC {
			continue
		}
		// Learn all spells (set spell value to indicate learned)
		for _, spell := range char.SpellsByIndex {
			if spell != nil {
				spell.Value = 1
			}
		}
	}
	return nil
}
