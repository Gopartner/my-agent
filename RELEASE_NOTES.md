## v0.1.0 — Rilis Perdana

AI coding agent di terminal. Cross-platform: Windows, macOS, Linux, Termux.

### Install

**Windows** — PowerShell:
```powershell
powershell -c "irm https://raw.githubusercontent.com/Gopartner/my-agent/main/install.ps1 | iex"
```

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/Gopartner/my-agent/main/install.sh | bash
```
**Termux (Android):**
```bash
pkg update && pkg install golang -y
go install github.com/gopartner/my-agent@latest
```

**Kalau sudah punya Go (semua OS):**
```bash
go install github.com/gopartner/my-agent@latest
```

### Binary yang disediakan
- `my-agent-windows-amd64.exe`
- `my-agent-linux-amd64`
- `my-agent-linux-arm64` (Termux / Raspberry Pi)
- `my-agent-darwin-amd64` (macOS Intel)
- `my-agent-darwin-arm64` (macOS Apple Silicon)

### Fitur
- TUI chat streaming realtime
- 14 tools: file, git, shell, web search, dll
- First-run wizard (token disimpan otomatis)
