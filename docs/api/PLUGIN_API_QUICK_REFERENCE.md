# Plugin API Quick Reference

## Character Data Access

### Get Character by Name
```go
char, err := api.GetCharacter(ctx, "Terra")
if err != nil {
    // Handle error
}
```

### Find Character with Predicate
```go
// Find first character with HP < 50%
char := api.FindCharacter(ctx, func(c *models.Character) bool {
    return c.HP.Current < (c.HP.Max / 2)
})
```

## Available Character Fields

### Basic Info
```go
char.ID             // int - Character ID (0-13)
char.Name           // string - Character name
char.RootName       // string - Internal name
char.Level          // int - Current level
char.Exp            // int - Experience points
char.IsEnabled      // bool - In party?
char.IsNPC          // bool - Is NPC character?
```

### HP/MP Stats
```go
char.HP.Current     // int - Current HP
char.HP.Max         // int - Maximum HP (includes base offset)
char.MP.Current     // int - Current MP
char.MP.Max         // int - Maximum MP (includes base offset)
```

### Combat Stats
```go
char.Vigor          // int - Physical attack power
char.Stamina        // int - Physical defense
char.Speed          // int - Turn frequency/agility
char.Magic          // int - Magical power
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

Default empty equipment IDs:
- Weapon/Shield: 93
- Armor: 199
- Helmet: 198
- Relics: 200

### Spells
```go
char.SpellsByID     // map[int]*models.Spell - Spell ID -> Spell
char.SpellsByIndex  // []*models.Spell - Spells by index
char.SpellsSorted   // []*models.Spell - Spells sorted by name
```

Spell struct:
```go
spell.Name          // string - Spell name
spell.Index         // int - Spell index
spell.Value         // int - Proficiency/mastery level (0-100)
```

Spell ID range: 1-54

### Commands
```go
char.Commands       // []*models.Command - Array of commands
```

Command struct:
```go
cmd.Name            // string - Command name
cmd.Value           // int - Command ID
cmd.SortedIndex     // int - Sorted position
```

### Status Effects
```go
char.StatusEffects  // []*consts.NameSlotMask8 - Status conditions
```

## Common Patterns

### Check Character Has Spell
```go
if spell, ok := char.SpellsByID[54]; ok && spell.Value > 0 {
    // Character knows Ultima and has learned it
}
```

### Calculate HP Percentage
```go
hpPercent := float64(char.HP.Current) / float64(char.HP.Max) * 100
```

### Check Equipment Slot Empty
```go
const EmptyWeapon = 93
if char.Equipment.WeaponID == EmptyWeapon {
    // Weapon slot is empty
}
```

### Count Learned Spells
```go
learnedCount := 0
for _, spell := range char.SpellsByID {
    if spell.Value > 0 {
        learnedCount++
    }
}
```

### Find Best Stat
```go
// Find character with highest vigor
bestVigor := api.FindCharacter(ctx, func(c *models.Character) bool {
    if !c.IsEnabled {
        return false
    }
    // Compare against other characters to find max
    return true  // Simplified - compare in loop
})
```

## Inventory Access

### Get Full Inventory
```go
inventory := api.GetInventory(ctx)
```

### Find Items by Predicate
```go
// Find items with count > 10
items := api.FindItems(ctx, func(item *models.InventoryRow) bool {
    return item.Count > 10
})
```

## Party Access

### Get Current Party
```go
party := api.GetParty(ctx)
```

## Equipment Reference

### Get Character Equipment Info
```go
equipment := api.GetEquipment(ctx, "Terra")
```

## Permissions Required

### Read Operations
All read operations require:
```go
CommonPermissions.ReadSave
```

Check in plugin:
```go
if !api.HasPermission(CommonPermissions.ReadSave) {
    return errors.New("insufficient permissions")
}
```

## Error Handling

### Standard Pattern
```go
char, err := api.GetCharacter(ctx, "Terra")
if err != nil {
    if err == ErrCharacterNotFound {
        // Handle not found
    } else if err == ErrInsufficientPermissions {
        // Handle permission denied
    } else {
        // Handle other errors
    }
    return err
}
```

### Common Errors
```go
ErrCharacterNotFound        // Character doesn't exist
ErrInsufficientPermissions  // Missing required permission
ErrNilPRData               // Save data not loaded
```

## Character Names

Valid character names for GetCharacter():
```
"Terra", "Locke", "Edgar", "Sabin", "Celes", "Shadow",
"Cyan", "Gau", "Setzer", "Strago", "Relm", "Mog",
"Gogo", "Umaro"
```

## Example: Complete Character Analysis

```go
func analyzeCharacter(api plugins.API, name string) error {
    ctx := context.Background()
    
    // Get character
    char, err := api.GetCharacter(ctx, name)
    if err != nil {
        return err
    }
    
    // Check if in party
    if !char.IsEnabled {
        return fmt.Errorf("%s not in party", name)
    }
    
    // Analyze HP
    hpPercent := float64(char.HP.Current) / float64(char.HP.Max) * 100
    fmt.Printf("HP: %.1f%%\n", hpPercent)
    
    // Analyze equipment
    emptySlots := 0
    if char.Equipment.WeaponID == 93 { emptySlots++ }
    if char.Equipment.ShieldID == 93 { emptySlots++ }
    if char.Equipment.ArmorID == 199 { emptySlots++ }
    if char.Equipment.HelmetID == 198 { emptySlots++ }
    if char.Equipment.Relic1ID == 200 { emptySlots++ }
    if char.Equipment.Relic2ID == 200 { emptySlots++ }
    fmt.Printf("Empty equipment slots: %d\n", emptySlots)
    
    // Count spells
    spellCount := 0
    for _, spell := range char.SpellsByID {
        if spell.Value > 0 {
            spellCount++
        }
    }
    fmt.Printf("Spells learned: %d\n", spellCount)
    
    // Check best stat
    stats := map[string]int{
        "Vigor": char.Vigor,
        "Stamina": char.Stamina,
        "Speed": char.Speed,
        "Magic": char.Magic,
    }
    var bestStat string
    var bestValue int
    for stat, value := range stats {
        if value > bestValue {
            bestStat = stat
            bestValue = value
        }
    }
    fmt.Printf("Best stat: %s (%d)\n", bestStat, bestValue)
    
    return nil
}
```

## Performance Tips

1. **Cache character data** if using multiple times:
   ```go
   char, _ := api.GetCharacter(ctx, "Terra")
   // Use char multiple times instead of calling API again
   ```

2. **Use FindCharacter** for searches instead of manual iteration:
   ```go
   // Good
   char := api.FindCharacter(ctx, predicate)
   
   // Avoid
   for _, name := range allNames {
       char, _ := api.GetCharacter(ctx, name)
       if predicate(char) { break }
   }
   ```

3. **Check IsEnabled** before processing to skip unavailable characters:
   ```go
   if !char.IsEnabled {
       return nil // Skip
   }
   ```

4. **Batch similar operations** to reduce API calls

## Debugging Tips

### Print Complete Character State
```go
fmt.Printf("%+v\n", char)  // Prints all fields
```

### Check Spell IDs
```go
for id, spell := range char.SpellsByID {
    if spell.Value > 0 {
        fmt.Printf("Spell %d: %s (mastery: %d)\n", 
            id, spell.Name, spell.Value)
    }
}
```

### Verify Equipment IDs
```go
fmt.Printf("Equipment: W:%d S:%d A:%d H:%d R1:%d R2:%d\n",
    char.Equipment.WeaponID,
    char.Equipment.ShieldID,
    char.Equipment.ArmorID,
    char.Equipment.HelmetID,
    char.Equipment.Relic1ID,
    char.Equipment.Relic2ID)
```

---

**Version:** 1.0 (Phase 4C+)  
**Last Updated:** Phase 4C+ Completion  
**API Coverage:** Read operations (100%)  
