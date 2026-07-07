# my-agent

AI coding agent di terminal. Baca/tulis file, git, search web, jalankan command — semua dari TUI.

## Cara Install

### Opsi 1: PowerShell (paling mudah) ⭐

Buka PowerShell (`Win + R` → ketik `powershell` → Enter), lalu:

```powershell
irm https://raw.githubusercontent.com/Gopartner/my-agent/main/install.ps1 | iex
```

Tunggu proses selesai. Kalau diminta, ketik `y` untuk lanjut.
Selesai langsung ketik: `my-agent`

### Opsi 2: Download EXE langsung

1. Buka https://github.com/Gopartner/my-agent/releases
2. Klik `my-agent-windows-amd64.exe`
3. Klik 2x file hasil download

**SmartScreen muncul?** Klik tulisan **"More info"** → tombol **"Run anyway"** muncul → klik. Ini wajar karena file belum di-sign.

### Opsi 3: Dari source (butuh Go)

```powershell
git clone https://github.com/Gopartner/my-agent.git
cd my-agent
go install .
my-agent
```

## Pertama Kali Jalan

Muncul layar selamat datang:

```
  Token: hf_...
```

Paste Hugging Face token kamu → Enter.
Token disimpan otomatis. Besoknya tinggal jalanin lagi, langsung masuk TUI.

> Token bisa dibuat gratis di https://huggingface.co/settings/tokens

## Cara Pakai

```powershell
my-agent
```

Tinggal ketik perintah, AI ngerjain pakai tools. Contoh:

```
You: baca file main.go
You: cari semua function yang pake error handling
You: git status
You: commit semua perubahan dengan pesan "fix bug"
```

## Keybindings

| Key | Aksi |
|---|---|
| `Enter` | Kirim pesan |
| `Ctrl+C` / `Esc` | Keluar |

## Build manual

```bash
git clone https://github.com/Gopartner/my-agent.git
cd my-agent
go build -ldflags="-s -w" -o my-agent.exe .
```

## Lisensi

MIT
