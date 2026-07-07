package main

import (
	"fmt"
	"os"

	"github.com/gopartner/my-agent/internal"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	apiKey := os.Getenv("HF_TOKEN")
	if apiKey == "" {
		fmt.Println("❌ HF_TOKEN environment variable tidak ditemukan.")
		fmt.Println("   Set dengan: $env:HF_TOKEN = \"hf_...\"")
		os.Exit(1)
	}

	m := internal.InitialModel()
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	m.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
