package internal

import (
	"fmt"
	"runtime"
	"strings"
)

const Version = "v0.1.0"

func InfoText() string {
	var b strings.Builder
	b.WriteString("╔══════════════════════════════════════════╗\n")
	b.WriteString(fmt.Sprintf("║  my-agent %s", Version))
	for i := 0; i < 26-len(Version); i++ { b.WriteString(" ") }
	b.WriteString("║\n")
	b.WriteString("║  AI coding agent di terminal             ║\n")
	b.WriteString("╠══════════════════════════════════════════╣\n")
	b.WriteString(fmt.Sprintf("║  OS:   %s/%s", runtime.GOOS, runtime.GOARCH))
	for i := 0; i < 24-len(runtime.GOOS)-len(runtime.GOARCH); i++ { b.WriteString(" ") }
	b.WriteString("║\n")
	b.WriteString(fmt.Sprintf("║  Go:   %s", runtime.Version()))
	for i := 0; i < 28-len(runtime.Version()); i++ { b.WriteString(" ") }
	b.WriteString("║\n")
	b.WriteString(fmt.Sprintf("║  API:  DeepSeek-V3.1 (HF Router)"))
	b.WriteString("       ║\n")
	b.WriteString("╠══════════════════════════════════════════╣\n")
	b.WriteString("║  Tools:                                   ║\n")
	b.WriteString("║  • read_file, write_file, edit_file      ║\n")
	b.WriteString("║  • delete_file, list_dir, project_tree   ║\n")
	b.WriteString("║  • run_command, search_code              ║\n")
	b.WriteString("║  • git_status, git_diff, git_commit      ║\n")
	b.WriteString("║  • web_search, web_fetch, http_request   ║\n")
	b.WriteString("╠══════════════════════════════════════════╣\n")
	b.WriteString("║  Commands:                                ║\n")
	b.WriteString("║  /info    — info ini                      ║\n")
	b.WriteString("║  /aturan  — lihat aturan AI               ║\n")
	b.WriteString("║  /baru    — hapus sesi, mulai baru        ║\n")
	b.WriteString("╚══════════════════════════════════════════╝\n")
	return b.String()
}

func ToolListText() string {
	defs := GetToolDefs()
	var b strings.Builder
	b.WriteString("Tools yang tersedia:\n\n")
	for _, d := range defs {
		b.WriteString(fmt.Sprintf("  • %s — %s\n", d.Function.Name, d.Function.Description))
	}
	return b.String()
}
