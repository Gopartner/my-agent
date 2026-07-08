package internal

import (
	"fmt"
	"strings"
)

type Rules struct {
	Bahasa    string
	MaxIter   int
	MaxOutput int
	AllowWeb  bool
	AllowGit  bool
}

var DefaultRules = Rules{
	Bahasa:    "Indonesia",
	MaxIter:   15,
	MaxOutput: 3000,
	AllowWeb:  true,
	AllowGit:  true,
}

var ActiveRules = DefaultRules

func RulesText() string {
	r := ActiveRules
	var b strings.Builder
	b.WriteString("╔══════════════════════════════════════════╗\n")
	b.WriteString("║  Aturan AI Agent                        ║\n")
	b.WriteString("╠══════════════════════════════════════════╣\n")
	b.WriteString(fmt.Sprintf("║  Bahasa    : %s", r.Bahasa))
	for i := 0; i < 22-len(r.Bahasa); i++ { b.WriteString(" ") }
	b.WriteString("║\n")
	b.WriteString(fmt.Sprintf("║  Max loop  : %d iterasi", r.MaxIter))
	b.WriteString("              ║\n")
	b.WriteString(fmt.Sprintf("║  Max output: %d chars", r.MaxOutput))
	b.WriteString("             ║\n")
	b.WriteString(fmt.Sprintf("║  Web tools : "))

	if r.AllowWeb { b.WriteString("✅ Aktif") } else { b.WriteString("❌ Nonaktif") }
	b.WriteString("              ║\n")
	b.WriteString(fmt.Sprintf("║  Git tools : "))
	if r.AllowGit { b.WriteString("✅ Aktif") } else { b.WriteString("❌ Nonaktif") }
	b.WriteString("              ║\n")
	b.WriteString("╚══════════════════════════════════════════╝\n")
	return b.String()
}

func SystemPrompt() string {
	r := ActiveRules
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Kamu adalah AI coding agent yang cerdas dan membantu.

Gunakan Bahasa %s dalam setiap interaksi.

Aturan:
- Gunakan tools yang tersedia untuk menyelesaikan task
- Kerjakan langkah demi langkah
- Analisa project structure dulu sebelum memulai task besar
- Jika ada error, cari solusi pakai web_search
- Tanyakan jika ada yang kurang jelas
`, r.Bahasa))

	if r.AllowWeb {
		b.WriteString("- Kamu diizinkan mencari informasi di web\n")
	}
	if r.AllowGit {
		b.WriteString("- Kamu diizinkan menggunakan git (status, diff, commit)\n")
	}

	b.WriteString(fmt.Sprintf("\nMaksimal %d iterasi tool calling per task.\n", r.MaxIter))
	b.WriteString(fmt.Sprintf("Output tool dibatasi %d karakter.\n", r.MaxOutput))
	b.WriteString("\nTools yang tersedia:\n")

	defs := GetToolDefs()
	for _, d := range defs {
		b.WriteString(fmt.Sprintf("- %s: %s\n", d.Function.Name, d.Function.Description))
	}

	return b.String()
}
