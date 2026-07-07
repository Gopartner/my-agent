#!/usr/bin/env bash
set -euo pipefail

REPO="github.com/gopartner/my-agent"

echo ""
echo "  ╔═══════════════════════════════════╗"
echo "  ║       my-agent Installer          ║"
echo "  ║   AI coding agent di terminal     ║"
echo "  ╚═══════════════════════════════════╝"
echo ""

# Cek Go
if ! command -v go &>/dev/null; then
    echo "  Menginstall Go..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &>/dev/null; then
            brew install go
        else
            echo "  ✖ Install Homebrew dulu: https://brew.sh/"
            echo "    Atau download Go: https://go.dev/dl/"
            exit 1
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt &>/dev/null; then
            sudo apt update && sudo apt install -y golang-go
        elif command -v pacman &>/dev/null; then
            sudo pacman -S --noconfirm go
        else
            echo "  ✖ Install Go manual: https://go.dev/dl/"
            exit 1
        fi
    else
        echo "  ✖ OS tidak dikenal. Install Go manual: https://go.dev/dl/"
        exit 1
    fi
fi

echo "  ✔ Go $(go version)"

# Install my-agent
echo "  Menginstall my-agent..."
go install "$REPO@latest"
echo "  ✔ my-agent terinstall"

# Cek PATH
GOPATH=$(go env GOPATH 2>/dev/null || echo "$HOME/go")
GOBIN=$(go env GOBIN 2>/dev/null || echo "$GOPATH/bin")

if [[ ":$PATH:" != *":$GOBIN:"* ]]; then
    echo ""
    echo "  ⚠ Tambahkan ini ke ~/.bashrc atau ~/.zshrc:"
    echo "    export PATH=\"\$PATH:$GOBIN\""
    echo ""
fi

echo ""
echo "  ───────────────────────────────────"
echo "  ✔ my-agent siap dipakai!"
echo ""
echo "  Jalankan di terminal:"
echo "    my-agent"
echo ""
echo "  Pertama kali jalan, tinggal masukkin"
echo "  Hugging Face token. Sekali doang."
echo "  ───────────────────────────────────"
echo ""
