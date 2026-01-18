# Character Ability Swap Changelog

## [1.0.0] - 2026-01-16

### Added
- Initial release of Character Ability Swap plugin
- **Command ability swapping system** for all 14 characters
- **18 command ability library**:
  1. Fight - Basic attack (all characters)
  2. Magic - Cast spells (Terra/Celes)
  3. Morph - Esper form (Terra)
  4. Steal - Steal items (Locke)
  5. Capture - Colosseum capture (Locke)
  6. SwdTech - Bushido techniques (Cyan)
  7. Throw - Throw weapons (Shadow)
  8. Tools - Mechanical tools (Edgar)
  9. Blitz - Martial arts (Sabin)
  10. Runic - Absorb magic (Celes)
  11. Lore - Enemy Lores (Strago)
  12. Sketch - Mimic enemy (Relm)
  13. Control - Control enemy (Relm)
  14. Slot - Slot attacks (Setzer)
  15. Rage - Enemy Rages (Gau)
  16. Leap - Jump to Veldt (Gau)
  17. Mimic - Repeat action (Gogo)
  18. Dance - Perform Dances (Mog)
- **Character editing interface**:
  - Select any character (14 total)
  - Edit any of 4 command slots
  - Choose from complete ability library
  - Set empty slots
- **4 preset builds**:
  1. **All-Magic Party** - Every character has Magic command
  2. **All-Physical Party** - Focus on physical abilities (Blitz/Throw/SwdTech/Tools)
  3. **Utility Focused** - Emphasize Steal/Control/support abilities
  4. **Chaos Mode** - Random unusual combinations (Rage/Dance/Lore mixes)
- **View current setup** - Display all 14 characters with their 4 command slots
- **Reset to defaults** - Restore vanilla ability setups
- **Save current setup** - Store custom configurations with names
- **Load saved setup** - Recall previously saved configurations
- **Export configuration** - Generate text file with complete setup
- Command descriptions and original character associations
- Safe API wrappers with pcall error handling
- Persistent storage using Lua table serialization

### Default Command Setups
- Terra: Fight, Magic, Morph, (Item)
- Locke: Fight, Steal, Capture, (Item)
- Cyan: Fight, SwdTech, (Item), (empty)
- Shadow: Fight, Throw, (Item), (empty)
- Edgar: Fight, Tools, (Item), (empty)
- Sabin: Fight, Blitz, (Item), (empty)
- Celes: Fight, Magic, Runic, (Item)
- Strago: Fight, Magic, Lore, (Item)
- Relm: Fight, Sketch, Control, (Item)
- Setzer: Fight, Slot, (Item), (empty)
- Mog: Fight, Magic, Dance, (Item)
- Gau: Fight, Rage, Leap, (Item)
- Gogo: Fight, Mimic, (Item), (empty) - can customize all slots
- Umaro: Fight only (auto-battle character)

### Features
- Slot-by-slot command customization
- Complete ability library with descriptions
- Preset build application
- Configuration save/load system
- Export functionality
- Reset protection

### Technical Details
- **Lines of Code**: ~800
- **Permissions**: read_save, write_save, ui_display
- **Dependencies**: None (command library hardcoded)
- **Data Format**: Lua table literal serialization
- **Storage Files**:
  - `ability_swap/configurations.lua` - active and saved setups
  - `ability_swap/export_*.txt` - exported configurations

### Configuration Structure
```lua
{
  active_setup = {
    [character_name] = {slot1_id, slot2_id, slot3_id, slot4_id}
  },
  saved_setups = {
    {name = "Setup Name", setup = {...}, saved = timestamp}
  }
}
```

### Experimental Warning
⚠️ This plugin requires NEW APIs for actual in-game ability modification:
- `character.get_command_abilities(char_id)` → read commands (NEW API)
- `character.set_command_ability(char_id, slot, ability_id)` → write commands (NEW API)

Without these APIs:
- ✅ Configuration tracking and planning works
- ✅ Export and documentation works
- ❌ Direct in-game ability modification not possible

### Known Limitations
- **API Dependency**: Requires command ability APIs for actual modification
- Without API, acts as planning/documentation tool
- Some abilities character-specific (Morph, Leap)
- Balance considerations: certain combinations may be overpowered
- Character identity alteration
- Save compatibility concerns

### Notes
- Phase 6, Plugin 1 of 18
- Part of Gameplay-Altering Plugins phase
- Experimental feature with fundamental gameplay impact
- Test on backup saves recommended
- Community preset sharing potential
- Slot 4 typically reserved for Item command (not modifiable in vanilla game)

### Example Use Cases
1. **All-Mage Party**: Give everyone Magic for spell-heavy runs
2. **Hybrid Builds**: Mix physical/magical abilities on each character
3. **Challenge Runs**: Restrict to specific ability types
4. **Experimental Gameplay**: Try unusual combinations (Terra with Rage, Gau with Magic)
5. **Balance Testing**: Test ability synergies and combinations

### Future API Requirements
- `character.get_command_abilities(char_id)` - Read current command setup
- `character.set_command_ability(char_id, slot, ability_id)` - Modify command slots
- Command ability validation (check if ability works with character)
- Battle system integration for custom ability execution
