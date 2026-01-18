# ROM Sprite Loading - Quick Start Guide

## Overview

The FF6 Save Editor now supports loading authentic character sprites directly from your Final Fantasy VI ROM file. This enables viewing and exporting actual character animations instead of placeholder graphics.

## Prerequisites

- A valid Final Fantasy VI ROM file (USA, Japan, or European version)
- Supported formats: `.smc`, `.sfc`, `.fig`
- Expected ROM size: 3-4 MB

## Setup Instructions

### 1. Configure ROM Path

1. Launch the FF6 Save Editor (`ffvi_editor.exe`)
2. Go to **Tools → Configure ROM...**
3. Click **Browse...** to select your ROM file
4. Click **Save**
5. **Restart the application** for changes to take effect

### 2. Load a Save File

1. Go to **File → Open**
2. Navigate to your save file (e.g., `save_data/76561198072182150/`)
3. Select and open a `.ff6` save file

### 3. View Character Sprites

Once your ROM is configured and a save file is loaded, you can:

#### Animation Player
- Go to **Tools → Animation Player...**
- View the first character's walking animation
- The sprite will be loaded from the ROM automatically
- Adjust playback speed and loop settings

#### Frame Editor
- Go to **Tools → Frame Editor...**
- View and edit individual sprite frames
- Each frame is 16×24 pixels from the ROM

#### Export Animation
- Go to **Tools → Export Animation...**
- Export character sprites as GIF or sprite sheets
- Authentic ROM graphics will be exported

#### Batch Export All Characters
- Go to **Tools → Export All Animations...**
- Export all party members' sprites at once
- Select a destination folder

## Supported Characters

The ROM loader supports all 14 playable characters:
- Terra (0), Locke (1), Cyan (2), Shadow (3)
- Edgar (4), Sabin (5), Celes (6), Strago (7)
- Relm (8), Setzer (9), Mog (10), Gau (11)
- Gogo (12), Umaro (13)

## Technical Details

### ROM Detection
- The editor validates ROM files using SHA256 checksums
- Supports ROMs with or without SMC headers (512-byte headers are automatically detected)
- Known ROM versions: USA 1.0, Japan 1.0/1.1, European release

### Sprite Format
- **Field sprites**: 16×24 pixels, 3 frames, 4bpp (4 bits per pixel)
- **Compression**: Custom LZ77 variant used by FF6
- **Palettes**: 16 colors per sprite, RGB555 format (15-bit color)

### ROM Offsets
Character sprites are located at fixed offsets in the ROM:
- **Field sprites**: 0x150000–0x15D000
- **Battle sprites**: 0x268000–0x282000 (not yet implemented)

### Fallback Behavior
If no ROM is configured or the ROM fails to load:
- The editor falls back to placeholder sprites (blank 16×24 sprites)
- All features remain functional, but sprites will appear empty
- You can configure the ROM at any time without losing your edits

## Troubleshooting

### ROM Not Loading
- **Check the file path**: Ensure the ROM path in settings is correct
- **Verify ROM size**: Should be 3-4 MB
- **Check file extension**: Must be `.smc`, `.sfc`, or `.fig`
- **Try another ROM**: Different ROM versions have different checksums

### Sprites Still Blank
- **Restart required**: ROM loading only happens at startup
- **Check character ID**: Only IDs 0-13 are supported
- **Verify ROM validity**: Check console output for ROM loading errors

### Console Output
The editor prints ROM loading status to the console:
```
ROM loaded successfully from: C:\path\to\ff6.smc
Loaded sprite for character 0 from ROM
```

If you see warnings:
```
Warning: Failed to load ROM from C:\path\to\ff6.smc: invalid ROM file
Warning: Failed to load sprite for character 0 from ROM: decompression error
```

This indicates a problem with the ROM file or sprite extraction.

## Known Limitations

1. **Battle sprites not implemented**: Only field sprites (16×24) are currently supported. Battle sprites (32×32) will be added in a future update.

2. **Palette loading**: Currently uses placeholder palettes. Full palette extraction from ROM will be implemented later.

3. **ROM path is global**: The ROM path is shared across all save files. You cannot use different ROMs for different save files.

4. **No ROM browser**: You must manually locate your ROM file. The editor does not scan for ROMs automatically.

## Future Enhancements

- Battle sprite extraction (32×32, 6 frames)
- Palette extraction from ROM
- Automatic ROM detection/scanning
- Support for hacked/modified ROMs
- Sprite editing with ROM preview
- Export back to ROM format

## Legal Notice

This editor does NOT include any ROM files. You must provide your own legally obtained Final Fantasy VI ROM. The editor only reads sprite data from ROM files that you already own.

## See Also

- [Combat Pack Quick Start](COMBAT_PACK_QUICK_START.md)
- [Phase 13 Sprite Editor](PHASE_13_SPRITE_EDITOR_COMPLETE.md)
- [Documentation Index](DOCUMENTATION_INDEX.md)
