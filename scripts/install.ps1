# prgen installer -- Windows (PowerShell)
# Usage: PowerShell -ExecutionPolicy Bypass -File scripts\install.ps1
$ErrorActionPreference = "Stop"

$RepoRoot    = Split-Path -Parent $PSScriptRoot
$BinName     = "prgen.exe"
$VersionFile = Join-Path $RepoRoot "VERSION"

$Version = "dev"
if (Test-Path $VersionFile) {
    $Version = (Get-Content $VersionFile -Raw).Trim()
}
$BuildDate = (Get-Date -Format "yyyy-MM-dd")

Write-Host ""
Write-Host "[prgen] Building v$Version ..." -ForegroundColor Cyan

# Build
Push-Location $RepoRoot
try {
    $ldflags = "-s -w -X github.com/alann-estrada-KSH/ai-pr-generator/internal/version.Version=$Version -X github.com/alann-estrada-KSH/ai-pr-generator/internal/version.BuildDate=$BuildDate"
    go build -ldflags $ldflags -o $BinName ./cmd/prgen
    if ($LASTEXITCODE -ne 0) { throw "Build failed with exit code $LASTEXITCODE" }
    Write-Host "[prgen] Build OK" -ForegroundColor Green
} finally {
    Pop-Location
}

# Install location
$InstallDir = Join-Path $env:USERPROFILE "AppData\Local\prgen"
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

$BinPath = Join-Path $InstallDir $BinName
Copy-Item (Join-Path $RepoRoot $BinName) $BinPath -Force
Write-Host "[prgen] Binary installed to: $BinPath" -ForegroundColor Green

# PATH (permanent, user scope)
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    $NewPath = "$CurrentPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    $env:PATH = "$env:PATH;$InstallDir"
    Write-Host "[prgen] Added $InstallDir to user PATH" -ForegroundColor Green
} else {
    Write-Host "[prgen] Install dir already in PATH" -ForegroundColor DarkGray
}

# Default config
$ConfigDir = Join-Path $env:USERPROFILE ".prgen"
if (-not (Test-Path $ConfigDir)) {
    New-Item -ItemType Directory -Path $ConfigDir | Out-Null
}

$ConfigDest = Join-Path $ConfigDir "config.yaml"
if (-not (Test-Path $ConfigDest)) {
    Copy-Item (Join-Path $RepoRoot "config.yaml") $ConfigDest
    Write-Host "[prgen] Config copied to $ConfigDest" -ForegroundColor Green
}

# Cleanup
Remove-Item (Join-Path $RepoRoot $BinName) -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "[prgen] Installation complete!" -ForegroundColor Green
Write-Host ""

try {
    & $BinPath version
} catch {
    Write-Host "[prgen] Restart your terminal and run: prgen version" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Tip: edit your config at: $ConfigDest" -ForegroundColor DarkCyan
Write-Host ""
