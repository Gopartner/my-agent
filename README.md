# my-agent

AI coding agent di terminal — kayak ChatGPT tapi bisa baca/tulis file, git, search web, dan menjalankan perintah.

## Install

### Windows

Buka **PowerShell** (`Win + R` → `powershell` → Enter), lalu:

```powershell
irm https://raw.githubusercontent.com/Gopartner/my-agent/main/install.ps1 | iex
```

Selesai, ketik: `my-agent`

### Mac / Linux

Buka **Terminal**, lalu:

```bash
curl -fsSL https://raw.githubusercontent.com/Gopartner/my-agent/main/install.sh | bash
```

Selesai, ketik: `my-agent`

### Kalau sudah punya Go

```bash
go install github.com/gopartner/my-agent@latest
```

## Cara Pakai

```bash
my-agent
```

Pertama kali jalan, tinggal paste Hugging Face token. Sekali doang, tersimpan otomatis.

> Token: https://huggingface.co/settings/tokens

Lanjut ngobrol sama AI:

```
You: baca file main.go
You: cari bug di folder src/
You: git commit "fix typo"
```

## Keybindings

| Key | Aksi |
|---|---|
| `Enter` | Kirim |
| `Ctrl+C` / `Esc` | Keluar |

## Tools

`read_file` · `write_file` · `edit_file` · `delete_file` · `list_dir` · `project_tree` ·
`run_command` · `search_code` · `git_status` · `git_diff` · `git_commit` ·
`web_search` · `web_fetch` · `http_request`

## Build dari source

```bash
git clone https://github.com/Gopartner/my-agent.git
cd my-agent
go build -ldflags="-s -w" -o my-agent .
```

## Lisensi

MIT
