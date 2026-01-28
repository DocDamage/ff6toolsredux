# Poverty Mode

**Version:** 1.0.0  
**Category:** Experimental / Challenge  
**Phase:** 6, Batch 4  
**Author:** FF6 Editor Plugin System

---

## Overview

The **Poverty Mode** plugin creates an extreme challenge by removing all or most resources from your Final Fantasy VI save file. Experience true survival gameplay where every item, every gil piece matters. This plugin is designed for hardcore players seeking the ultimate challenge: beating FF6 with zero or near-zero resources.

### Key Features

- ✅ **Zero Gil Mode** - Start with absolutely no money
- ✅ **Remove Consumables** - No potions, ethers, or consumable items
- ✅ **Equipment-Only Survival** - Keep only what's equipped
- ✅ **5 Difficulty Presets** - From Light to Extreme poverty
- ✅ **Custom Poverty Config** - Create your own challenge rules
- ✅ **Challenge Tracking** - Track found items and compliance
- ✅ **Backup & Restore** - Undo if challenge becomes unwinnable
- ✅ **Challenge Proof Export** - Export completion evidence

⚠️ **WARNING:** This plugin is **EXTREMELY CHALLENGING** and can make the game unwinnable if used incorrectly. **ALWAYS BACKUP YOUR SAVE FILE** before use. This is NOT recommended for first playthroughs.

---

## Installation

1. Copy the `poverty-mode` folder to your FF6 Editor `plugins/` directory
2. Restart FF6 Editor or reload plugins
3. The plugin will appear in the **Experimental** category

**Requirements:**
- FF6 Editor v3.4.0 or higher
- Valid FF6 save file (recommend mid-game or later)
- External save backup (REQUIRED)
- Permissions: `read_save`, `write_save`, `ui_display`

---

## Poverty Levels (Presets)

### Level 1: Light Poverty
**Difficulty:** ★☆☆☆☆  
**Gil:** 1,000  
**Items:** Max 10 of each  
**Equipment:** Allowed  

Perfect for first-time poverty challenges. You can carry limited items and have some Gil for emergencies.

```lua
PovertyMode.applyPovertyLevel("LIGHT")
```

---

### Level 2: Moderate Poverty
**Difficulty:** ★★☆☆☆  
**Gil:** 500  
**Items:** Max 5 of each  
**Equipment:** Allowed  

A step up in difficulty. Item limits are strict but manageable.

```lua
PovertyMode.applyPovertyLevel("MODERATE")
```

---

### Level 3: Strict Poverty
**Difficulty:** ★★★☆☆  
**Gil:** 100  
**Items:** No consumables  
**Equipment:** Allowed (found only)  

No consumables allowed. You must rely entirely on equipment and character abilities.

```lua
PovertyMode.applyPovertyLevel("STRICT")
```

---

### Level 4: Hardcore Poverty
**Difficulty:** ★★★★☆  
**Gil:** 0  
**Items:** No consumables  
**Equipment:** Allowed (equipped only)  

Zero Gil. No consumables. Only equipped items remain. Extremely challenging.

```lua
PovertyMode.applyPovertyLevel("HARDCORE")
```

---

### Level 5: Extreme Poverty
**Difficulty:** ★★★★★  
**Gil:** 0  
**Items:** None  
**Equipment:** Equipped only  

The ultimate challenge. Zero resources, zero consumables, zero unequipped items. Only your equipped gear and skills.

```lua
PovertyMode.applyPovertyLevel("EXTREME")
```

---

## Custom Poverty Configuration

Create your own challenge rules:

```lua
PovertyMode.applyCustomPoverty({
    gil = 250,                          -- Set Gil amount
    max_item_quantity = 3,              -- Max items per slot
    remove_consumables = true,          -- Remove consumables?
    remove_unequipped_equipment = false, -- Remove unequipped gear?
})
```

**Examples:**

**No-Shop Challenge (Gil zero, items allowed):**
```lua
PovertyMode.applyCustomPoverty({
    gil = 0,
    max_item_quantity = 99,
    remove_consumables = false,
    remove_unequipped_equipment = false,
})
```

**Magic-Only Challenge (remove physical equipment):**
```lua
PovertyMode.applyCustomPoverty({
    gil = 0,
    max_item_quantity = 10,
    remove_consumables = true,
    remove_unequipped_equipment = true,
})
```

**Speedrun Poverty (minimal resources, fast start):**
```lua
PovertyMode.applyCustomPoverty({
    gil = 100,
    max_item_quantity = 5,
    remove_consumables = false,
    remove_unequipped_equipment = true,
})
```

---

## Challenge Tracking

### Track Found Items

When you find items during poverty mode, track them:

```lua
PovertyMode.trackFoundItem(itemId, quantity, "treasure chest")
PovertyMode.trackFoundItem(itemId, quantity, "boss drop")
PovertyMode.trackFoundItem(itemId, quantity, "event reward")
```

This builds a log of all resources acquired during the challenge.

---

### Check Compliance

Verify your save still complies with poverty rules:

```lua
local compliant, message = PovertyMode.checkCompliance()

if compliant then
    print("✅ Challenge still valid!")
else
    print("❌ Violations detected:")
    for _, violation in ipairs(message) do
        print("  - " .. violation)
    end
end
```

**Violations Include:**
- Gil exceeding allowed amount
- Consumables when they're forbidden
- Item quantities exceeding limits

---

### Export Challenge Proof

Generate proof of challenge completion:

```lua
local proof = PovertyMode.exportChallengeProof()
```

**Proof Contains:**
```lua
{
    poverty_level = "HARDCORE",
    challenge_start_time = 1704672000,
    challenge_duration = 14400,  -- 4 hours
    found_items_count = 23,
    operations_count = 156,
    compliance_status = true,
    current_gil = 0,
}
```

Share this proof with the community to verify your challenge run!

---

## Restoration

### Restore from Backup

Undo poverty mode if the challenge becomes unwinnable:

```lua
PovertyMode.restoreBackup()
```

**Restores:**
- Original Gil amount
- All items and quantities
- Pre-poverty state

**Important:** Backup is created automatically when applying poverty mode.

---

### Check Status

View current poverty mode status:

```lua
local status = PovertyMode.getStatus()
```

**Returns:**
```lua
{
    enabled = true,
    poverty_level = "HARDCORE",
    challenge_duration = 3600,  -- seconds
    has_backup = true,
    found_items_count = 12,
    compliance_violations = 0,
}
```

---

## Strategies for Poverty Mode

### Surviving Extreme Poverty

1. **Prioritize Equipment** - Every piece of gear matters
2. **Learn Attack Magic Early** - Free damage without consumables
3. **Master Defensive Spells** - Protect, Shell, Regen
4. **Use Steal Command** - Only source of items
5. **Avoid Optional Battles** - Conserve resources
6. **Plan Boss Strategies** - No retries without items
7. **Exploit Enemy Weaknesses** - Maximize damage efficiency

### Recommended Party Composition

**Physical Focus:**
- Sabin (Blitzes require no items)
- Edgar (Tools are unlimited)
- Cyan (Bushido uses no items)
- Locke (Steal for item acquisition)

**Magic Focus:**
- Terra (Natural magic)
- Celes (Runic absorbs magic)
- Strago (Lores are reusable)
- Relm (Sketch for utility)

### Critical Spells

- **Cure Series** - Free healing
- **Life/Life 2** - Revival without Phoenix Downs
- **Osmose** - MP recovery
- **Rasp** - MP damage to enemies
- **Regen** - Passive healing
- **Reflect** - Bounce enemy magic

---

## Use Cases

### 1. No-Shop Challenge

**Goal:** Beat the game without ever using shops.

```lua
PovertyMode.applyCustomPoverty({
    gil = 0,
    remove_consumables = false,  -- Keep found items
    remove_unequipped_equipment = false,
})
```

**Rules:** Zero Gil, cannot buy anything, must rely on found items only.

---

### 2. Hardcore Nuzlocke-Style

**Goal:** Extreme difficulty with permadeath elements.

```lua
PovertyMode.applyPovertyLevel("HARDCORE")
```

**Rules:** Zero Gil, no consumables, equipped items only. Combine with permadeath tracking for ultimate challenge.

---

### 3. Equipment-Only Run

**Goal:** Rely solely on equipment stats and abilities.

```lua
PovertyMode.applyPovertyLevel("EXTREME")
```

**Rules:** No items at all. Equipment and character abilities only.

---

### 4. Speedrun Practice

**Goal:** Optimize routing with limited resources.

```lua
PovertyMode.applyPovertyLevel("MODERATE")
```

**Rules:** Limited items force optimal routing and boss strategies.

---

### 5. Community Challenge

**Goal:** Complete poverty run with proof for community recognition.

```lua
PovertyMode.applyPovertyLevel("HARDCORE")
-- Play through game
local proof = PovertyMode.exportChallengeProof()
-- Share proof with community
```

---

## Technical Details

### Resource Removal Process

1. **Backup Creation** - Save current state
2. **Confirmation** - Verify user intent
3. **Gil Removal** - Set Gil to 0 or limit
4. **Consumable Removal** - Remove items ID 0-99
5. **Equipment Handling** - Remove unequipped items (ID 100-254)
6. **Item Limiting** - Cap quantities to max allowed
7. **Logging** - Record all operations

### Item Categories

- **Consumables (ID 0-99):** Potions, Ethers, Tonic, etc.
- **Equipment (ID 100-254):** Weapons, Armor, Relics

### Equipped Item Detection

Plugin scans all 14 characters for:
- Weapon
- Shield
- Helmet
- Armor
- Relic 1
- Relic 2

These items are preserved in poverty mode (except Extreme).

---

## Troubleshooting

### Issue: "No backup available to restore"

**Cause:** Trying to restore without applying poverty mode first.

**Solution:** Backup is created when poverty mode is applied. Cannot restore if never applied.

---

### Issue: Game becomes unwinnable

**Cause:** Applied extreme poverty too early in the game.

**Solution:** Use `restoreBackup()` to undo. Recommend applying poverty mode mid-game (after Narshe/Figaro).

---

### Issue: Cannot defeat bosses without items

**Cause:** Poverty mode removes healing/revival items.

**Solution:** Learn Cure/Life spells. Use defensive buffs. Optimize equipment. Consider lower poverty level.

---

### Issue: Compliance violations detected

**Cause:** Acquired Gil or items beyond poverty rules.

**Solution:** This is expected from found items. Tracking helps document your challenge. Violations don't invalidate the run unless you intentionally break rules.

---

### Issue: Cannot buy anything (no shops)

**Cause:** Gil is zero in poverty mode.

**Solution:** This is intended. Poverty mode is a no-purchase challenge. Find items or restore backup.

---

## Frequently Asked Questions

### Q: Is poverty mode reversible?
**A:** Yes! Use `restoreBackup()` to undo all changes. Backup is created automatically.

### Q: Can I use poverty mode on a new game?
**A:** Technically yes, but EXTREMELY difficult. Recommend starting at least after Narshe or mid-game.

### Q: What happens if I find items during poverty mode?
**A:** You can keep them within your poverty rules. Use `trackFoundItem()` to log acquisitions.

### Q: Can I combine poverty mode with other plugins?
**A:** Yes! Combine with Hard Mode Creator, Permadeath Tracker, or Challenge Timer for extreme difficulty.

### Q: Does poverty mode affect story progression?
**A:** No. It only removes resources. Story events proceed normally.

### Q: Can I change poverty level mid-challenge?
**A:** Not recommended. Creates inconsistency. Finish challenge at chosen level or restore backup.

### Q: Is there a "poverty mode leaderboard"?
**A:** Not officially, but export challenge proof to share with the community. Compare completion times and compliance.

### Q: What's the hardest poverty level?
**A:** Extreme Poverty (zero everything). Many bosses become nearly impossible without consumables.

### Q: Can I adjust poverty rules after starting?
**A:** No. Rules are set when poverty mode is applied. Restore backup to start fresh with different rules.

### Q: Does poverty mode work with challenge runs?
**A:** Yes! It's designed specifically for challenge runs: Nuzlocke, solo character, speedrun, etc.

---

## Known Limitations

1. **Item Categories** - Item ID ranges approximate; some edge cases may exist
2. **Equipped Detection** - Only detects currently equipped items per character
3. **Story Items** - Some key items cannot be removed (game requirement)
4. **Shop Access** - Shops remain accessible but Gil is zero (cannot buy)
5. **Found Items** - Manual tracking required; not automatic detection

---

## Safety & Warnings

⚠️ **ALWAYS BACK UP YOUR SAVE FILE EXTERNALLY** before using this plugin.

⚠️ **CAN MAKE GAME UNWINNABLE** if applied at wrong game stage.

⚠️ **NOT RECOMMENDED FOR:**
- First playthroughs
- Early-game saves (pre-Narshe)
- Players unfamiliar with FF6 mechanics
- Casual players

✅ **RECOMMENDED FOR:**
- Experienced FF6 players
- Challenge run enthusiasts
- Hardcore difficulty seekers
- Community challenge participants
- Speedrunners (practice mode)

---

## Changelog

### Version 1.0.0 (January 16, 2026)
- Initial release
- 5 poverty level presets
- Custom poverty configuration
- Challenge tracking and compliance
- Backup and restoration
- Challenge proof export
- Found items logging

---

## Credits

**Developed by:** FF6 Editor Plugin System  
**Plugin Category:** Experimental / Challenge  
**Phase 6, Batch 4** - Part of the 44-plugin expansion

**Special Thanks:**
- FF6 challenge run community
- No-shop run pioneers
- Poverty mode playtesters

---

## License

This plugin is part of the FF6 Editor Plugin System and follows the main project's license terms.

---

**Good luck, and may the crystals guide you! 💎🔥**
