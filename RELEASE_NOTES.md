## v0.1.0 — Rilis Perdana

AI coding agent di terminal. Baca/tulis file, git, search web, jalankan command — dari terminal.

### Install

**Windows** — Buka PowerShell, paste:
```powershell
irm https://raw.githubusercontent.com/Gopartner/my-agent/main/install.ps1 | iex
```

**Mac / Linux** — Buka Terminal, paste:
```bash
curl -fsSL https://raw.githubusercontent.com/Gopartner/my-agent/main/install.sh | bash
```

**Kalau sudah punya Go:**
```bash
go install github.com/gopartner/my-agent@latest
```

### Cara pakai
```bash
my-agent
```

Pertama jalan tinggal paste Hugging Face token. Sekali doang.

### Fitur
- TUI chat dengan streaming realtime
- 14 tools: file, git, shell, web search, dll
- First-run wizard (token disimpan otomatis)
