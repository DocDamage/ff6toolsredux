# World of Balance Checklist

Completion tracker focused on everything you can (and should) finish before the World of Ruin transition: recruits, espers, missable items, rages, dances, and key event flags. The plugin is intentionally **read-only** and errs on the side of caution—if an API is missing, it will flag the item as "Unknown" and provide guidance instead of guessing.

> Version: 1.0.0  \
> Permissions: `read_save`, `ui_display`  \
> Scope: World of Balance (with optional context for WoR-only recruits so you know what cannot be done yet)

---

## Why This Plugin Exists

The World of Balance is the biggest concentration of missables in FF6. Once you depart the Floating Continent, any unchecked box turns into a permanent regret: lost rages, unlearned dances, missed relics, and event rewards you cannot redo. This checklist keeps everything visible in one place, driven by your save data where possible and clearly marked as reference-only when the current API cannot verify a status.

Design goals:
- **Prevent regrets**: Early warning for point-of-no-return content.
- **Stay honest**: If the editor cannot read a flag, it tells you instead of faking a status.
- **Fast to use**: One menu, seven options, minimal typing.
- **Version aware**: Notes for differences across SNES/GBA/Pixel Remaster.
- **Safe**: Read-only; no save mutations.

---

## Feature Checklist

- **Summary dashboard**: World state, characters recruited, espers acquired, and a floating continent warning if you are still in WoB.
- **Character checklist**: 14-slot roster with WoR-only characters explicitly labeled to avoid confusion.
- **Esper checklist**: Expected WoB-accessible espers with dynamic ownership detection; marks acquired vs. missing.
- **Missable content view**: Items, banquet rewards, rages, dances, and key events in one list with quick tips.
- **Event flags**: Reads known flags (opera, banquet score, factory clear, world state) when the API exposes them.
- **Export view**: Copyable text snapshot for planning/streaming notes.
- **Help view**: Quick reminders of what each option does and how unknown statuses are handled.

---

## Quick Start

1. Open a WoB save (or any save if you only want the reference data).
2. Run the plugin and choose `1) Summary` for a sanity check.
3. Visit `2) Characters` and confirm all WoB recruits are in your roster.
4. Visit `3) Espers` to see which WoB magicite you have secured.
5. Visit `4) Missables` for items/rages/dances and make a short to-do list.
6. Visit `5) Event Flags` if you want to verify opera/banquet/factory completion.
7. Use `6) Export` to copy the state into a note or tracker of your choice.

Tip: If a status shows **Unknown**, it means the API needed to verify it was not available. Use the note next to the item to double-check in-game.

---

## Data Sources and Safety

- **Roster**: Pulled from `GetCharacter(index)` for slots 0-15. A character is considered recruited if any of `Enabled`, `enabled`, or `Available` is truthy.
- **Espers**: Pulled from `GetEspers()`. A magicite is considered owned if any of `Owned`, `owned`, or `Acquired` is truthy.
- **World state**: First tries `GetWorldState()`, then `GetFlag("world_state")`. If both fail, it assumes **World of Balance** by default (configurable).
- **Event flags**: Reads named flags directly via `GetFlag(flagName)`. Missing flags render as `Unknown`.
- **Missables**: Reference-only. No reliable flag coverage exists yet for many items/dances/rages; the plugin shows tips instead of statuses.
- **Write operations**: None. The plugin never calls write APIs.

Fallback behavior is always explicit: Unknown statuses are surfaced plainly with short guidance so you can verify in-game.

---

## Menu Options in Detail

### 1) Summary
- Shows world detection result (WoB, WoR, or Unknown).
- Character completion: count and percentage.
- Esper completion: count and percentage for expected WoB magicite.
- Missable note reminding you to open the dedicated view.
- Floating Continent warning if you are still flagged as WoB (or world unknown).

### 2) Characters
- Lists 14 roster entries. Umaro and Gogo are labeled **(WoR)** so you do not chase them prematurely.
- Uses your save data to mark recruited vs. missing.
- Shows a short note for each recruit (where/how they join).

### 3) Espers
- Uses the **expected WoB list**: Ramuh, Kirin, Siren, Stray, Ifrit, Shiva, Unicorn, Maduin, Carbuncle, Bismarck, Phantom, Seraphim, Golem, ZoneSeek, Shoat, Tritoch, Fenrir, Terrato, Starlet.
- Marks an esper as acquired if it appears in your `GetEspers()` list.
- Displays totals and percentage for quick confidence checking.
- If your version differs (e.g., Pixel Remaster or GBA additions), the reference list still works, but feel free to extend it in `ESPER_TARGETS` near the top of `plugin.lua`.

### 4) Missables
- Groups into four buckets: **Missable Items/Rewards**, **Rages**, **Dances**, **Key Events**.
- Every entry includes a location and a one-line tip.
- Status shows `?` because the current APIs do not expose those flags; treat it as a reminder list.
- Pagination keeps the view readable (20 entries per page by default).

### 5) Event Flags
- Reads specific flags (opera completion, banquet score, floating continent reach, WoR transition, factory clear).
- If a flag is unavailable, the status shows `Unknown`.
- Use this to confirm you have satisfied key story beats before moving on.

### 6) Export
- Produces a copy-friendly text block with world state, character checks, esper checks, and the missable list.
- Handy for streaming, speedrun notes, or sharing with friends.

### 7) Help
- One-page recap of what each option does and how unknowns are handled.

---

## Configuration

Edit the `config` table at the top of `plugin.lua` to tune behavior:

- `assumeWorldOfBalanceWhenUnknown` (default `true`): When world detection fails, treat the save as WoB and keep warnings on. Set to `false` if you prefer neutrality.
- `showCompletionPercentages` (default `true`): Toggle percentage display in summary and category views.
- `warnPointOfNoReturn` (default `true`): Shows a Floating Continent warning whenever world is WoB or Unknown.
- `showMissableGuides` (default `true`): If set to `false`, missable entries will not show the extra tip lines.
- `showEventFlagDetails` (default `true`): Include flag descriptions in the Event Flags view.
- `maxListPerPage` (default `20`): Page size for paginated views.

All options are read at runtime; no rebuild is needed—just reload the plugin after saving the file.

---

## World Detection Logic

1. Try `GetWorldState()`.
2. If nil, try `GetFlag("world_state")`.
3. If both fail, fall back to `World of Balance` when `assumeWorldOfBalanceWhenUnknown` is true; otherwise show `Unknown`.

Why the fallback? The checklist is designed to be conservative. It is safer to warn about missables when unsure than to stay silent and let you proceed into WoR with missing content.

---

## Character Checklist Notes

- **Shadow**: The plugin cannot verify whether you waited for him on the Floating Continent; it only checks whether he exists in your roster. Treat the roster check as necessary but not sufficient.
- **Gau**: Requires visiting the Veldt and using Leap. If he is not in the roster, do it before boarding the Blackjack for the Floating Continent.
- **Mog**: Save him in Narshe. This is your chance to learn early dances; see Missables.
- **WoR-only recruits (Umaro, Gogo)**: Shown as informational entries. They will remain unchecked in WoB and that is expected.

---

## Esper Checklist Notes

- The reference list covers common WoB magicite. Version differences:
  - **SNES/PSX**: Matches the list above; Fenrir/Starlet are WoR in vanilla but shown for visibility (treated as missing until acquired).
  - **GBA/Pixel Remaster**: Leviathan/Cactuar/Gilgamesh are WoR and not in the default list. Add them if you want them tracked before WoR.
- If an esper appears in your save data with a slightly different name (localization), the plugin will still match if the name contains the target substring.
- If `GetEspers()` is unavailable, all espers will show as unchecked; use the notes to verify manually.

---

## Missable Content Coverage

### Items and Rewards
- **Charm Bangle (South Figaro path)**: Steal it while wearing the merchant disguise.
- **Tintinabar**: Speak to Arvis after Banon events; easy to forget.
- **Memento Ring**: In Relm's house during the Thamasa fire sequence; miss it and you cannot re-enter.
- **Banquet Rewards**: Optimal answers yield Hero Ring and Charm Bangle. Partial answers reduce the reward tier.

### Rages to Grab Early
- **Tyranosaur**: Hard to find later; great for power leveling in WoR but get it now if possible.
- **Intangir**: Triangle Island. Disappears in WoR.
- **Brontaur**: Solid early Rage from Narshe cliffs.

### Dances to Learn Early
- **Wind Rhapsody** (grass), **Forest Suite** (forest), **Desert Aria** (desert), **Snowman Jazz** (snow field). Learn them before world-state changes, especially Snowman Jazz while the Narshe snowfield is accessible.

### Key Events
- **Opera Success**: The flag should read completed; if Unknown, make sure the sequence is done before moving the plot.
- **Imperial Banquet**: Optimal dialog grants better rewards and helps trust.
- **Magitek Factory Clear**: Ensures Ifrit/Shiva and the escape are done before proceeding.

> Note: Because the current API does not expose chest flags or dance/rage learning flags, these entries are reminder-only. The plugin deliberately labels statuses with `?` instead of guessing.

---

## Event Flags

The plugin attempts to read a handful of high-value flags:
- `event_opera_complete`: Success state for the opera sequence.
- `event_banquet_score`: Numerical or boolean indicator of banquet performance.
- `event_floating_continent`: Whether the Floating Continent has been reached.
- `world_state`: WoB/WoR indicator.
- `event_magitek_factory`: Completion flag for the factory arc.

If a flag is missing, the status is `Unknown`. Use in-game confirmation in that case.

---

## Export Workflow

`6) Export` produces a plaintext summary suitable for:
- Sharing progress with co-op partners or friends.
- Posting in Discord for help.
- Keeping alongside a speedrun route.
- Archiving your state before attempting a challenge run.

Copy the dialog contents directly from the editor.

---

## Limitations and Known Gaps

- **Chest flags**: No API coverage yet; treasure tracking will arrive in Phase 2.2.
- **Bestiary/encounter flags**: Not used here; slated for Phase 2.3.
- **Rage/Dance learning flags**: Not exposed; the plugin lists tips only.
- **Shadow survival**: Cannot detect whether you waited on the Floating Continent; verify manually.
- **Esper list variance**: Localization/version differences may require you to adjust `ESPER_TARGETS`.
- **World detection**: Falls back to WoB by design when unknown; disable in config if you prefer neutrality.

---

## Troubleshooting

- **Everything shows Unknown**: The save may not be loaded, or the editor build may lack the needed APIs. Load a save and retry. If still unknown, check the API changelog for your build.
- **Espers not detected**: Verify `GetEspers()` exists in your build. If the API uses different property names, update `esperOwned()` in the plugin (search for `esperOwned` near the top).
- **Characters mislabelled**: If your build uses different fields than `Enabled/enabled/Available`, add the correct field to `isCharacterRecruited()`.
- **Wrong world state**: Toggle `assumeWorldOfBalanceWhenUnknown` to `false` so you can manually judge when detection fails.
- **Pagination too small/large**: Change `maxListPerPage` to your preferred page size.

---

## How to Extend or Customize

- Add or remove espers from `ESPER_TARGETS` to match your version or house rules.
- Add new missable items/rages/dances to the `MISSABLES` table with a location and note.
- Add more event flags to the `EVENT_FLAGS` table; the plugin will render them automatically.
- Adjust config options to control warnings, percentages, and pagination.

All data lives near the top of `plugin.lua` for ease of editing.

---

## Roadmap (Phase 2 and Beyond)

- **Treasure Chest Tracker (Phase 2.2)**: Add real chest flag checks when the API ships.
- **Bestiary (Phase 2.3)**: Use encounter flags to reduce Unknown statuses for rages/dances.
- **Colosseum Guide (Phase 2.4)**: Cross-link with WoB preparation steps for optional gearing.
- **Write-mode variant** (future, opt-in): Ability to mark checklist items manually if an API is permanently unavailable.

---

## FAQ

**Q: Why are Umaro and Gogo listed if they cannot be recruited in WoB?**  
A: They are shown to prevent confusion. They will remain unchecked in WoB and that is expected.

**Q: My esper names differ (localization). Will they match?**  
A: Matching is case-insensitive and substring-based. If still unmatched, add the localized name to `ESPER_TARGETS`.

**Q: Can I hide the Floating Continent warning?**  
A: Set `warnPointOfNoReturn = false` in `config`.

**Q: Does this plugin modify my save?**  
A: No. It is entirely read-only.

**Q: Why do some entries show `?`?**  
A: The current API does not expose reliable flags for those items. The plugin refuses to guess to avoid misleading you.

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and planned enhancements.

---

## Support and Feedback

- If you spot a mismatch between your version and the reference lists, update the tables and open a PR or issue.
- If a new API ships for chest flags or event flags, extend the plugin and bump the version.
- Bug reports: include your editor version, platform, and whether `GetWorldState`, `GetEspers`, and `GetFlag` are available.

---

## Release Checklist (for maintainers)

- [ ] Verify on a WoB save: roster detection works and no crashes occur when an API is missing.
- [ ] Validate `GetEspers()` naming: ensure Owned/Acquired detection matches your build.
- [ ] Run through each menu option to confirm pagination and dialogs render correctly.
- [ ] Update `CHANGELOG.md` with any code or data changes.
- [ ] Tag release `v1.0.0` (or later) in the registry when stable.

Happy hunting—leave no missable behind!
