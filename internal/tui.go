package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const sessionFile = ".agent_session.json"

type ChatMsg struct {
	Role    string
	Content string
}

type ToolEvent struct {
	ToolName string
	Args     string
	Result   string
}

type StreamEvent struct {
	Content string
	Done    bool
	Err     error
}

type model struct {
	prog         *tea.Program
	ready        bool
	viewport     viewport.Model
	input        string
	chatHistory  []ChatMsg
	messages     []Message
	agent        *Agent
	styles       Styles
	spinner      spinner.Model
	loading      bool
	width        int
	height       int
	streamBuffer strings.Builder
	toolResults  []ToolEvent
	statusMsg    string
	cursorPos    int
}

func (m *model) SetProgram(p *tea.Program) {
	m.prog = p
}

func InitialModel() model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	s.Spinner = spinner.Dot

	m := model{
		styles:  NewStyles(),
		spinner: s,
		agent:   NewAgent("."),
	}

	saved, _ := LoadSession(sessionFile)
	if len(saved) > 1 {
		m.messages = saved
		m.chatHistory = extractChatHistory(saved)
		m.statusMsg = "Sesi sebelumnya dilanjutkan"
	}

	return m
}

func extractChatHistory(msgs []Message) []ChatMsg {
	var history []ChatMsg
	for _, msg := range msgs {
		if msg.Role == "user" && msg.Content != "" {
			history = append(history, ChatMsg{Role: "user", Content: msg.Content})
		} else if msg.Role == "assistant" && msg.Content != "" {
			history = append(history, ChatMsg{Role: "assistant", Content: msg.Content})
		} else if msg.Role == "tool" && msg.Content != "" {
			trunc := msg.Content
			if len(trunc) > 100 {
				trunc = trunc[:100]
			}
			history = append(history, ChatMsg{Role: "tool", Content: "  ⚡ " + trunc + " ..."})
		}
	}
	return history
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, msg.Height-8)
			m.viewport.YPosition = 2
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 8
		}
		m.viewport.SetContent(m.renderChat())
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.input == "" {
				return m, nil
			}
			userMsg := m.input
			m.input = ""
			m.cursorPos = 0
			m.chatHistory = append(m.chatHistory, ChatMsg{Role: "user", Content: userMsg})
			m.messages = append(m.messages, Message{Role: "user", Content: userMsg})
			m.loading = true
			m.streamBuffer.Reset()
			m.toolResults = nil
			m.statusMsg = ""

			SaveSession(m.messages, sessionFile)

			go m.processUserMessage(userMsg)

			return m, m.spinner.Tick

		case "backspace":
			if m.cursorPos > 0 {
				m.input = m.input[:m.cursorPos-1] + m.input[m.cursorPos:]
				m.cursorPos--
			}
		case "delete":
			if m.cursorPos < len(m.input) {
				m.input = m.input[:m.cursorPos] + m.input[m.cursorPos+1:]
			}
		case "left":
			if m.cursorPos > 0 {
				m.cursorPos--
			}
		case "right":
			if m.cursorPos < len(m.input) {
				m.cursorPos++
			}
		case "home":
			m.cursorPos = 0
		case "end":
			m.cursorPos = len(m.input)
		case "ctrl+u":
			m.input = ""
			m.cursorPos = 0
		default:
			if len(msg.String()) == 1 {
				r := msg.Runes
				if len(r) > 0 {
					before := m.input[:m.cursorPos]
					after := m.input[m.cursorPos:]
					m.input = before + string(r) + after
					m.cursorPos += len(r)
				}
			}
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case StreamEvent:
		if msg.Err != nil {
			m.chatHistory = append(m.chatHistory, ChatMsg{Role: "error", Content: "❌ " + msg.Err.Error()})
			m.loading = false
			m.viewport.SetContent(m.renderChat())
			m.viewport.GotoBottom()
			return m, nil
		}
		if msg.Done {
			if m.streamBuffer.Len() > 0 {
				content := m.streamBuffer.String()
				m.chatHistory = append(m.chatHistory, ChatMsg{Role: "assistant", Content: content})
				m.streamBuffer.Reset()
			}
			for _, tr := range m.toolResults {
				truncArgs := tr.Args
				if len(truncArgs) > 80 {
					truncArgs = truncArgs[:80]
				}
				m.chatHistory = append(m.chatHistory, ChatMsg{Role: "tool", Content: fmt.Sprintf("  %s(%s)", tr.ToolName, truncArgs)})
				firstLine := strings.Split(tr.Result, "\n")[0]
				if len(firstLine) > 120 {
					firstLine = firstLine[:120]
				}
				m.chatHistory = append(m.chatHistory, ChatMsg{Role: "tool", Content: "  ✅ " + firstLine})
			}
			m.loading = false
			SaveSession(m.messages, sessionFile)
			if m.statusMsg == "" {
				m.statusMsg = "Selesai ✅"
			}
			m.viewport.SetContent(m.renderChat())
			m.viewport.GotoBottom()
			return m, nil
		}
		m.streamBuffer.WriteString(msg.Content)
		m.viewport.SetContent(m.renderChat())
		m.viewport.GotoBottom()
		return m, nil

	case ToolEvent:
		m.toolResults = append(m.toolResults, msg)
		m.statusMsg = fmt.Sprintf("🔧 %s...", msg.ToolName)
		m.viewport.SetContent(m.renderChat())
		return m, nil
	}

	return m, nil
}

func (m *model) processUserMessage(userMsg string) {
	onContent := func(content string) {
		m.prog.Send(StreamEvent{Content: content})
	}

	onTool := func(name, args, result string) {
		m.prog.Send(ToolEvent{ToolName: name, Args: args, Result: result})
	}

	err := m.agent.ProcessFull(onContent, onTool)

	done := StreamEvent{Done: true}
	if err != nil {
		done.Err = err
	}
	m.prog.Send(done)
}

func (m model) renderChat() string {
	var b strings.Builder
	for _, chat := range m.chatHistory {
		switch chat.Role {
		case "user":
			b.WriteString(m.styles.UserMsg.Render("You: "))
			b.WriteString(chat.Content + "\n\n")
		case "assistant":
			rendered, err := glamour.Render(chat.Content, "dark")
			if err != nil {
				b.WriteString(m.styles.AIMsg.Render(chat.Content + "\n\n"))
			} else {
				b.WriteString(rendered + "\n")
			}
		case "stream":
			b.WriteString(m.streamBuffer.String())
		case "tool":
			b.WriteString(m.styles.ToolMsg.Render(chat.Content) + "\n")
		case "error":
			b.WriteString(m.styles.ErrorMsg.Render(chat.Content) + "\n\n")
		}
	}
	return b.String()
}

func (m model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}

	modelName := lipgloss.NewStyle().Foreground(lipgloss.Color("#A5B4FC")).Render("DeepSeek-V3.1")
	status := "● Idle"
	if m.loading {
		status = m.spinner.View() + " Working..."
	}
	header := m.styles.Header.Width(m.width - 4).Render(fmt.Sprintf("  🚀 AI Agent  |  %s  |  %s", modelName, status))

	m.viewport.Width = m.width - 4
	m.viewport.Height = m.height - 8
	m.viewport.SetContent(m.renderChat())
	chatView := m.styles.Chat.Width(m.width - 4).Render(m.viewport.View())

	inputStyle := m.styles.Input.Width(m.width - 10)
	inputLine := inputStyle.Render(fmt.Sprintf(" %s█", m.input))
	charCount := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(fmt.Sprintf(" %d", len(m.input)))

	help := m.styles.Help.Width(m.width - 4).Render(" Ctrl+C: Keluar  |  Enter: Kirim  |  Ketik 'new' untuk sesi baru")

	sections := []string{header, "", chatView, "", inputLine + charCount, help}
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}
