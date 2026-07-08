package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SessionData struct {
	Messages  []Message `json:"messages"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Count     int       `json:"message_count"`
}

func sessionDir() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		home, _ := os.UserHomeDir()
		appData = filepath.Join(home, ".config")
	}
	dir := filepath.Join(appData, "my-agent", "sessions")
	os.MkdirAll(dir, 0755)
	return dir
}

func sessionPath() string {
	return filepath.Join(sessionDir(), "current.json")
}

type SessionSummary struct {
	TotalMessages int
	TotalTools    int
	Duration      string
}

func SaveSession(messages []Message) error {
	now := time.Now().Format(time.RFC3339)
	data := SessionData{
		Messages:  messages,
		UpdatedAt: now,
		Count:     len(messages),
	}

	// Set created_at if first save
	if existing, err := os.ReadFile(sessionPath()); err == nil {
		var old SessionData
		if json.Unmarshal(existing, &old) == nil && old.CreatedAt != "" {
			data.CreatedAt = old.CreatedAt
		}
	}
	if data.CreatedAt == "" {
		data.CreatedAt = now
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sessionPath(), b, 0644)
}

func LoadSession() ([]Message, error) {
	data, err := os.ReadFile(sessionPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var session SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return session.Messages, nil
}

func EndSession() (SessionSummary, error) {
	data, err := os.ReadFile(sessionPath())
	if err != nil {
		if os.IsNotExist(err) {
			return SessionSummary{}, nil
		}
		return SessionSummary{}, err
	}

	var session SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return SessionSummary{}, err
	}

	created, _ := time.Parse(time.RFC3339, session.CreatedAt)
	duration := time.Since(created).Round(time.Second).String()

	toolCount := 0
	for _, m := range session.Messages {
		if m.Role == "tool" {
			toolCount++
		}
	}

	summary := SessionSummary{
		TotalMessages: session.Count,
		TotalTools:    toolCount,
		Duration:      duration,
	}

	// Archive current session
	archiveName := fmt.Sprintf("session-%s.json", time.Now().Format("2006-01-02-150405"))
	archivePath := filepath.Join(sessionDir(), archiveName)
	os.Rename(sessionPath(), archivePath)

	return summary, nil
}

func SessionSummaryText(summary SessionSummary) string {
	var b strings.Builder
	b.WriteString("╔══════════════════════════════════════════╗\n")
	b.WriteString("║            Sesi Selesai                 ║\n")
	b.WriteString("╠══════════════════════════════════════════╣\n")
	b.WriteString(fmt.Sprintf("║  Pesan     : %d", summary.TotalMessages))
	for i := 0; i < 22-len(fmt.Sprint(summary.TotalMessages)); i++ { b.WriteString(" ") }
	b.WriteString("║\n")
	b.WriteString(fmt.Sprintf("║  Tool call : %d", summary.TotalTools))
	for i := 0; i < 22-len(fmt.Sprint(summary.TotalTools)); i++ { b.WriteString(" ") }
	b.WriteString("║\n")
	b.WriteString(fmt.Sprintf("║  Durasi    : %s", summary.Duration))
	for i := 0; i < 22-len(summary.Duration); i++ { b.WriteString(" ") }
	b.WriteString("║\n")
	b.WriteString("╚══════════════════════════════════════════╝\n")
	return b.String()
}

func NewSession() {
	os.Remove(sessionPath())
}
