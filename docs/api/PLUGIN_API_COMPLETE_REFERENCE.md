# Plugin API Complete Reference - Phase 4D

## Quick Start: Read & Write Character Data

### Reading Character Data
```go
// Get character with full stats
char, err := api.GetCharacter(ctx, "Terra")
if err != nil {
    return err
}

// Access all character data
fmt.Printf("HP: %d/%d\n", char.HP.Current, char.HP.Max)
fmt.Printf("Level: %d (EXP: %d)\n", char.Level, char.Exp)
fmt.Printf("Vigor: %d, Magic: %d\n", char.Vigor, char.Magic)
fmt.Printf("Weapon ID: %d\n", char.Equipment.WeaponID)
fmt.Printf("Spells learned: %d\n", len(char.SpellsByID))
```

### Writing Character Data
```go
// Modify character data
char.HP.Current = 9999
char.HP.Max = 9999
char.MP.Current = 999
char.MP.Max = 999
char.Level = 99
char.Vigor = 128
char.Stamina = 128
char.Speed = 128
char.Magic = 128

// Update equipment
char.Equipment.WeaponID = 168  // Ultima Weapon
char.Equipment.ShieldID = 52   // Flame Shield

// Learn spells
if spell, ok := char.SpellsByID[54]; ok {
    spell.Value = 100  // Master Ultima spell
}

// Save changes
err = api.SetCharacter(ctx, "Terra", char)
```

## Complete API Reference

### Character Operations (100% Complete)

#### GetCharacter(ctx, name) ✅
**Purpose:** Read complete character data

**Returns:**
- Character struct with all fields populated
- HP/MP (current + max)
- Level, Experience
- Combat stats (Vigor, Stamina, Speed, Magic)
- Equipment (all 6 slots)
- Spells (with proficiency 0-100)
- Commands

**Example:**
```go
terra, err := api.GetCharacter(ctx, "Terra")
if err != nil {
    log.Fatal(err)
}
```

#### SetCharacter(ctx, name, character) ✅
**Purpose:** Write complete character data

**Updates:**
- All basic fields (name, enabled, experience)
- HP/MP (current + max)
- Level
- Combat stats (Vigor, Stamina, Speed, Magic)
- Equipment (all 6 slots)
- Spells (with proficiency)
- Commands

**Example:**
```go
// Modify character
terra.Level = 99
terra.HP.Max = 9999
terra.Equipment.WeaponID = 168

// Save
err = api.SetCharacter(ctx, "Terra", terra)
```

**Requires:** `CommonPermissions.WriteSave`

#### FindCharacter(ctx, predicate) ✅
**Purpose:** Search characters with custom criteria

**Example:**
```go
// Find character with low HP
lowHP := api.FindCharacter(ctx, func(c *models.Character) bool {
    return c.HP.Current < (c.HP.Max / 2)
})

// Find character with specific equipment
hasUltima := api.FindCharacter(ctx, func(c *models.Character) bool {
    return c.Equipment.WeaponID == 168
})

// Find character who knows specific spell
knowsUltima := api.FindCharacter(ctx, func(c *models.Character) bool {
    spell, ok := c.SpellsByID[54]
    return ok && spell.Value > 0
})
```

## Character Fields Reference

### Basic Info
```go
char.ID             // int - Character ID (0-13)
char.Name           // string - Character name
char.RootName       // string - Internal name
char.IsEnabled      // bool - In party?
char.IsNPC          // bool - Is NPC?
```

### Progression
```go
char.Level          // int - Current level (1-99)
char.Exp            // int - Experience points
```

### HP/MP Stats
```go
char.HP.Current     // int - Current HP
char.HP.Max         // int - Maximum HP (includes base)
char.MP.Current     // int - Current MP
char.MP.Max         // int - Maximum MP (includes base)
```

### Combat Stats
```go
char.Vigor          // int - Physical attack (0-128)
char.Stamina        // int - Physical defense (0-128)
char.Speed          // int - Turn frequency (0-128)
char.Magic          // int - Magical power (0-128)
```

### Equipment
```go
char.Equipment.WeaponID     // int - Weapon ID
char.Equipment.ShieldID     // int - Shield ID
char.Equipment.ArmorID      // int - Armor ID
char.Equipment.HelmetID     // int - Helmet ID
char.Equipment.Relic1ID     // int - First relic ID
char.Equipment.Relic2ID     // int - Second relic ID
```

**Empty Equipment IDs:**
- Weapon/Shield: 93
- Armor: 199
- Helmet: 198
- Relics: 200

### Spells
```go
char.SpellsByID     // map[int]*models.Spell - Spell ID → Spell
char.SpellsByIndex  // []*models.Spell - Spells by index
char.SpellsSorted   // []*models.Spell - Spells sorted by name
```

**Spell Struct:**
```go
spell.Name          // string - Spell name
spell.Index         // int - Spell index
spell.Value         // int - Proficiency (0-100, 0=not learned)
```

**Spell ID Range:** 1-54

### Commands
```go
char.Commands       // []*models.Command - Array of commands
```

**Command Struct:**
```go
cmd.Name            // string - Command name
cmd.Value           // int - Command ID
cmd.SortedIndex     // int - Sorted position
```

## Common Patterns

### Max Out Character
```go
char, _ := api.GetCharacter(ctx, "Terra")

// Max stats
char.Level = 99
char.HP.Current = 9999
char.HP.Max = 9999
char.MP.Current = 999
char.MP.Max = 999
char.Vigor = 128
char.Stamina = 128
char.Speed = 128
char.Magic = 128

api.SetCharacter(ctx, "Terra", char)
```

### Heal Character
```go
char, _ := api.GetCharacter(ctx, "Terra")
char.HP.Current = char.HP.Max
char.MP.Current = char.MP.Max
api.SetCharacter(ctx, "Terra", char)
```

### Equip Best Gear
```go
char, _ := api.GetCharacter(ctx, "Terra")
char.Equipment.WeaponID = 168  // Ultima Weapon
char.Equipment.ShieldID = 52   // Flame Shield
char.Equipment.ArmorID = 196   // Snow Scarf
char.Equipment.HelmetID = 198  // No helmet
char.Equipment.Relic1ID = 208  // Ribbon
char.Equipment.Relic2ID = 230  // Genji Glove
api.SetCharacter(ctx, "Terra", char)
```

### Learn All Spells
```go
char, _ := api.GetCharacter(ctx, "Terra")

// Learn all spells at 100% proficiency
for spellID := 1; spellID <= 54; spellID++ {
    if spell, ok := char.SpellsByID[spellID]; ok {
        spell.Value = 100
    }
}

api.SetCharacter(ctx, "Terra", char)
```

### Level Up Character
```go
char, _ := api.GetCharacter(ctx, "Terra")

// Level up
char.Level = 50
char.Exp = 500000

// Increase HP/MP with level
char.HP.Max += 200
char.MP.Max += 20

// Increase stats
char.Vigor += 5
char.Stamina += 4
char.Speed += 3
char.Magic += 6

api.SetCharacter(ctx, "Terra", char)
```

### Copy Equipment Between Characters
```go
source, _ := api.GetCharacter(ctx, "Terra")
target, _ := api.GetCharacter(ctx, "Celes")

// Copy equipment
target.Equipment = source.Equipment

api.SetCharacter(ctx, "Celes", target)
```

### Remove All Equipment
```go
char, _ := api.GetCharacter(ctx, "Terra")

// Set to empty equipment IDs
char.Equipment.WeaponID = 93
char.Equipment.ShieldID = 93
char.Equipment.ArmorID = 199
char.Equipment.HelmetID = 198
char.Equipment.Relic1ID = 200
char.Equipment.Relic2ID = 200

api.SetCharacter(ctx, "Terra", char)
```

## Plugin Examples

### Example 1: Max Stats Plugin
```go
func MaxStatsPlugin(api plugins.API) error {
    ctx := context.Background()
    
    characters := []string{
        "Terra", "Locke", "Edgar", "Sabin", "Celes", "Shadow",
        "Cyan", "Gau", "Setzer", "Strago", "Relm", "Mog",
    }
    
    for _, name := range characters {
        char, err := api.GetCharacter(ctx, name)
        if err != nil || !char.IsEnabled {
            continue
        }
        
        // Max everything
        char.Level = 99
        char.HP.Current = 9999
        char.HP.Max = 9999
        char.MP.Current = 999
        char.MP.Max = 999
        char.Vigor = 128
        char.Stamina = 128
        char.Speed = 128
        char.Magic = 128
        
        if err := api.SetCharacter(ctx, name, char); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Example 2: HP/MP Restore Plugin
```go
func RestoreAllPlugin(api plugins.API) error {
    ctx := context.Background()
    
    characters := []string{
        "Terra", "Locke", "Edgar", "Sabin", "Celes", "Shadow",
        "Cyan", "Gau", "Setzer", "Strago", "Relm", "Mog",
    }
    
    healed := 0
    for _, name := range characters {
        char, err := api.GetCharacter(ctx, name)
        if err != nil || !char.IsEnabled {
            continue
        }
        
        // Restore HP/MP
        char.HP.Current = char.HP.Max
        char.MP.Current = char.MP.Max
        
        if err := api.SetCharacter(ctx, name, char); err != nil {
            return err
        }
        healed++
    }
    
    fmt.Printf("Restored HP/MP for %d characters\n", healed)
    return nil
}
```

### Example 3: Auto-Equip Plugin
```go
func AutoEquipPlugin(api plugins.API) error {
    ctx := context.Background()
    
    // Best equipment IDs
    bestGear := struct {
        Weapon  int
        Shield  int
        Armor   int
        Helmet  int
        Relic1  int
        Relic2  int
    }{
        Weapon: 168,  // Ultima Weapon
        Shield: 52,   // Flame Shield
        Armor:  196,  // Snow Scarf
        Helmet: 198,  // None
        Relic1: 208,  // Ribbon
        Relic2: 230,  // Genji Glove
    }
    
    characters := []string{
        "Terra", "Locke", "Edgar", "Sabin", "Celes", "Shadow",
    }
    
    for _, name := range characters {
        char, err := api.GetCharacter(ctx, name)
        if err != nil || !char.IsEnabled {
            continue
        }
        
        // Equip best gear
        char.Equipment.WeaponID = bestGear.Weapon
        char.Equipment.ShieldID = bestGear.Shield
        char.Equipment.ArmorID = bestGear.Armor
        char.Equipment.HelmetID = bestGear.Helmet
        char.Equipment.Relic1ID = bestGear.Relic1
        char.Equipment.Relic2ID = bestGear.Relic2
        
        if err := api.SetCharacter(ctx, name, char); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Example 4: Magic Teacher Plugin
```go
func TeachAllMagicPlugin(api plugins.API) error {
    ctx := context.Background()
    
    characters := []string{
        "Terra", "Celes", "Strago", "Relm",
    }
    
    for _, name := range characters {
        char, err := api.GetCharacter(ctx, name)
        if err != nil || !char.IsEnabled {
            continue
        }
        
        // Teach all spells at 100% proficiency
        learnedCount := 0
        for spellID := 1; spellID <= 54; spellID++ {
            if spell, ok := char.SpellsByID[spellID]; ok {
                if spell.Value == 0 {
                    learnedCount++
                }
                spell.Value = 100
            }
        }
        
        if err := api.SetCharacter(ctx, name, char); err != nil {
            return err
        }
        
        fmt.Printf("%s learned %d new spells\n", name, learnedCount)
    }
    
    return nil
}
```

### Example 5: Level Booster Plugin
```go
func LevelBoostPlugin(api plugins.API, targetLevel int) error {
    ctx := context.Background()
    
    characters := []string{
        "Terra", "Locke", "Edgar", "Sabin", "Celes", "Shadow",
        "Cyan", "Gau", "Setzer", "Strago", "Relm", "Mog",
    }
    
    for _, name := range characters {
        char, err := api.GetCharacter(ctx, name)
        if err != nil || !char.IsEnabled {
            continue
        }
        
        if char.Level < targetLevel {
            oldLevel := char.Level
            char.Level = targetLevel
            
            // Calculate stat increases (approximate)
            levelGain := targetLevel - oldLevel
            char.HP.Max += levelGain * 40
            char.MP.Max += levelGain * 5
            char.Vigor += levelGain / 4
            char.Stamina += levelGain / 4
            char.Speed += levelGain / 5
            char.Magic += levelGain / 3
            
            // Update experience (approximate)
            char.Exp = targetLevel * 10000
            
            if err := api.SetCharacter(ctx, name, char); err != nil {
                return err
            }
            
            fmt.Printf("%s leveled up from %d to %d\n", name, oldLevel, targetLevel)
        }
    }
    
    return nil
}
```

## Error Handling

### Check Permission
```go
if !api.HasPermission(CommonPermissions.WriteSave) {
    return errors.New("plugin requires write permission")
}
```

### Handle Errors
```go
char, err := api.GetCharacter(ctx, "Terra")
if err != nil {
    if err == ErrCharacterNotFound {
        fmt.Println("Character not found")
    } else if err == ErrInsufficientPermissions {
        fmt.Println("Permission denied")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return err
}
```

## Permissions Required

### Read Operations
- `CommonPermissions.ReadSave`

### Write Operations
- `CommonPermissions.WriteSave`

### Plugin Manifest
```json
{
  "name": "My Plugin",
  "version": "1.0.0",
  "permissions": ["read_save", "write_save"],
  "enabled": true
}
```

## Performance Tips

1. **Batch updates** - Get once, modify, set once
2. **Check IsEnabled** - Skip disabled characters
3. **Cache character data** - Avoid repeated API calls
4. **Use predicates** - Let FindCharacter filter for you
5. **Validate before setting** - Check HP/MP bounds

## Best Practices

1. **Always check errors** from API calls
2. **Verify permissions** before attempting writes
3. **Validate stat ranges** (HP ≤ 9999, stats ≤ 128)
4. **Backup before modifying** critical data
5. **Test with copies** of save files first
6. **Document** what your plugin modifies
7. **Provide feedback** to users via logging

---

**API Version:** 1.0 (Phase 4D Complete)  
**Status:** Production Ready ✅  
**Completeness:** 100% for Character Operations  
**Last Updated:** January 16, 2026  
