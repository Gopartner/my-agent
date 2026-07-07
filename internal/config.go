package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Token string `json:"token"`
}

func configDir() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		home, _ := os.UserHomeDir()
		appData = filepath.Join(home, ".config")
	}
	return filepath.Join(appData, "my-agent")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func LoadToken() string {
	if t := os.Getenv("MY_AGENT_TOKEN"); t != "" {
		return t
	}
	data, err := os.ReadFile(configPath())
	if err != nil {
		return ""
	}
	var cfg Config
	if json.Unmarshal(data, &cfg) != nil {
		return ""
	}
	return cfg.Token
}

func SaveToken(token string) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	cfg := Config{Token: token}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0600)
}

func PromptToken() string {
	fmt.Println("")
	fmt.Println("  ╔══════════════════════════════════════════╗")
	fmt.Println("  ║          Selamat datang di my-agent!     ║")
	fmt.Println("  ╠══════════════════════════════════════════╣")
	fmt.Println("  ║  Masukkan Hugging Face token kamu.       ║")
	fmt.Println("  ║  Token bisa dibuat di:                   ║")
	fmt.Println("  ║  https://huggingface.co/settings/tokens  ║")
	fmt.Println("  ║                                          ║")
	fmt.Println("  ║  Token akan disimpan di:                 ║")
	fmt.Println("  ║  " + configPath() + "   ║")
	fmt.Println("  ╚══════════════════════════════════════════╝")
	fmt.Print("  Token: ")

	var token string
	fmt.Scanln(&token)

	if token == "" {
		fmt.Println("  Token tidak boleh kosong.")
		os.Exit(1)
	}

	if err := SaveToken(token); err != nil {
		fmt.Printf("  Gagal menyimpan token: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("  ✔ Token tersimpan! Memulai my-agent...")
	fmt.Println("")
	return token
}
