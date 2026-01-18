# Build Instructions - Final Fantasy VI Save Editor

## Prerequisites

### 1. Go 1.25.6 (Installed ✅)
Located at: `C:\Program Files\Go`

### 2. MinGW-w64 GCC Compiler (Installed ✅)
Located at: `C:\msys64\mingw64\bin`

The GCC compiler is required for CGO (C bindings) which Fyne uses for OpenGL rendering.

## Quick Build

### Option 1: Using Setup Script (Recommended)
```powershell
.\setup-build-env.ps1
go build -o ffvi-save-editor.exe
```

### Option 2: Manual Build
```powershell
$env:Path = "C:\msys64\mingw64\bin;$env:Path"
$env:CGO_ENABLED = 1
go build -o ffvi-save-editor.exe
```

## Build Output
- **Executable**: `ffvi-save-editor.exe` (~40 MB)
- **Build Time**: ~30-60 seconds on first build
- **Subsequent Builds**: ~5-10 seconds (cached)

## Running the Application

Simply double-click `ffvi-save-editor.exe` or run:
```powershell
.\ffvi-save-editor.exe
```

## Testing

Run all tests:
```powershell
go test ./...
```

Run specific package tests:
```powershell
go test ./io/pr -v
go test ./models/pr -v
```

Run benchmarks:
```powershell
go test ./io/pr -bench=. -benchmem
```

## Build Options

### Release Build (Smaller Size)
```powershell
go build -ldflags="-s -w" -o ffvi-save-editor.exe
```
- `-s`: Strip symbol table
- `-w`: Strip DWARF debug info
- Result: ~30 MB (instead of 40 MB)

### Debug Build
```powershell
go build -gcflags="all=-N -l" -o ffvi-save-editor-debug.exe
```
- `-N`: Disable optimizations
- `-l`: Disable inlining
- Better for debugging with delve

## Troubleshooting

### "gcc: not found"
Run the setup script or manually add MinGW to PATH:
```powershell
$env:Path = "C:\msys64\mingw64\bin;$env:Path"
```

### "cgo: C compiler not found"
Ensure CGO is enabled:
```powershell
$env:CGO_ENABLED = 1
```

### Build Constraints Error
This means GCC is not in PATH. Run `.\setup-build-env.ps1` first.

### Slow First Build
The first build downloads and compiles all dependencies. Subsequent builds use Go's build cache and are much faster.

## Clean Build

Remove build cache and rebuild:
```powershell
go clean -cache
go build -o ffvi-save-editor.exe
```

## Cross-Compilation

To build for different platforms:

### Windows (from Windows)
```powershell
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o ffvi-save-editor-windows.exe
```

### Linux (requires Linux GCC)
```powershell
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o ffvi-save-editor-linux
```

Note: Cross-compilation with CGO requires platform-specific toolchains.

## Development Workflow

1. Make code changes
2. Run tests: `go test ./...`
3. Build: `go build -o ffvi-save-editor.exe`
4. Test the application
5. Commit changes

## Project Structure

```
ffvi_editor/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── browser/                # Auto-update functionality
├── global/                 # Global variables
├── io/                     # I/O operations
│   ├── config/            # Configuration management
│   ├── file/              # File I/O with encryption
│   └── pr/                # Save file parsing (main logic)
├── models/                 # Data models
│   ├── character.go       # Character data structures
│   ├── equipment.go       # Equipment definitions
│   ├── misc.go           # Miscellaneous game data
│   └── pr/               # PR-specific models
└── ui/                    # Fyne GUI components
    └── forms/            # Form editors
        └── editors/      # Specific data editors
```

## Dependencies

Major dependencies (auto-downloaded during build):
- **fyne.io/fyne/v2**: GUI framework (v2.5.2)
- **github.com/go-gl/gl**: OpenGL bindings
- **github.com/kiamev/ffpr-save-cypher**: Save encryption (v1.0.0)
- **gitlab.com/c0b/go-ordered-json**: JSON with key ordering

View all dependencies:
```powershell
go list -m all
```

## Performance

Benchmark results on 11th Gen Intel Core i7-1185G7 @ 3.00GHz:
- GetInt: 28.71 ns/op, 0 allocations
- GetString: 14.22 ns/op, 0 allocations
- LoadEquipment: 14.9 µs/op, 161 allocations
- LoadMapData: 31.8 µs/op, 189 allocations

Type-safe refactoring adds <5% overhead while providing significant safety benefits.
