package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SessionData struct {
	Messages []Message `json:"messages"`
}

func SaveSession(messages []Message, path string) error {
	os.MkdirAll(filepath.Dir(path), 0755)
	data := SessionData{Messages: messages}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func LoadSession(path string) ([]Message, error) {
	data, err := os.ReadFile(path)
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
