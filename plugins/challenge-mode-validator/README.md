# Challenge Mode Validator Plugin

Verify challenge run compliance with preset templates and custom rules for Final Fantasy VI.

## Features

- **9 Preset Challenge Modes**: Popular challenge templates ready to use
- **Custom Rule Creation**: Define your own challenge rules and restrictions
- **Real-time Validation**: Check compliance against active challenge rules
- **Violation Detection**: Automatically detect and log rule violations
- **Violation Log**: Timestamped log of all detected violations
- **Challenge Proof Export**: Generate verification report for sharing
- **Challenge Archiving**: Archive completed challenges with full history

## Preset Challenge Modes

### 1. Low Level Run
Complete game with all characters under level 20

### 2. Natural Magic Block
No magic learning allowed (Runic and natural spells only, no espers)

### 3. Solo Character Run
Use only one character for the entire game

### 4. No Equipment Challenge
No weapons or armor equipped

### 5. No Shop Run
Cannot buy items from shops (starting gil only)

### 6. Fixed Party Challenge
Use the same 4 characters for entire game

### 7. No Esper Run
Cannot equip espers in battle

### 8. Ancient Cave Mode
Start fresh, only use items found in dungeons

### 9. Minimalist Run
Complete with minimum battles and items

## Usage

1. **Start Challenge**: Select a preset or create custom challenge
2. **Check Compliance**: Validate save file against active challenge rules
3. **View Violations**: Review timestamped violation log
4. **Export Proof**: Generate challenge verification report
5. **End Challenge**: Archive challenge and reset for new run

## Rule Categories

- **Level Restrictions**: Maximum level caps
- **Equipment Restrictions**: Weapon, armor, esper limits
- **Magic Restrictions**: Spell learning and usage limits
- **Party Restrictions**: Party size and composition rules
- **Inventory Restrictions**: Shop purchases and item usage
- **Gameplay Restrictions**: Battle count and progression limits

## Validation Logic

The plugin checks rules against current save state:
- **max_level**: Checks all characters against level cap
- **no_learned_magic**: Verifies no esper-taught spells
- **no_espers**: Checks no espers equipped
- **max_party_size**: Validates party size limit
- **no_weapons**: Checks no weapons equipped
- **no_armor**: Checks no armor equipped

Additional rules use placeholder validation pending API support.

## File Structure

```
plugins/challenge-mode-validator/
├── metadata.json          # Plugin metadata
├── plugin.lua             # Main plugin code
├── README.md              # This file
└── CHANGELOG.md           # Version history

challenge_validator/
├── active_challenge.lua   # Current active challenge
├── violations.lua         # Violation log
├── archived_*.lua         # Completed challenge archives
└── proof_*.txt            # Exported challenge proofs
```

## Challenge Proof Format

Exported proof includes:
- Challenge name and description
- Start date and time
- Complete rule list
- Violation count and details
- Verification status (VERIFIED / NOT VERIFIED)

## API Requirements

- `GetParty()`: Party roster information
- `GetCharacter(id)`: Character stats, equipment, spells
- `ReadFile()`, `WriteFile()`: Challenge data persistence

## Permissions

- `read_save`: Read save file for validation
- `ui_display`: Display challenge interface
- `file_io`: Persist challenge data and export proofs

## Notes

- Many rules use placeholder validation pending API availability
- Real-time event tracking requires event hook API
- Some tracking features (gil spent, battles, etc.) need additional APIs
- Challenge proof is snapshot-based, not continuous monitoring
- Violations are logged at time of validation check

## Custom Challenge Creation

Custom challenges allow you to:
1. Define challenge name and description
2. Add rules from available categories
3. Set rule parameters (level caps, restrictions, etc.)
4. Track compliance with same validation system

## Future Enhancements

- Real-time event hook integration for continuous monitoring
- Battle count tracking
- Gil transaction tracking
- Item source tracking
- Community challenge sharing
- Popular challenge mode library expansion

## Version

**Version**: 1.0.0  
**Phase**: 5 (Challenge & Advanced Tools)  
**Author**: FF6 Save Editor Plugin System
