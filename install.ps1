#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Instal my-agent — AI coding agent di terminal.
    Mirip kayak install Node.js: satu perintah, langsung bisa dipakai.
.DESCRIPTION
    - Otomatis install Go lewat winget kalo belum ada
    - Install my-agent via go install
    - Tambah PATH jika perlu
    - Token diminta pas pertama kali jalan, bukan pas install
.EXAMPLE
    irm https://raw.githubusercontent.com/Gopartner/my-agent/main/install.ps1 | iex
#>

$ErrorActionPreference = "Stop"
$repo = "github.com/gopartner/my-agent"

function Print-Banner {
    Write-Host ""
    Write-Host "  ╔═══════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "  ║       my-agent Installer          ║" -ForegroundColor Cyan
    Write-Host "  ║   AI coding agent di terminal     ║" -ForegroundColor Cyan
    Write-Host "  ╚═══════════════════════════════════╝" -ForegroundColor Cyan
    Write-Host ""
}

function Ensure-Go {
    if (Get-Command go -ErrorAction SilentlyContinue) {
        Write-Host "  ✔ Go $(go version)" -ForegroundColor Green
        return
    }
    Write-Host "  Menginstall Go..." -ForegroundColor Yellow
    if (-not (Get-Command winget -ErrorAction SilentlyContinue)) {
        Write-Host "  ✖ winget tidak tersedia. Install Go dulu:" -ForegroundColor Red
        Write-Host "    https://go.dev/dl/" -ForegroundColor White
        exit 1
    }
    winget install GoLang.Go --silent --accept-package-agreements 2>&1 | Out-Null
    # Reload PATH
    $env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [Environment]::GetEnvironmentVariable("Path", "User")
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Host "  ✖ Gagal install Go otomatis." -ForegroundColor Red
        Write-Host "    Install manual: https://go.dev/dl/" -ForegroundColor White
        exit 1
    }
    Write-Host "  ✔ Go $(go version)" -ForegroundColor Green
}

function Install-MyAgent {
    Write-Host "  Menginstall my-agent..." -ForegroundColor Yellow
    $output = & { go install "$repo@latest" 2>&1 } | %{ "$_" }
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  ✖ Gagal install: $($output -join "`n")" -ForegroundColor Red
        exit 1
    }
    Write-Host "  ✔ my-agent terinstall" -ForegroundColor Green
}

function Ensure-PATH {
    $targetDir = "$HOME\go\bin"
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$targetDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$targetDir;$userPath", "User")
        $env:Path = "$targetDir;$env:Path"
        Write-Host "  ✔ PATH ditambahkan: $targetDir" -ForegroundColor Green
    }
}

function Show-Done {
    Write-Host ""
    Write-Host "  ───────────────────────────────────" -ForegroundColor Gray
    Write-Host "  ✔ my-agent siap dipakai!" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Jalankan di terminal:" -ForegroundColor White
    Write-Host "    my-agent" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  Pertama kali jalan, tinggal masukkin" -ForegroundColor Gray
    Write-Host "  Hugging Face token. Sekali doang." -ForegroundColor Gray
    Write-Host "  Token: https://huggingface.co/settings/tokens" -ForegroundColor Gray
    Write-Host "  ───────────────────────────────────" -ForegroundColor Gray
    Write-Host ""
}

Print-Banner
Ensure-Go
Install-MyAgent
Ensure-PATH
Show-Done
