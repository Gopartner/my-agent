param(
    [string]$InstallDir = "$HOME\my-agent"
)

$ErrorActionPreference = "Stop"
$repo = "https://github.com/Gopartner/my-agent"

Write-Host ""
Write-Host "  ========================================" -ForegroundColor Cyan
Write-Host "     my-agent Installer" -ForegroundColor Cyan
Write-Host "  ========================================" -ForegroundColor Cyan
Write-Host ""

# Cek Go
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Go tidak ditemukan. Install Go dulu..." -ForegroundColor Yellow
    if (Get-Command winget -ErrorAction SilentlyContinue) {
        winget install GoLang.Go --silent
        $env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [Environment]::GetEnvironmentVariable("Path", "User")
        if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
            Write-Host "Gagal install Go otomatis. Install manual: https://go.dev/dl/" -ForegroundColor Red
            exit 1
        }
    } else {
        Write-Host "Download Go dari: https://go.dev/dl/" -ForegroundColor Red
        exit 1
    }
}

# Build dari source (paling aman, no SmartScreen)
Write-Host "Mengunduh & build my-agent..." -ForegroundColor Cyan
$tmpDir = "$env:TEMP\my-agent-install"
if (Test-Path $tmpDir) { Remove-Item -Recurse -Force $tmpDir }
git clone --depth 1 $repo $tmpDir 2>&1 | Out-Null

Push-Location $tmpDir
go build -ldflags="-s -w" -o my-agent.exe .
if ($LASTEXITCODE -ne 0) {
    Pop-Location
    Write-Host "Build gagal!" -ForegroundColor Red
    exit 1
}
Pop-Location

# Buat folder install
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

Copy-Item -Force "$tmpDir\my-agent.exe" "$InstallDir\my-agent.exe"
Remove-Item -Recurse -Force $tmpDir

# Tambah PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$InstallDir;$userPath", "User")
    $env:Path = "$InstallDir;$env:Path"
    Write-Host "  ✔ Folder ditambahkan ke PATH: $InstallDir" -ForegroundColor Green
}

Write-Host ""
Write-Host "  ✔ my-agent siap!" -ForegroundColor Green
Write-Host ""
Write-Host "  Jalankan: my-agent" -ForegroundColor White
Write-Host ""
Write-Host "  Pertama kali jalan, kamu tinggal masukkan" -ForegroundColor Gray
Write-Host "  Hugging Face token, sisanya otomatis." -ForegroundColor Gray
