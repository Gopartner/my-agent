## v0.1.0 — Rilis Perdana

AI coding agent di terminal. Bisa baca/tulis file, git, search web, jalankan command, dll.

### Cara Install (paling mudah)

Buka PowerShell, paste ini:

```powershell
irm https://raw.githubusercontent.com/Gopartner/my-agent/main/install.ps1 | iex
```

Selesai, ketik: `my-agent`

### Atau download langsung

| Platform | File |
|---|---|
| Windows | `my-agent-windows-amd64.exe` |
| Linux | `my-agent-linux-amd64` |
| macOS Intel | `my-agent-darwin-amd64` |
| macOS Apple Silicon | `my-agent-darwin-arm64` |

> SmartScreen muncul? Klik **More info** → **Run anyway** — ini wajar karena file belum di-sign.

### Fitur
- TUI chat dengan streaming token realtime
- 14 tools: read_file, write_file, edit_file, run_command, search_code, git, web_search, dll
- First-run wizard: tinggal masukkan token sekali, otomatis tersimpan
- Session persistence
