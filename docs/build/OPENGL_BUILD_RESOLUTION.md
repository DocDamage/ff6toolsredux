# OpenGL Build Issue - Resolution

## Problem
The main application fails to build on Windows due to missing CGO/GCC dependencies required by the Fyne OpenGL backend:
```
imports github.com/go-gl/gl/v2.1/gl: build constraints exclude all Go files
```

## Root Cause
- Fyne UI framework requires OpenGL bindings via CGO
- Windows Go installations have `CGO_ENABLED=0` by default
- OpenGL bindings require a C compiler (GCC/MinGW)
- No C compiler found in PATH

## Solution: CLI-Only Build

Created `main_cli.go` - a CLI-only entry point that bypasses GUI dependencies entirely.

### Building CLI Version
```bash
go build -o ffvi_editor_cli.exe .\main_cli.go
```

**Result:** ✅ Builds successfully, no OpenGL/CGO required

### Testing CLI Build
```bash
.\ffvi_editor_cli.exe combat-pack --mode smoke
```

**Result:** ✅ All Combat Pack features working (9/9 tests passed)

## Available Solutions

### Option 1: CLI-Only Build (Recommended for Development)
**Pros:**
- ✅ No additional dependencies
- ✅ Fast compilation
- ✅ All Combat Pack features work
- ✅ Suitable for automation/scripting

**Cons:**
- ❌ No GUI interface
- ❌ Manual command-line usage required

**Build:**
```bash
go build -o ffvi_editor_cli.exe .\main_cli.go
```

### Option 2: Enable CGO + Install MinGW (Full GUI)
**Pros:**
- ✅ Full GUI interface
- ✅ All features accessible
- ✅ User-friendly UI

**Cons:**
- ❌ Requires C compiler installation
- ❌ Slower compilation
- ❌ Additional ~500MB dependencies

**Setup:**
1. Download MinGW-w64: https://www.mingw-w64.org/downloads/
2. Add MinGW `bin` to PATH
3. Enable CGO: `$env:CGO_ENABLED = "1"`
4. Build: `go build -o ffvi_editor.exe .`

### Option 3: Pre-built Binaries (Distribution)
**Recommended for end users**
- Build on a system with CGO/GCC
- Distribute compiled `.exe` files
- No build dependencies for end users

## Current State

### ✅ Working (CLI Build)
- Combat Depth Pack (all modes)
- Save file manipulation
- Lua script execution
- Smoke tests
- Encounter/Boss/Companion operations

### ⚠️ Not Available (CLI Build)
- GUI dialogs
- Visual save editor
- Drag-and-drop operations
- UI-based plugin management

## Recommendation

**For Development/Testing:**
Use CLI build (`main_cli.go`) - fastest and no dependencies.

**For Distribution:**
Build GUI version on a system with MinGW installed, then distribute the compiled `.exe`.

**For End-to-End Testing:**
CLI build is sufficient to verify all Combat Pack functionality including save manipulation.

## CLI Usage Examples

### Run Smoke Tests
```bash
.\ffvi_editor_cli.exe combat-pack --mode smoke
```

### With Save File
```bash
.\ffvi_editor_cli.exe combat-pack --mode smoke --file mysave.json
```

### Encounter Tuning
```bash
.\ffvi_editor_cli.exe combat-pack --mode encounter --zone "Mt. Kolts" --rate 1.2 --elite 0.15 --file mysave.json
```

### Boss Remix
```bash
.\ffvi_editor_cli.exe combat-pack --mode boss --affixes "enraged,arcane_shield" --file mysave.json
```

### Companion Director
```bash
.\ffvi_editor_cli.exe combat-pack --mode companion --profile aggressive --risk high --file mysave.json
```

## Status
✅ **OpenGL issue bypassed** - CLI build fully functional for all Combat Pack operations.
