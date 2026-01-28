# Setup Build Environment for FF6 Save Editor
# Adds MinGW-w64 GCC to PATH for CGO support

Write-Host "Setting up build environment..." -ForegroundColor Cyan

# Add MinGW-w64 to PATH
$mingwPath = "C:\msys64\mingw64\bin"
if (Test-Path $mingwPath) {
    $env:Path = "$mingwPath;$env:Path"
    Write-Host "✓ Added MinGW-w64 to PATH" -ForegroundColor Green
} else {
    Write-Host "✗ MinGW-w64 not found at $mingwPath" -ForegroundColor Red
    Write-Host "  Run: winget install MSYS2.MSYS2" -ForegroundColor Yellow
    exit 1
}

# Enable CGO
$env:CGO_ENABLED = 1
Write-Host "✓ CGO_ENABLED=1" -ForegroundColor Green

# Verify GCC is available
try {
    $gccVersion = gcc --version 2>&1 | Select-Object -First 1
    Write-Host "✓ GCC found: $gccVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ GCC not available" -ForegroundColor Red
    exit 1
}

# Verify Go is available
try {
    $goVersion = go version
    Write-Host "✓ Go found: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ Go not available" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Build environment ready!" -ForegroundColor Green
Write-Host "You can now run: go build -o ffvi-save-editor.exe" -ForegroundColor Cyan
