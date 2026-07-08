# my-agent — Project Guide

## Overview

AI coding agent berbasis **DeepSeek-V3.1** (via Hugging Face Router) dengan **TUI modern** (bubbletea). Berjalan di terminal, bisa baca/tulis file, git, search web, dan menjalankan perintah.

- **Language:** Go 1.21+
- **Module:** `github.com/gopartner/my-agent`
- **License:** MIT
- **Platform:** Windows, macOS, Linux, Termux (Android)

## Architecture

```
main.go                  # Entry point, token setup, starts TUI
internal/
├── agent.go             # Core agent: tool calling loop, message history
├── hfapi.go             # Hugging Face API client (SSE streaming)
├── tools.go             # 14 tool definitions + execution
├── tui.go               # Bubbletea TUI (viewport, input, spinner)
├── config.go            # First-run wizard, token storage
├── session.go           # Session save/load/archive, summary
├── info.go              # /info command — version, tools, system info
├── rules.go             # /aturan command — AI rules configuration
└── styles.go            # Lipgloss dark theme
```

## Development Rules

1. **Bahasa Indonesia** untuk semua interaksi AI, kode komentar, dan dokumentasi
2. **No exported functions** unless necessary — keep package internal
3. **Error messages** harus informatif dan dalam Bahasa Indonesia
4. **Minimal dependencies** — prefer standard library
5. **All files** in `internal/` package
6. **Session persistence** via `%APPDATA%\my-agent\sessions\` (Windows) or `~/.config/my-agent/sessions/` (Unix)

## Key Components

### HFClient (`hfapi.go`)
- Endpoint: `https://router.huggingface.co/v1/chat/completions`
- Auth: `Bearer MY_AGENT_TOKEN` (env var or config file)
- Streaming SSE, accumulates tool calls via `ToolCallAccum`
- Content deltas via `onContent`, tool call deltas via `onToolCall`

### Agent (`agent.go`)
- `ProcessRun()` — one iteration: send messages → stream → execute tools → return
- `ProcessFull()` — loop up to 15 iterations until no more tool calls
- Message history grows with `role: tool` results appended

### Tools (`tools.go`)
14 tools, all in `ExecuteTool()` switch:
- File: `read_file`, `write_file`, `edit_file`, `delete_file`
- Directory: `list_dir`, `project_tree`
- Shell: `run_command` (auto-detect `cmd /C` vs `sh -c`)
- Search: `search_code`
- Git: `git_status`, `git_diff`, `git_commit`
- Web: `web_search` (DuckDuckGo HTML scrape), `web_fetch`
- HTTP: `http_request`

### TUI (`tui.go`)
- Bubbletea model with viewport (chat), text input, spinner
- Commands:
  - `/info` — system info
  - `/aturan` — rules display
  - `/baru` — new session (archives old one)
- Streaming: goroutine sends `StreamEvent` / `ToolEvent` via `prog.Send()`
- Markdown rendering via glamour

### Config (`config.go`)
- Token resolution chain: `MY_AGENT_TOKEN` env var → `config.json` in app data dir → first-run prompt
- Config path: `%APPDATA%\my-agent\config.json` or `~/.config/my-agent/config.json`

### Session (`session.go`)
- Current session: `sessions/current.json`
- On `/baru`: archive to `sessions/session-{timestamp}.json`, show summary
- Summary includes: message count, tool calls, duration

## Build & Release

```bash
# Build local
go build -ldflags="-s -w" -o my-agent .

# Cross-compile
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o my-agent-linux-amd64 .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o my-agent.exe .
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o my-agent-darwin-arm64 .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o my-agent-linux-arm64 .

# Release: tag v* → GitHub Actions builds all platforms
git tag v1.0.0 && git push origin v1.0.0
```

## Dependencies

- `github.com/charmbracelet/bubbletea` — TUI framework
- `github.com/charmbracelet/lipgloss` — Styling
- `github.com/charmbracelet/bubbles` — Viewport & spinner widgets
- `github.com/charmbracelet/glamour` — Markdown rendering
- `github.com/PuerkitoBio/goquery` — HTML parsing (web tools)

## Install Methods

| Method | Command |
|---|---|
| PowerShell (Windows) | `powershell -c "irm https://raw.githubusercontent.com/Gopartner/my-agent/main/install.ps1 \| iex"` |
| Shell (Mac/Linux) | `curl -fsSL https://raw.githubusercontent.com/Gopartner/my-agent/main/install.sh \| bash` |
| Go install | `go install github.com/gopartner/my-agent@latest` |
| Termux | `pkg install golang -y && go install github.com/gopartner/my-agent@latest` |

## Future Improvements

- [ ] Configurable rules via TUI (toggle web/git)
- [ ] Multi-session management (list, resume)
- [ ] Streaming output per-tool (not just final result)
- [ ] Syntax highlighting in chat viewport
- [ ] Custom model selection (not just DeepSeek-V3.1)
- [ ] Context window management (trim old messages)
- [ ] Plugin system for custom tools
