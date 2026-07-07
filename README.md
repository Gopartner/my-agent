# my-agent

AI coding agent di terminal — kayak ChatGPT tapi bisa baca/tulis file, git, search web, dan menjalankan perintah.

Cross-platform: **Windows 10/11** · **macOS** · **Linux** · **Termux (Android)**

## Install

### Windows

```powershell
powershell -c "irm https://raw.githubusercontent.com/Gopartner/my-agent/main/install.ps1 | iex"
```

Lalu: `my-agent`

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/Gopartner/my-agent/main/install.sh | bash
```

Lalu: `my-agent`

### Termux (Android)

```bash
pkg update && pkg install golang -y
go install github.com/gopartner/my-agent@latest
my-agent
```

### Kalau sudah punya Go (semua OS)

```bash
go install github.com/gopartner/my-agent@latest
```

## Cara Pakai

```bash
my-agent
```

Pertama kali jalan, paste Hugging Face token. Sekali doang.

> Token: https://huggingface.co/settings/tokens

## Keybindings

| Key | Aksi |
|---|---|
| `Enter` | Kirim |
| `Ctrl+C` / `Esc` | Keluar |
| `←` / `→` | Gerak cursor |

## Tools

`read_file` · `write_file` · `edit_file` · `delete_file` · `list_dir` · `project_tree` ·
`run_command` · `search_code` · `git_status` · `git_diff` · `git_commit` ·
`web_search` · `web_fetch` · `http_request`

## Lisensi

MIT
