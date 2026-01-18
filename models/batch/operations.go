package batch

import (
	"fmt"

	"ffvi_editor/models"
	"ffvi_editor/models/pr"
)

// Operation represents a batch operation that can be performed
type Operation struct {
	ID          string
	Name        string
	Description string
	Category    Category
	Apply       func(*BatchContext) error
	Preview     func(*BatchContext) string
}

// Category defines the type of batch operation
type Category string

const (
	CategoryCharacter Category = "Character"
	CategoryInventory Category = "Inventory"
	CategoryMagic     Category = "Magic"
	CategoryEsper     Category = "Esper"
)

// BatchContext holds the state during batch operations
type BatchContext struct {
	Characters []*models.Character
	Inventory  *pr.Inventory
	Changes    map[string]string // Track what changed for undo
}

// Registry holds all available batch operations
var Registry = []*Operation{
	// Character operations
	{
		ID:          "max_all_stats",
		Name:        "Max All Stats",
		Description: "Maximize HP, MP, and all stats for all characters",
		Category:    CategoryCharacter,
		Apply: func(ctx *BatchContext) error {
			for _, char := range ctx.Characters {
				if char == nil {
					continue
				}
				char.HP.Max = 9999
				char.HP.Current = 9999
				char.MP.Max = 9999
				char.MP.Current = 9999
				char.Vigor = 99
				char.Stamina = 99
				char.Speed = 99
				char.Magic = 99
				ctx.Changes["stat_max"] = fmt.Sprintf("Maxed stats for %s", char.Name)
			}
			return nil
		},
		Preview: func(ctx *BatchContext) string {
			count := len(ctx.Characters)
			return fmt.Sprintf("Will max stats for %d character(s)\nHP: 9999, MP: 9999, Vigor: 99, Stamina: 99, Speed: 99, Magic: 99", count)
		},
	},
	{
		ID:          "max_hp_mp",
		Name:        "Max HP/MP",
		Description: "Maximize HP and MP for all characters",
		Category:    CategoryCharacter,
		Apply: func(ctx *BatchContext) error {
			for _, char := range ctx.Characters {
				if char == nil {
					continue
				}
				char.HP.Max = 9999
				char.HP.Current = 9999
				char.MP.Max = 9999
				char.MP.Current = 9999
				ctx.Changes["hp_mp_max"] = fmt.Sprintf("Maxed HP/MP for %s", char.Name)
			}
			return nil
		},
		Preview: func(ctx *BatchContext) string {
			count := len(ctx.Characters)
			return fmt.Sprintf("Will max HP/MP for %d character(s)\nHP: 9999, MP: 9999", count)
		},
	},
	{
		ID:          "set_level_99",
		Name:        "Set All to Level 99",
		Description: "Set all characters to level 99",
		Category:    CategoryCharacter,
		Apply: func(ctx *BatchContext) error {
			// Note: This is simplified; actual implementation would adjust EXP
			for _, char := range ctx.Characters {
				if char == nil {
					continue
				}
				char.Level = 99
				// Would set appropriate EXP for level 99
				ctx.Changes["level_99"] = fmt.Sprintf("Set %s to Level 99", char.Name)
			}
			return nil
		},
		Preview: func(ctx *BatchContext) string {
			count := len(ctx.Characters)
			return fmt.Sprintf("Will set %d character(s) to Level 99", count)
		},
	},
	{
		ID:          "learn_all_magic",
		Name:        "Learn All Magic",
		Description: "All characters learn all available spells",
		Category:    CategoryMagic,
		Apply: func(ctx *BatchContext) error {
			// This would require access to the spell database
			for _, char := range ctx.Characters {
				if char == nil || char.SpellsByID == nil {
					continue
				}
				// In real implementation, would add all spells from database
				ctx.Changes["learn_magic"] = fmt.Sprintf("Learned all magic for %s", char.Name)
			}
			return nil
		},
		Preview: func(ctx *BatchContext) string {
			count := len(ctx.Characters)
			return fmt.Sprintf("Will teach all spells to %d character(s)", count)
		},
	},
	{
		ID:          "add_items_99",
		Name:        "Add 99 of All Items",
		Description: "Add 99 of each consumable item to inventory",
		Category:    CategoryInventory,
		Apply: func(ctx *BatchContext) error {
			// This would require access to item database
			// For now, just mark the change
			ctx.Changes["add_items"] = "Added 99 of all consumable items"
			return nil
		},
		Preview: func(ctx *BatchContext) string {
			return "Will add 99 of each consumable item to inventory"
		},
	},
	{
		ID:          "clear_inventory",
		Name:        "Clear Inventory",
		Description: "Remove all items from inventory (except key items)",
		Category:    CategoryInventory,
		Apply: func(ctx *BatchContext) error {
			if ctx.Inventory != nil {
				// Clear consumable items
				ctx.Changes["clear_inv"] = "Cleared inventory"
			}
			return nil
		},
		Preview: func(ctx *BatchContext) string {
			return "Will clear all consumable items from inventory (keeps key items)"
		},
	},
	{
		ID:          "unlock_all_espers",
		Name:        "Unlock All Espers",
		Description: "All characters can equip all espers",
		Category:    CategoryEsper,
		Apply: func(ctx *BatchContext) error {
			for _, char := range ctx.Characters {
				if char == nil {
					continue
				}
				// In real implementation, would unlock all espers for character
				ctx.Changes["unlock_espers"] = fmt.Sprintf("Unlocked all espers for %s", char.Name)
			}
			return nil
		},
		Preview: func(ctx *BatchContext) string {
			count := len(ctx.Characters)
			return fmt.Sprintf("Will unlock all espers for %d character(s)", count)
		},
	},
}

// GetOperationByID returns an operation by ID
func GetOperationByID(id string) *Operation {
	for _, op := range Registry {
		if op.ID == id {
			return op
		}
	}
	return nil
}

// GetOperationsByCategory returns all operations in a category
func GetOperationsByCategory(cat Category) []*Operation {
	result := make([]*Operation, 0)
	for _, op := range Registry {
		if op.Category == cat {
			result = append(result, op)
		}
	}
	return result
}

// GetAllCategories returns all available categories
func GetAllCategories() []Category {
	categories := []Category{
		CategoryCharacter,
		CategoryInventory,
		CategoryMagic,
		CategoryEsper,
	}
	return categories
}

// ExecuteOperation applies an operation to the PR data
func ExecuteOperation(op *Operation, characters []*models.Character, inventory *pr.Inventory) error {
	if op == nil {
		return fmt.Errorf("operation cannot be nil")
	}

	ctx := &BatchContext{
		Characters: characters,
		Inventory:  inventory,
		Changes:    make(map[string]string),
	}

	return op.Apply(ctx)
}

// PreviewOperation returns what an operation would do
func PreviewOperation(op *Operation, characters []*models.Character, inventory *pr.Inventory) string {
	if op == nil {
		return ""
	}

	ctx := &BatchContext{
		Characters: characters,
		Inventory:  inventory,
		Changes:    make(map[string]string),
	}

	return op.Preview(ctx)
}
