package main

import (
	"fmt"
	"os"

	"github.com/gopartner/my-agent/internal"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	apiKey := os.Getenv("MY_AGENT_TOKEN")
	if apiKey == "" {
		fmt.Println("❌ MY_AGENT_TOKEN environment variable tidak ditemukan.")
		fmt.Println("   Set dengan: $env:MY_AGENT_TOKEN = \"hf_...\"")
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
