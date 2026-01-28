# Setup Build Environment for FF6 Save Editor
# Adds MinGW-w64 GCC to PATH for CGO support

Write-Host "Setting up build environment..." -ForegroundColor Cyan

# Add MSYS2 toolchain to PATH (mingw64 or ucrt64)
$toolchainCandidates = @(
    "C:\msys64\mingw64\bin",
    "C:\msys64\ucrt64\bin"
)

$toolchainPath = $null
foreach ($candidate in $toolchainCandidates) {
    if (Test-Path (Join-Path $candidate "gcc.exe")) {
        $toolchainPath = $candidate
        break
    }
}

if ($null -eq $toolchainPath) {
    Write-Host "✗ GCC toolchain not found." -ForegroundColor Red
    Write-Host "  Expected one of:" -ForegroundColor Yellow
    foreach ($candidate in $toolchainCandidates) { Write-Host "    - $candidate" -ForegroundColor Yellow }
    Write-Host ""
    Write-Host "Install prerequisites:" -ForegroundColor Yellow
    Write-Host "  winget install MSYS2.MSYS2" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Then in the MSYS2 shell, install GCC:" -ForegroundColor Yellow
    Write-Host "  pacman -S --needed mingw-w64-x86_64-gcc" -ForegroundColor Yellow
    Write-Host "  # or" -ForegroundColor Yellow
    Write-Host "  pacman -S --needed mingw-w64-ucrt-x86_64-gcc" -ForegroundColor Yellow
    exit 1
}

$env:Path = "$toolchainPath;$env:Path"
Write-Host "✓ Added GCC toolchain to PATH: $toolchainPath" -ForegroundColor Green

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
    Write-Host "  Install with: winget install GoLang.Go" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "Build environment ready!" -ForegroundColor Green
Write-Host "You can now run: go build -o ffvi-save-editor.exe" -ForegroundColor Cyan
