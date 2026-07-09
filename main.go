package main

import (
	"fmt"
	"os"

	"github.com/gopartner/my-agent/internal"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	token := internal.LoadToken()
	if token == "" {
		token = internal.PromptToken()
	}
	os.Setenv("MY_AGENT_TOKEN", token)

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
