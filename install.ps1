param(
    [string]$InstallDir = "$HOME\go\bin"
)

$ErrorActionPreference = "Stop"

# Cek Go
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Go belum terinstall. Download: https://go.dev/dl/" -ForegroundColor Red
    exit 1
}

Write-Host "Building my-agent..." -ForegroundColor Cyan
go build -ldflags="-s -w" -o my-agent.exe .
if ($LASTEXITCODE -ne 0) { exit 1 }

# Pastikan folder tujuan ada
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

Copy-Item -Force my-agent.exe "$InstallDir\my-agent.exe"
Remove-Item my-agent.exe

# Cek apakah InstallDir ada di PATH user
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$InstallDir*") {
    $newPath = "$InstallDir;$userPath"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + $newPath
    Write-Host "Path ditambahkan: $InstallDir" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "✔ my-agent terinstall!" -ForegroundColor Green
Write-Host ""
Write-Host "Cara pakai:" -ForegroundColor Cyan
Write-Host "  1. Set token:" -ForegroundColor White
Write-Host "     `$env:MY_AGENT_TOKEN = `"hf_...`"" -ForegroundColor Gray
Write-Host "  2. Jalankan:" -ForegroundColor White
Write-Host "     my-agent" -ForegroundColor Gray
Write-Host ""
Write-Host "Tips: Simpan token permanent:" -ForegroundColor Yellow
Write-Host "  [Environment]::SetEnvironmentVariable(""MY_AGENT_TOKEN"", ""hf_..."", ""User"")" -ForegroundColor Gray
