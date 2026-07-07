# my-agent

AI coding agent berbasis **DeepSeek-V3.1** (via Hugging Face Router) dengan **TUI modern** (bubbletea) dan **tool calling** lengkap.

## Fitur

- **Chat TUI** — viewport + markdown rendering (glamour) + input with cursor
- **Streaming realtime** — token langsung muncul tanpa nunggu selesai
- **Tool calling** — AI otonom menggunakan 14 tools
- **Session persistence** — history tersimpan otomatis
- **Dark theme** — lipgloss styling dengan accent `#7C3AED`

## Tools

| Tool | Deskripsi |
|---|---|
| `read_file` | Baca file |
| `write_file` | Tulis file baru |
| `edit_file` | Edit file (find & replace) |
| `delete_file` | Hapus file/folder |
| `list_dir` | Lihat isi folder |
| `project_tree` | Struktur project |
| `run_command` | Jalankan shell command |
| `search_code` | Cari teks dalam project |
| `git_status` | Status git |
| `git_diff` | Lihat perubahan |
| `git_commit` | Commit semua perubahan |
| `web_search` | Cari informasi di web |
| `web_fetch` | Ambil konten halaman web |
| `http_request` | Kirim HTTP request |

## Persiapan

Buat Hugging Face token di https://huggingface.co/settings/tokens (permission `read`).

## Instalasi

### Opsi 1: Install Script (Windows)

```powershell
git clone https://github.com/Gopartner/my-agent.git
cd my-agent
.\install.ps1
```

Script akan build & copy `my-agent.exe` ke `$HOME\go\bin\` (otomatis ditambah ke PATH).

### Opsi 2: Manual (semua OS)

```bash
git clone https://github.com/Gopartner/my-agent.git
cd my-agent
go install .
```

Binary terinstall di `$GOPATH/bin/my-agent`.

### Opsi 3: Download Release

Download binary dari [Releases](https://github.com/Gopartner/my-agent/releases), extract, dan letakkan di folder yang ada di PATH.

## Menyetel Token

### Temporer (per sesi)

```powershell
$env:MY_AGENT_TOKEN = "hf_..."
my-agent
```

### Permanen (Windows)

```powershell
[Environment]::SetEnvironmentVariable("MY_AGENT_TOKEN", "hf_...", "User")
```

Buka terminal baru, lalu `my-agent`.

### Permanen (Linux/Mac)

Bash:
```bash
echo 'export MY_AGENT_TOKEN="hf_..."' >> ~/.bashrc
source ~/.bashrc
my-agent
```

## Keybindings

| Key | Aksi |
|---|---|
| `Enter` | Kirim pesan |
| `Backspace` | Hapus karakter sebelumnya |
| `Delete` | Hapus karakter setelah cursor |
| `←` / `→` | Gerakkan cursor |
| `Home` / `End` | Lompat ke awal/akhir input |
| `Ctrl+U` | Hapus seluruh input |
| `Ctrl+C` / `Esc` | Keluar |

## Struktur Project

```
internal/
├── agent.go    # Agent logic & tool calling loop
├── hfapi.go    # Hugging Face API client (SSE streaming)
├── session.go  # Save/load session JSON
├── styles.go   # Lipgloss theme (dark)
├── tools.go    # 14 tool definitions + execution
└── tui.go      # Bubbletea TUI (viewport, input, spinner)
```

## Build dari Source

```bash
go build -ldflags="-s -w" -o my-agent .
```

Cross-compile:

```bash
GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -o my-agent-linux .
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o my-agent.exe .
GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w" -o my-agent-mac .
```

## Lisensi

MIT
