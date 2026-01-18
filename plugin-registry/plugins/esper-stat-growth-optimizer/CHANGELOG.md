# Changelog - Esper Stat Growth Optimizer

## 1.0.0
- Added stat planning menu with per-stat esper ranking and target-level projection.
- Added comparison view for quick esper per-level bonuses.
- Added help dialog and safe, read-only menu loop.
- Aligned esper bonus data with models/game/esper_growth.go (flat Vigor/Speed/Stamina/Magic/Defense/MagicDef).
- Prefer runtime esper list and growth fields from plugin APIs when available; fallback to bundled table if absent.
