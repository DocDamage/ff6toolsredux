# Natural Magic Block Enforcer Plugin

Enforce the "Natural Magic Block" challenge run rules - remove all magic from FF6 for a physical-only playthrough.

## Overview

This plugin helps you set up and validate the popular "Natural Magic Block" (NMB) challenge run. It removes all learned magic spells from characters, helps track esper removal, and validates that your run remains magic-free.

## Features

### Challenge Enforcement
- **Remove All Magic**: Strips all learned spells from all characters
- **Backup System**: Creates automatic backup before enforcement
- **Instant Application**: One-click magic removal across entire party
- **Reversible**: Restore wizard can undo enforcement using backup

### Validation System
- **Magic Detection**: Scans all characters for learned spells
- **Violation Reporting**: Lists any characters who have magic
- **Challenge Proof Export**: Generate timestamped validation reports
- **Real-Time Checking**: Validate at any point during your playthrough

### Restore Wizard
- **Backup Restoration**: Restore all magic from pre-enforcement backup
- **Selective Recovery**: (Future) Restore specific characters or spells
- **Safety Checks**: Confirmation dialogs prevent accidental restoration

### Status Tracking
- **Enforcement Status**: Check if NMB is currently active
- **Enforcement Date**: Track when challenge was started
- **Backup Availability**: Verify backup exists for restoration

## What is Natural Magic Block?

Natural Magic Block is a popular FF6 challenge run where:
- ✅ **No Magic Spells**: Characters cannot learn or use magic
- ✅ **No Espers**: Remove all espers from inventory (prevents spell learning)
- ✅ **Physical Only**: Win using only physical attacks, items, and command abilities
- ✅ **Natural Commands**: SwdTech, Blitz, Tools, Slots, Sketch, Rage, etc. are all allowed

### Why NMB?

- **Different Strategy**: Forces you to use neglected abilities (SwdTech, Blitz, Tools)
- **Character Specialization**: Each character's unique commands become essential
- **Increased Difficulty**: No Ultima spam, no Quick cheese, no Vanish-Doom
- **Fresh Experience**: Play FF6 in a completely different way

## Usage

### Starting an NMB Run

1. **Backup Your Save**: Always work with a backup save file first
2. **Run the Plugin**: Launch Natural Magic Block Enforcer
3. **Enforce Challenge**: Select "Enforce Natural Magic Block"
4. **Confirm**: Plugin removes all magic and creates restoration backup
5. **Validate**: Use "Validate Challenge" to confirm clean slate

### During Your Run

1. **Regular Validation**: Check periodically that no magic has been learned
2. **Export Proofs**: Generate timestamped proof files for documentation
3. **Monitor Status**: Check enforcement status and backup availability

### After Your Run

1. **Final Validation**: Generate final challenge proof
2. **Optional Restore**: Use Restore Wizard to get your magic back
3. **Export Proof**: Keep timestamped proof as completion evidence

## Validation Reports

### Clean Report
```
=== Natural Magic Block Validation ===

✓ CHALLENGE VALID
No characters have learned magic.
```

### Violation Report
```
=== Natural Magic Block Validation ===

✗ VIOLATIONS DETECTED

Terra: 3 spells (Fire, Cure, Haste)
Celes: 1 spells (Ice)
```

## Technical Details

### Spell Removal Process

1. **Backup Phase**: Iterate all characters, save spell lists to backup table
2. **Removal Phase**: Set all spell flags to `false` in character data
3. **Persistence**: Write backup to `natural_magic_block/state.lua`
4. **Validation**: Scan characters to confirm all spells removed

### Backup Data Structure

```lua
{
  enforced = true,
  enforcement_date = "2024-01-15 14:30:00",
  backup = {
    ["Terra"] = {"Fire", "Ice", "Cure", "Haste"},
    ["Locke"] = {"Fire", "Cure"},
    -- etc.
  }
}
```

### Validation Algorithm

1. Iterate all 14 playable characters
2. Check `character.spells` table
3. Count learned spells (value = `true`)
4. Report violations or confirm clean state

### Export Format

```
=== FF6 Natural Magic Block Challenge Proof ===
Export Date: 2024-01-15 18:45:30
Enforcement Status: ACTIVE
Enforcement Date: 2024-01-15 14:30:00

[Validation results]
```

## API Requirements

**Current Status**: Partially functional - spell removal works, esper removal requires additional APIs.

### Working Features

1. **Spell Removal**: ✅ `character.spells[spell] = false` - Works with existing character data API
2. **Spell Detection**: ✅ Read character spell lists for validation
3. **Configuration Storage**: ✅ File I/O for backups and state

### Conceptual Features

1. **Esper Inventory Manipulation**: ⚠️ Requires inventory API to remove espers
2. **Prevent Spell Learning**: ⚠️ Requires game engine hook to block spell learning from battles
3. **In-Game Validation Hook**: ⚠️ Requires event system to validate at save points

## Warnings

⚠️ **EXPERIMENTAL PLUGIN**

This plugin fundamentally alters your save file:
- Removes ALL magic spells from ALL characters
- Changes are immediate and affect save file
- Backup system protects against accidents
- Manual esper removal required (inventory API not available)

**Recommendations**:
- Use on a backup save file first
- Validate immediately after enforcement
- Export proof regularly during your run
- Remove espers manually from inventory
- Don't save over your main save until you're sure

## Manual Steps Required

Since some APIs aren't available, you must manually:

1. **Remove Espers**: Go to inventory and remove all espers
2. **Monitor Learning**: Don't fight enemies that teach spells
3. **Avoid Magic Points**: Don't use espers in battle (since you removed them, this is automatic)

## Challenge Tips

### Strongest Characters for NMB

1. **Edgar**: Tools are incredibly powerful (Flash, Drill, Chainsaw)
2. **Sabin**: Blitz commands deal massive damage (Bum Rush, Phantom Rush)
3. **Cyan**: SwdTech is viable (Quadra Slam, Cleave)
4. **Setzer**: Slots can still win (Joker Doom works)
5. **Gau**: Rage remains effective

### Recommended Party Composition

- **Edgar** (Tools) + **Sabin** (Blitz) + **Setzer** (Slots) + **Gau** (Rage)
- Focus on characters with strong command abilities
- Avoid Terra/Celes (their strength is magic)

### Key Items

- **Atma Weapon**: Powered by HP, not magic
- **Valiant Knife**: Powered by missing HP
- **Fixed Dice**: Random damage ignores stats
- **Offering**: Double physical hits

## FAQ

**Q: Can I undo enforcement mid-run?**  
A: Yes, use the Restore Wizard. However, this invalidates your NMB challenge.

**Q: What if I accidentally learn a spell?**  
A: Run validation to detect it, then re-enforce to remove it. Document for honesty.

**Q: Do command abilities count as magic?**  
A: No! SwdTech, Blitz, Tools, etc. are all allowed in NMB.

**Q: Can I use items that replicate magic effects?**  
A: Yes! Tinctures, Fenix Downs, Remedies are all fair game.

**Q: What about Gau's Rages that use magic?**  
A: Community rule: allowed since it's natural to the enemy, not learned magic.

## Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

## License

Part of the FF6 Save Editor plugin system. See main LICENSE file.
