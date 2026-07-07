package internal

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Header   lipgloss.Style
	Sidebar  lipgloss.Style
	Chat     lipgloss.Style
	Input    lipgloss.Style
	Help     lipgloss.Style
	Spinner  lipgloss.Style
	UserMsg  lipgloss.Style
	AIMsg    lipgloss.Style
	ToolMsg  lipgloss.Style
	ErrorMsg lipgloss.Style
	Divider  lipgloss.Style
}

func NewStyles() Styles {
	accent := lipgloss.Color("#7C3AED")
	success := lipgloss.Color("#10B981")
	err := lipgloss.Color("#EF4444")
	dim := lipgloss.Color("#6B7280")
	bg := lipgloss.Color("#0F172A")
	surface := lipgloss.Color("#1E293B")
	border := lipgloss.Color("#334155")
	text := lipgloss.Color("#E2E8F0")
	muted := lipgloss.Color("#94A3B8")

	return Styles{
		Header: lipgloss.NewStyle().
			Background(accent).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 2),

		Sidebar: lipgloss.NewStyle().
			Background(surface).
			Border(lipgloss.NormalBorder()).
			BorderRight(true).
			BorderForeground(border).
			Foreground(text).
			Padding(1, 1),

		Chat: lipgloss.NewStyle().
			Background(bg).
			Foreground(text).
			Padding(1, 2),

		Input: lipgloss.NewStyle().
			Background(surface).
			Border(lipgloss.NormalBorder()).
			BorderForeground(border).
			Padding(0, 1).
			Foreground(text),

		Help: lipgloss.NewStyle().
			Foreground(muted).
			Padding(0, 2),

		Spinner: lipgloss.NewStyle().
			Foreground(accent),

		UserMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A5B4FC")).
			Bold(true),

		AIMsg: lipgloss.NewStyle().
			Foreground(text).
			Padding(0, 0, 1, 0),

		ToolMsg: lipgloss.NewStyle().
			Foreground(success).
			Italic(true),

		ErrorMsg: lipgloss.NewStyle().
			Foreground(err).
			Bold(true),

		Divider: lipgloss.NewStyle().
			Foreground(dim).
			Padding(0, 1),
	}
}
