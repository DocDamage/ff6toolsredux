# Enhanced Plugin API Examples

## Example 1: Character Stats Display Plugin

```go
package main

import (
    "context"
    "fmt"
    "ffvi_editor/plugins"
    "ffvi_editor/models"
)

// StatsDisplayPlugin shows comprehensive character information
type StatsDisplayPlugin struct {
    plugins.BasePlugin
}

func (p *StatsDisplayPlugin) Execute(api plugins.API) error {
    ctx := context.Background()
    
    // Get Terra's full stats
    terra, err := api.GetCharacter(ctx, "Terra")
    if err != nil {
        return fmt.Errorf("failed to get Terra: %w", err)
    }
    
    // Display comprehensive stats
    fmt.Printf("=== %s (Level %d) ===\n", terra.Name, terra.Level)
    fmt.Printf("HP: %d/%d\n", terra.HP.Current, terra.HP.Max)
    fmt.Printf("MP: %d/%d\n", terra.MP.Current, terra.MP.Max)
    fmt.Printf("EXP: %d\n", terra.Exp)
    fmt.Printf("\nCombat Stats:\n")
    fmt.Printf("  Vigor: %d\n", terra.Vigor)
    fmt.Printf("  Stamina: %d\n", terra.Stamina)
    fmt.Printf("  Speed: %d\n", terra.Speed)
    fmt.Printf("  Magic: %d\n", terra.Magic)
    
    // Display equipment
    fmt.Printf("\nEquipment:\n")
    fmt.Printf("  Weapon: ID %d\n", terra.Equipment.WeaponID)
    fmt.Printf("  Shield: ID %d\n", terra.Equipment.ShieldID)
    fmt.Printf("  Armor: ID %d\n", terra.Equipment.ArmorID)
    fmt.Printf("  Helmet: ID %d\n", terra.Equipment.HelmetID)
    fmt.Printf("  Relic 1: ID %d\n", terra.Equipment.Relic1ID)
    fmt.Printf("  Relic 2: ID %d\n", terra.Equipment.Relic2ID)
    
    // Display learned spells
    fmt.Printf("\nLearned Spells: %d\n", len(terra.SpellsByID))
    for id, spell := range terra.SpellsByID {
        if spell.Value > 0 {
            fmt.Printf("  %s (ID:%d) - Mastery: %d%%\n", 
                spell.Name, id, spell.Value)
        }
    }
    
    // Display commands
    fmt.Printf("\nCommands:\n")
    for _, cmd := range terra.Commands {
        fmt.Printf("  - %s\n", cmd.Name)
    }
    
    return nil
}
```

## Example 2: Low HP Party Alert Plugin

```go
// LowHPAlertPlugin warns about characters with low HP
type LowHPAlertPlugin struct {
    plugins.BasePlugin
    Threshold float64 // HP percentage threshold (e.g., 0.5 for 50%)
}

func (p *LowHPAlertPlugin) Execute(api plugins.API) error {
    ctx := context.Background()
    
    // Find all characters with HP below threshold
    var lowHPChars []*models.Character
    
    characterNames := []string{
        "Terra", "Locke", "Edgar", "Sabin", "Celes", "Shadow",
        "Cyan", "Gau", "Setzer", "Strago", "Relm", "Mog",
        "Gogo", "Umaro",
    }
    
    for _, name := range characterNames {
        char, err := api.GetCharacter(ctx, name)
        if err != nil || !char.IsEnabled {
            continue
        }
        
        // Check HP percentage
        hpPercent := float64(char.HP.Current) / float64(char.HP.Max)
        if hpPercent < p.Threshold {
            lowHPChars = append(lowHPChars, char)
        }
    }
    
    // Display alerts
    if len(lowHPChars) > 0 {
        fmt.Printf("⚠️  WARNING: %d character(s) have low HP!\n", len(lowHPChars))
        for _, char := range lowHPChars {
            hpPercent := float64(char.HP.Current) / float64(char.HP.Max) * 100
            fmt.Printf("  - %s: %d/%d HP (%.1f%%)\n",
                char.Name, char.HP.Current, char.HP.Max, hpPercent)
        }
    } else {
        fmt.Println("✅ All party members have healthy HP levels")
    }
    
    return nil
}
```

## Example 3: Equipment Audit Plugin

```go
// EquipmentAuditPlugin checks for unequipped slots
type EquipmentAuditPlugin struct {
    plugins.BasePlugin
}

func (p *EquipmentAuditPlugin) Execute(api plugins.API) error {
    ctx := context.Background()
    
    // Default "empty" equipment IDs
    const (
        EmptyWeapon = 93
        EmptyShield = 93
        EmptyArmor  = 199
        EmptyHelmet = 198
        EmptyRelic  = 200
    )
    
    characterNames := []string{
        "Terra", "Locke", "Edgar", "Sabin", "Celes", "Shadow",
        "Cyan", "Gau", "Setzer", "Strago", "Relm", "Mog",
        "Gogo", "Umaro",
    }
    
    fmt.Println("=== Equipment Audit ===")
    
    for _, name := range characterNames {
        char, err := api.GetCharacter(ctx, name)
        if err != nil || !char.IsEnabled {
            continue
        }
        
        // Check for empty slots
        var emptySlots []string
        if char.Equipment.WeaponID == EmptyWeapon {
            emptySlots = append(emptySlots, "Weapon")
        }
        if char.Equipment.ShieldID == EmptyShield {
            emptySlots = append(emptySlots, "Shield")
        }
        if char.Equipment.ArmorID == EmptyArmor {
            emptySlots = append(emptySlots, "Armor")
        }
        if char.Equipment.HelmetID == EmptyHelmet {
            emptySlots = append(emptySlots, "Helmet")
        }
        if char.Equipment.Relic1ID == EmptyRelic {
            emptySlots = append(emptySlots, "Relic 1")
        }
        if char.Equipment.Relic2ID == EmptyRelic {
            emptySlots = append(emptySlots, "Relic 2")
        }
        
        if len(emptySlots) > 0 {
            fmt.Printf("⚠️  %s has empty slots: %v\n", char.Name, emptySlots)
        }
    }
    
    return nil
}
```

## Example 4: Spell Coverage Analyzer

```go
// SpellCoveragePlugin analyzes spell distribution across party
type SpellCoveragePlugin struct {
    plugins.BasePlugin
}

func (p *SpellCoveragePlugin) Execute(api plugins.API) error {
    ctx := context.Background()
    
    // Track which characters know each spell
    spellKnowledge := make(map[int][]string) // spell ID -> character names
    
    characterNames := []string{
        "Terra", "Locke", "Edgar", "Sabin", "Celes", "Shadow",
        "Cyan", "Gau", "Setzer", "Strago", "Relm", "Mog",
        "Gogo", "Umaro",
    }
    
    for _, name := range characterNames {
        char, err := api.GetCharacter(ctx, name)
        if err != nil || !char.IsEnabled {
            continue
        }
        
        // Check learned spells
        for spellID, spell := range char.SpellsByID {
            if spell.Value > 0 {
                spellKnowledge[spellID] = append(spellKnowledge[spellID], char.Name)
            }
        }
    }
    
    // Analyze coverage
    fmt.Println("=== Spell Coverage Analysis ===")
    
    // Find unique spells
    fmt.Printf("\nTotal spells available: %d\n", len(spellKnowledge))
    
    // Find spells known by only one character
    var rareSpells []string
    for spellID, knowers := range spellKnowledge {
        if len(knowers) == 1 {
            // Get spell details from first character who knows it
            char, _ := api.GetCharacter(ctx, knowers[0])
            if spell, ok := char.SpellsByID[spellID]; ok {
                rareSpells = append(rareSpells, 
                    fmt.Sprintf("%s (only %s)", spell.Name, knowers[0]))
            }
        }
    }
    
    if len(rareSpells) > 0 {
        fmt.Printf("\n⚠️  Spells known by only 1 character:\n")
        for _, s := range rareSpells {
            fmt.Printf("  - %s\n", s)
        }
    }
    
    // Find most versatile mage
    spellCounts := make(map[string]int)
    for _, name := range characterNames {
        char, err := api.GetCharacter(ctx, name)
        if err != nil || !char.IsEnabled {
            continue
        }
        
        count := 0
        for _, spell := range char.SpellsByID {
            if spell.Value > 0 {
                count++
            }
        }
        if count > 0 {
            spellCounts[char.Name] = count
        }
    }
    
    fmt.Printf("\n📊 Spell counts per character:\n")
    for name, count := range spellCounts {
        fmt.Printf("  %s: %d spells\n", name, count)
    }
    
    return nil
}
```

## Example 5: Stat Optimizer Finder

```go
// StatOptimizerPlugin finds characters optimal for specific roles
type StatOptimizerPlugin struct {
    plugins.BasePlugin
}

func (p *StatOptimizerPlugin) Execute(api plugins.API) error {
    ctx := context.Background()
    
    // Define role requirements
    type Role struct {
        Name        string
        StatCheck   func(*models.Character) int
        Description string
    }
    
    roles := []Role{
        {
            Name: "Physical Attacker",
            StatCheck: func(c *models.Character) int {
                return c.Vigor + c.Speed/2
            },
            Description: "High Vigor + Speed",
        },
        {
            Name: "Tank",
            StatCheck: func(c *models.Character) int {
                return c.Stamina + c.HP.Max/50
            },
            Description: "High Stamina + HP",
        },
        {
            Name: "Mage",
            StatCheck: func(c *models.Character) int {
                return c.Magic + len(c.SpellsByID)/2
            },
            Description: "High Magic + Spell count",
        },
        {
            Name: "Speed Runner",
            StatCheck: func(c *models.Character) int {
                return c.Speed
            },
            Description: "Highest Speed",
        },
    }
    
    characterNames := []string{
        "Terra", "Locke", "Edgar", "Sabin", "Celes", "Shadow",
        "Cyan", "Gau", "Setzer", "Strago", "Relm", "Mog",
    }
    
    fmt.Println("=== Optimal Characters by Role ===")
    
    for _, role := range roles {
        var bestChar *models.Character
        var bestScore int
        
        for _, name := range characterNames {
            char, err := api.GetCharacter(ctx, name)
            if err != nil || !char.IsEnabled {
                continue
            }
            
            score := role.StatCheck(char)
            if score > bestScore {
                bestScore = score
                bestChar = char
            }
        }
        
        if bestChar != nil {
            fmt.Printf("\n🏆 Best %s: %s (Score: %d)\n", 
                role.Name, bestChar.Name, bestScore)
            fmt.Printf("   %s\n", role.Description)
            fmt.Printf("   Stats - Vigor:%d Stamina:%d Speed:%d Magic:%d\n",
                bestChar.Vigor, bestChar.Stamina, 
                bestChar.Speed, bestChar.Magic)
        }
    }
    
    return nil
}
```

## Plugin Configuration Examples

### Plugin Manifest for Stats Display
```json
{
  "name": "Character Stats Display",
  "version": "1.0.0",
  "author": "FF6 Community",
  "description": "Displays comprehensive character information",
  "permissions": ["read_save"],
  "entry_point": "stats_display.lua",
  "enabled": true
}
```

### Plugin Manifest for HP Alert
```json
{
  "name": "Low HP Alert",
  "version": "1.0.0",
  "author": "FF6 Community",
  "description": "Warns when characters have low HP",
  "permissions": ["read_save"],
  "entry_point": "hp_alert.lua",
  "config": {
    "threshold": 0.5
  },
  "enabled": true
}
```

### Plugin Manifest for Equipment Audit
```json
{
  "name": "Equipment Audit",
  "version": "1.0.0",
  "author": "FF6 Community",
  "description": "Checks for unequipped character slots",
  "permissions": ["read_save"],
  "entry_point": "equipment_audit.lua",
  "enabled": true
}
```

## API Usage Summary

### Reading Character Data
```go
// Get single character
char, err := api.GetCharacter(ctx, "Terra")

// Find character by criteria
lowHPChar := api.FindCharacter(ctx, func(c *models.Character) bool {
    return c.HP.Current < c.HP.Max/2
})

// Search inventory
items := api.FindItems(ctx, func(item *models.InventoryRow) bool {
    return item.Count > 10
})
```

### Available Character Fields
```go
char.Name           // Character name
char.Level          // Current level
char.Exp            // Experience points
char.HP.Current     // Current HP
char.HP.Max         // Maximum HP
char.MP.Current     // Current MP
char.MP.Max         // Maximum MP
char.Vigor          // Physical attack
char.Stamina        // Physical defense
char.Speed          // Turn frequency
char.Magic          // Magical power
char.Equipment      // Equipment slots struct
char.SpellsByID     // Map of learned spells
char.Commands       // Array of commands
char.IsEnabled      // In party?
```

### Permission Requirements
All read operations require: `CommonPermissions.ReadSave`

Write operations (future) will require: `CommonPermissions.WriteSave`

## Best Practices

1. **Always check errors** when calling API methods
2. **Verify character is enabled** before processing
3. **Use FindCharacter** for complex searches instead of iterating manually
4. **Check spell Value > 0** to confirm spell is learned (not just known)
5. **Cache character data** if using multiple times to avoid repeated API calls
6. **Handle missing characters gracefully** (not all saves have all characters recruited)

---

These examples demonstrate the power of the enhanced plugin API for analyzing and reporting on FF6 save data.
