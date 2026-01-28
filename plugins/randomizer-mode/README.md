# Randomizer Mode

**Phase:** 6 (Gameplay-Altering Plugins)  
**Batch:** 3  
**Version:** 1.0.0  
**Status:** ✅ Complete  
**Lines of Code:** ~900  

## Overview

The **Randomizer Mode** plugin transforms FF6 into a procedurally modified experience. Every playthrough can be completely different—different stats, equipment, abilities, and more. Create endless unique playthroughs while maintaining reproducibility through seed-based randomization.

### Key Features

- **Seed-Based Randomization** - Share seeds for identical randomizations
- **3 Intensity Levels** - From mild tweaks to complete chaos
- **Multiple Randomization Systems** - Stats, equipment, abilities, espers, spells
- **Balance Validation** - Ensure fair gameplay difficulty
- **Seed Sharing** - Export/import seeds with other players

## Features

### Randomization Intensity Levels

**3 Configurable Profiles:**

1. **Mild** (20% variance)
   - Small stat variations
   - Recognizable equipment
   - Familiar gameplay feel
   - Good for casual randomization

2. **Moderate** (50% variance)
   - Significant stat changes
   - Different equipment combos
   - Still balanced
   - Good middle ground

3. **Chaos** (100% variance)
   - Extreme randomization
   - Unpredictable equipment
   - Wild difficulty spikes
   - For hardcore fans

### Randomization Systems

**Character Stats Randomization:**
- Health Points (HP): 30-100
- Magic Points (MP): 0-60
- Vigor: 15-35
- Speed: 15-35
- Stamina: 15-35
- Magic: 15-35
- Defense: 10-30
- M. Defense: 10-30
- Evade: 0-20%
- M. Evade: 0-20%

**Equipment Randomization:**
- Random weapons per character
- Random armor per character
- Random helmets per character
- Random relics (optional)
- Balanced power scaling

**Command Abilities:**
- Shuffle special commands among characters
- Random command order
- New ability combinations

**Esper Assignments:**
- Random esper distribution
- 2-3 espers per character
- Random spell learning

**Inventory Randomization:**
- Random consumable items
- Starting item variations
- Random quantities

**Spell Learning:**
- Shuffle which espers teach which spells
- Different spell availability
- New magic combinations

### Seed-Based System

- **Reproducible Randomization** - Same seed = same randomization
- **Share Seeds** - Distribute seed codes to other players
- **Seed Export** - Export as shareable string format
- **Seed History** - Track all applied seeds
- **Seed Validation** - Check seed compatibility

### Balance Features

- **Balance Validation** - Checks for extreme imbalances
- **Statistics Analysis** - Calculates average difficulty
- **Warning System** - Alerts if too hard/easy
- **Configuration Export** - Share randomization settings
- **Intensity Presets** - Pre-tuned difficulty profiles

## Installation

1. Extract `randomizer-mode` to plugins/
2. Restart FF6 Save Editor
3. Available in Phase 6 - Experimental
4. Backup save files before use!

## Usage

### Basic Randomization (New Seed)

1. Open a save file
2. Navigate to Randomizer Mode
3. Click "Apply New Randomization"
4. Choose intensity level:
   - Mild (20%)
   - Moderate (50%)
   - Chaos (100%)
5. Confirm application
6. Randomization is applied with new seed

### Using Existing Seed

1. Click "Use Existing Seed"
2. Enter seed code (from friend or export)
3. Select intensity level
4. Apply
5. Same randomization as original seed

### Export and Share

1. After randomization, click "Export Seed"
2. Get seed code: `FF6_RANDOMIZER_1234567890_moderate_20260116`
3. Share with other players
4. They can apply same seed/intensity

### View Randomization Details

1. After applying, view results
2. See all randomized stats
3. View equipment changes
4. Check esper assignments
5. Review spell changes

## Examples

### Example 1: Fresh Playthrough (New Seed, Mild)

Goal: Slightly different experience on replay

1. Open save file
2. Apply randomization (New Seed, Mild)
3. Characters have 10-20% stat variations
4. Slightly different equipment
5. Still feels like FF6

### Example 2: Shared Seed Challenge

Goal: Play the same randomized game as a friend

1. Friend provides seed: `FF6_RANDOMIZER_1609459200_moderate_20260116`
2. Click "Use Existing Seed"
3. Enter seed code
4. Apply (Moderate intensity)
5. Both get identical randomization
6. Can compare playthroughs

### Example 3: Chaos Run

Goal: Completely unpredictable experience

1. Apply randomization (New Seed, Chaos)
2. 100% stat variance
3. Unrecognizable equipment
4. Random abilities on all characters
5. Wild difficulty curve

### Example 4: Speedrun with Randomizer

Goal: Speedrun with randomized conditions

1. Export seed from one playthrough
2. Multiple speedrunners use same seed
3. Compare times on identical randomizations
4. Fair competitive comparison

### Example 5: Balance Testing

Goal: Verify randomization isn't too extreme

1. Apply randomization
2. View balance validation results
3. Check if too hard/easy
4. Adjust intensity if needed
5. Re-randomize if unbalanced

## Seed Codes

### Seed Format
```
FF6_RANDOMIZER_{seed}_{intensity}_{date}
```

Example:
```
FF6_RANDOMIZER_1609459200_moderate_20260116
```

### Components
- **seed:** Unix timestamp when generated
- **intensity:** mild/moderate/chaos
- **date:** YYYYMMDD creation date

### Sharing Seeds
- Short format: `1609459200_moderate`
- Long format: Full seed code
- Both work in "Use Existing Seed"

## Technical Details

### Randomization Algorithm

Uses seeded pseudo-random number generator (LCG):
- **Seed:** Unix timestamp or custom value
- **Range:** Configurable per system
- **Variance:** By intensity level (20%, 50%, 100%)
- **Distribution:** Linear scaling within ranges

### Stat Randomization

Each character stat is randomized:
1. Get stat range from CONFIG
2. Calculate midpoint and variance
3. Add random adjustment based on intensity
4. Clamp result to valid range
5. Apply to character

### Equipment Randomization

Random equipment selection:
1. For each character
2. Random weapon (1-20)
3. Random shield (0-10, 0=none)
4. Random helmet (1-15)
5. Random armor (1-15)
6. Random relics (0-20 each)

### Esper Assignment

14 characters, 26 espers:
1. Create esper pool
2. Shuffle esper order
3. Assign 2-3 espers per character
4. Rest remain unassigned
5. Next playthrough different assignment

### Balance Validation

Checks for:
- Very low average stats (<50) - too easy warning
- Very high average stats (>300) - too hard warning
- Character power distribution
- Equipment cost balance
- Difficulty curve smoothness

### API Requirements

- `character.setStats(id, stats)` - Set character stats
- `character.setEquipment(id, equipment)` - Set equipment
- `character.setAbility(id, slot, ability)` - Set command
- `character.assignEsper(id, esper_id)` - Assign esper
- `character.learnSpell(id, spell_id)` - Grant spell
- `inventory.setItemQuantity(id, qty)` - Set items
- `player.setGil(amount)` - Set starting currency

## Randomization Examples

### Mild Randomization (Terra)

**Original:** HP 39, MP 6, Vigor 15, Speed 12, Magic 18  
**Randomized:** HP 41, MP 5, Vigor 16, Speed 11, Magic 19  
**Change:** +2 HP, -1 MP, +1 Vigor, -1 Speed, +1 Magic

### Moderate Randomization (Terra)

**Original:** HP 39, MP 6, Vigor 15, Speed 12, Magic 18  
**Randomized:** HP 52, MP 12, Vigor 22, Speed 8, Magic 28  
**Change:** +13 HP, +6 MP, +7 Vigor, -4 Speed, +10 Magic

### Chaos Randomization (Terra)

**Original:** HP 39, MP 6, Vigor 15, Speed 12, Magic 18  
**Randomized:** HP 95, MP 48, Vigor 32, Speed 6, Magic 9  
**Change:** +56 HP, +42 MP, +17 Vigor, -6 Speed, -9 Magic

## Seed Reproducibility

### Key Property: Same Seed = Same Randomization

```
Seed 1609459200 + Moderate Intensity
= ALWAYS produces same randomization
= Different system, same randomization
= Shareable across players
```

### Reusing Seeds

After first randomization:
1. Note seed code
2. Apply to different save
3. Get identical randomization
4. Test different playstyles on same random config

## Community Features

### Seed Sharing

- Share seed codes on forums/Discord
- Create seed competitions
- Community randomizer leagues
- Leaderboards by seed

### Seed Collections

- Archive favorite seeds
- Document "good" seeds
- Rate seeds by difficulty
- Curate seed libraries

### Seed Discussion

- "Best seeds for..." threads
- Difficulty ratings
- Character combo analysis
- Strategy discussions

## Safety Considerations

### Backup Before Randomizing

Randomization modifies:
- All character stats
- All equipment
- All abilities
- All esper assignments
- Inventory items

Always backup first!

### Extreme Randomization

⚠️ Chaos intensity can create:
- Overpowered characters
- Completely weak characters
- Unbalanced difficulty
- Possibly unwinnable scenarios

Validation helps detect these.

## Troubleshooting

### Randomization Too Hard

Problem: Characters are dying constantly
Solution:
1. Try Mild intensity instead
2. Use different seed
3. Manually adjust stats
4. Use Instant Mastery to buff characters

### Randomization Too Easy

Problem: Bosses fall over instantly
Solution:
1. Try Chaos intensity
2. Challenge Mode restrictions
3. Self-imposed handicaps
4. Different seed

### Can't Reproduce Seed

Problem: Same seed gave different results
Solution:
1. Verify exact seed code
2. Check intensity matches
3. Ensure plugin version matches
4. Report if issue persists

## Known Limitations

### Reproducibility

- Same seed + same intensity = same results
- Different plugin versions may vary
- System clock doesn't affect randomization
- Cross-platform compatible

### Customization

- Can't customize per-stat ranges
- Can't weight certain stats
- Limited esper distribution control
- No ability to exclude specific items

### Balance

- Validation is heuristic-based
- May not catch all imbalances
- Extreme cases possible
- Manual review recommended

## FAQ

### Q: Can I use randomizer seeds for speedruns?
**A:** Check speedrun.com rules. Most require unmodified seeds. Randomizer may not be allowed.

### Q: How do I know if a seed is "good"?
**A:** Try it! Run balance validation, look at stat changes, test difficulty yourself.

### Q: Can I combine randomizer with other Phase 6 plugins?
**A:** Yes, Instant Mastery + Randomizer + Custom Pacing = total chaos mode!

### Q: What if my randomization is unwinnable?
**A:** Use Undo, try different seed, or use Instant Mastery to buff.

### Q: How many unique seeds exist?
**A:** Theoretically infinite (Unix timestamps), practically millions of useful seeds.

### Q: Can I randomize just stats?
**A:** Future version - currently randomizes all systems together.

## Version History

### v1.0.0 (January 16, 2026) - Initial Release
- Seed-based randomization
- 3 intensity levels
- Character stats randomization
- Equipment randomization
- Command ability shuffling
- Esper assignment randomization
- Spell learning shuffling
- Balance validation
- Seed export/import
- Complete documentation

## Recommended Seed Presets

### For Casual Play
- Intensity: Mild
- Gives recognizable variations
- Familiar equipment
- Smooth difficulty

### For Challenging Play
- Intensity: Moderate
- Significant changes
- Still balanced
- Good difficulty curve

### For Hardcore Fans
- Intensity: Chaos
- Extreme randomization
- High variance
- Very unpredictable

### For Competition
- Intensity: Moderate
- Use shared seed
- Same experience for all
- Fair comparison

## Credits

**Phase:** 6 Batch 3  
**Developer:** FF6 Save Editor Team  
**Inspiration:** Randomizer community interests

---

**Status:** Stable v1.0.0  
**Last Updated:** January 16, 2026
