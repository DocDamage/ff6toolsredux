# Check for large files in the repository
# Run with: .\scripts\check-file-sizes.ps1

$MaxSizeBytes = 5MB
$MaxSizeHuman = "5MB"

Write-Host "Checking for large files (> $MaxSizeHuman) in repository..." -ForegroundColor Cyan

# Get all tracked files
$largeFiles = @()
$files = git ls-files

foreach ($file in $files) {
    if (Test-Path $file -PathType Leaf) {
        $size = (Get-Item $file).Length
        if ($size -gt $MaxSizeBytes) {
            $sizeMB = [math]::Round($size / 1MB, 2)
            $largeFiles += [PSCustomObject]@{
                File = $file
                Size = "$sizeMB MB"
                Bytes = $size
            }
        }
    }
}

if ($largeFiles.Count -gt 0) {
    Write-Host "`nWARNING: Found $($largeFiles.Count) large files (> $MaxSizeHuman):" -ForegroundColor Red
    Write-Host "These files should be removed from git history and added to .gitignore:" -ForegroundColor Yellow
    $largeFiles | Format-Table -AutoSize
    
    Write-Host "Recommendations:" -ForegroundColor Yellow
    Write-Host "  1. Remove large files: git rm --cached <file>"
    Write-Host "  2. Add to .gitignore: echo '<file>' >> .gitignore"
    Write-Host "  3. For build artifacts, use GitHub Releases instead"
    Write-Host "  4. For necessary large files, consider git-lfs"
    exit 1
} else {
    Write-Host "`n✓ No large files found. Repository size is healthy." -ForegroundColor Green
    exit 0
}
