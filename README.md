# my-agent

AI coding agent di terminal dengan **DeepSeek-V3.1**.

## Cara Pakai

### 1. Download

Download `my-agent-windows-amd64.exe` dari [Releases](https://github.com/Gopartner/my-agent/releases).

### 2. Taruh di folder mana aja

Bikin folder khusus, misal `C:\Users\Rafka\my-agent\`, taruh file `.exe` di situ.

### 3. Jalankan (double-click)

Klik 2x `my-agent.exe`. Muncul layar selamat datang → masukkan **Hugging Face token** → enter.

Token disimpan otomatis. Next time tinggal double-click lagi, langsung masuk TUI.

> **Biar bisa dipanggil dari terminal mana aja:** tambahin folder `my-agent` ke PATH. Atau jalankan `install.ps1` dari repo ini (butuh Go).

### Preview

```
  ╔══════════════════════════════════════════╗
  ║          Selamat datang di my-agent!     ║
  ║                                          ║
  ║  Masukkan Hugging Face token kamu.       ║
  ║  https://huggingface.co/settings/tokens  ║
  ║                                          ║
  ║  Token akan disimpan di:                 ║
  ║  C:\Users\Rafka\AppData\my-agent\...     ║
  ╚══════════════════════════════════════════╝
  Token: hf_...
```

## Keybindings

| Key | Aksi |
|---|---|
| `Enter` | Kirim pesan |
| `Ctrl+C` / `Esc` | Keluar |

## Build dari Source

```bash
git clone https://github.com/Gopartner/my-agent.git
cd my-agent
go build -ldflags="-s -w" -o my-agent.exe .
```

## Lisensi

MIT
