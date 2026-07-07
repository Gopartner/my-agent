# my-agent

AI coding agent berbasis **DeepSeek-V3.1** (via Hugging Face Router) dengan **TUI modern** (bubbletea) dan **tool calling** lengkap.

## Fitur

- **Chat TUI** — viewport + markdown rendering (glamour) + input with cursor
- **Streaming realtime** — konten langsung muncul token per token
- **Tool calling** — AI bisa menggunakan 14 tools secara otonom
- **Session persistence** — history obrolan tersimpan otomatis ke `.agent_session.json`
- **Dark theme** — lipgloss styling dengan accent `#7C3AED`

## Tools yang Didukung

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
| `web_search` | Cari informasi di web (DuckDuckGo) |
| `web_fetch` | Ambil konten halaman web |
| `http_request` | Kirim HTTP request |

## Struktur Project

```
agent-go/
├── main.go                  # Entry point
├── internal/
│   ├── agent.go             # Agent logic & tool loop
│   ├── hfapi.go             # Hugging Face API client (SSE streaming)
│   ├── session.go           # Save/load session JSON
│   ├── styles.go            # Lipgloss theme
│   ├── tools.go             # 14 tool definitions + execution
│   └── tui.go               # Bubbletea TUI (viewport, input, spinner)
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.21+
- Hugging Face API key ([dapatkan disini](https://huggingface.co/settings/tokens))

## Instalasi & Usage

```powershell
# Clone atau cd ke folder
cd agent-go

# Set API key
$env:HF_TOKEN = "hf_..."

# Build
go build -o agent-go.exe .

# Jalankan
./agent-go.exe
```

### Keybindings

| Key | Aksi |
|---|---|
| `Enter` | Kirim pesan |
| `Backspace` | Hapus karakter sebelumnya |
| `Delete` | Hapus karakter setelah cursor |
| `←` / `→` | Gerakkan cursor |
| `Home` / `End` | Lompat ke awal/akhir input |
| `Ctrl+U` | Hapus seluruh input |
| `Ctrl+C` / `Esc` | Keluar |

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) — Styling
- [bubbles](https://github.com/charmbracelet/bubbles) — Viewport & spinner
- [glamour](https://github.com/charmbracelet/glamour) — Markdown rendering
- [goquery](https://github.com/PuerkitoBio/goquery) — HTML parsing (web tools)

## API

Menggunakan Hugging Face Inference Router:

```
POST https://router.huggingface.co/v1/chat/completions
Authorization: Bearer $HF_TOKEN
Model: deepseek-ai/DeepSeek-V3.1
```

## Lisensi

MIT
