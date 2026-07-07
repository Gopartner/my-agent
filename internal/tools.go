package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetToolDefs() []ToolDef {
	return []ToolDef{
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "read_file",
				Description: "Baca isi file",
				Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string","description":"Path file"}},"required":["path"]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "write_file",
				Description: "Tulis file baru (folder otomatis dibuat)",
				Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"},"content":{"type":"string"}},"required":["path","content"]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "edit_file",
				Description: "Edit file dengan mencari dan mengganti teks",
				Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"},"old_string":{"type":"string"},"new_string":{"type":"string"}},"required":["path","old_string","new_string"]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "delete_file",
				Description: "Hapus file atau folder",
				Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"}},"required":["path"]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "list_dir",
				Description: "Lihat isi folder",
				Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string","description":"Path folder (default: .)"}},"required":[]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "project_tree",
				Description: "Lihat struktur project dalam bentuk tree",
				Parameters: json.RawMessage(`{"type":"object","properties":{},"required":[]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "run_command",
				Description: "Jalankan perintah shell (npm, pip, git, dll)",
				Parameters: json.RawMessage(`{"type":"object","properties":{"command":{"type":"string"}},"required":["command"]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "search_code",
				Description: "Cari teks dalam file project",
				Parameters: json.RawMessage(`{"type":"object","properties":{"pattern":{"type":"string"},"include":{"type":"string","description":"Glob pattern (contoh: *.go)"}},"required":["pattern"]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "git_status",
				Description: "Cek status git",
				Parameters: json.RawMessage(`{"type":"object","properties":{},"required":[]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "git_diff",
				Description: "Lihat perubahan yang belum di-commit",
				Parameters: json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"}},"required":[]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "git_commit",
				Description: "Commit semua perubahan",
				Parameters: json.RawMessage(`{"type":"object","properties":{"message":{"type":"string"}},"required":["message"]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "web_search",
				Description: "Cari informasi di web untuk dokumentasi/tutorial/solusi",
				Parameters: json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"}},"required":["query"]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "web_fetch",
				Description: "Ambil konten halaman web",
				Parameters: json.RawMessage(`{"type":"object","properties":{"url":{"type":"string"}},"required":["url"]}`),
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "http_request",
				Description: "Kirim HTTP request ke API",
				Parameters: json.RawMessage(`{"type":"object","properties":{"method":{"type":"string"},"url":{"type":"string"},"headers":{"type":"object"},"body":{"type":"string"}},"required":["method","url"]}`),
			},
		},
	}
}

func ExecuteTool(name string, args json.RawMessage, wd string) string {
	var result string
	switch name {
	case "read_file":
		var p struct{ Path string }
		json.Unmarshal(args, &p)
		result = doReadFile(p.Path, wd)
	case "write_file":
		var p struct{ Path, Content string }
		json.Unmarshal(args, &p)
		result = doWriteFile(p.Path, p.Content, wd)
	case "edit_file":
		var p struct{ Path, OldString, NewString string }
		json.Unmarshal(args, &p)
		result = doEditFile(p.Path, p.OldString, p.NewString, wd)
	case "delete_file":
		var p struct{ Path string }
		json.Unmarshal(args, &p)
		result = doDeleteFile(p.Path, wd)
	case "list_dir":
		var p struct{ Path string }
		json.Unmarshal(args, &p)
		if p.Path == "" { p.Path = "." }
		result = doListDir(p.Path, wd)
	case "project_tree":
		result = doProjectTree(wd)
	case "run_command":
		var p struct{ Command string }
		json.Unmarshal(args, &p)
		result = doRunCommand(p.Command, wd)
	case "search_code":
		var p struct{ Pattern, Include string }
		json.Unmarshal(args, &p)
		if p.Include == "" { p.Include = "**/*" }
		result = doSearchCode(p.Pattern, p.Include, wd)
	case "git_status":
		result = doRunCommand("git status", wd)
	case "git_diff":
		var p struct{ Path string }
		json.Unmarshal(args, &p)
		cmd := "git diff"
		if p.Path != "" { cmd += " " + p.Path }
		result = doRunCommand(cmd, wd)
	case "git_commit":
		var p struct{ Message string }
		json.Unmarshal(args, &p)
		r1 := doRunCommand("git add -A", wd)
		r2 := doRunCommand(fmt.Sprintf("git commit -m %s", escapeArg(p.Message)), wd)
		result = r1 + "\n" + r2
	case "web_search":
		var p struct{ Query string }
		json.Unmarshal(args, &p)
		result = doWebSearch(p.Query)
	case "web_fetch":
		var p struct{ URL string }
		json.Unmarshal(args, &p)
		result = doWebFetch(p.URL)
	case "http_request":
		var p struct {
			Method  string            `json:"method"`
			URL     string            `json:"url"`
			Headers map[string]string `json:"headers,omitempty"`
			Body    string            `json:"body,omitempty"`
		}
		json.Unmarshal(args, &p)
		result = doHTTPRequest(p.Method, p.URL, p.Headers, p.Body)
	default:
		result = fmt.Sprintf("Tool '%s' tidak dikenal", name)
	}
	if len(result) > 3000 {
		result = result[:3000] + "\n... (truncated)"
	}
	return result
}

func doReadFile(path, wd string) string {
	fp := filepath.Join(wd, path)
	data, err := os.ReadFile(fp)
	if err != nil { return fmt.Sprintf("ERROR: %v", err) }
	return string(data)
}

func doWriteFile(path, content, wd string) string {
	fp := filepath.Join(wd, path)
	os.MkdirAll(filepath.Dir(fp), 0755)
	if err := os.WriteFile(fp, []byte(content), 0644); err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return fmt.Sprintf("File ditulis: %s (%d bytes)", path, len(content))
}

func doEditFile(path, oldStr, newStr, wd string) string {
	fp := filepath.Join(wd, path)
	data, err := os.ReadFile(fp)
	if err != nil { return fmt.Sprintf("ERROR: %v", err) }
	content := string(data)
	if !strings.Contains(content, oldStr) {
		return fmt.Sprintf("ERROR: Teks tidak ditemukan di %s", path)
	}
	content = strings.ReplaceAll(content, oldStr, newStr)
	os.WriteFile(fp, []byte(content), 0644)
	return fmt.Sprintf("File diedit: %s", path)
}

func doDeleteFile(path, wd string) string {
	fp := filepath.Join(wd, path)
	if err := os.RemoveAll(fp); err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return fmt.Sprintf("Dihapus: %s", path)
}

func doListDir(path, wd string) string {
	fp := filepath.Join(wd, path)
	entries, err := os.ReadDir(fp)
	if err != nil { return fmt.Sprintf("ERROR: %v", err) }
	var b strings.Builder
	for _, e := range entries {
		prefix := "📄 "
		if e.IsDir() { prefix = "📁 " }
		b.WriteString(prefix + e.Name() + "\n")
	}
	return strings.TrimSpace(b.String())
}

func doProjectTree(wd string) string {
	var b strings.Builder
	filepath.Walk(wd, func(p string, fi os.FileInfo, err error) error {
		if err != nil { return nil }
		rel, _ := filepath.Rel(wd, p)
		if rel == "." { return nil }
		if strings.HasPrefix(rel, ".") && fi.IsDir() { return filepath.SkipDir }
		if strings.HasPrefix(filepath.Base(p), ".") { return nil }
		depth := strings.Count(rel, string(filepath.Separator))
		if depth > 4 { return nil }
		prefix := "  "
		if fi.IsDir() { prefix = "📁 " } else { prefix = "📄 " }
		b.WriteString(strings.Repeat("  ", depth) + prefix + filepath.Base(p) + "\n")
		return nil
	})
	return strings.TrimSpace(b.String())
}

func doRunCommand(cmdStr, wd string) string {
	var shell, flag string
	if ShellIsWindows() {
		shell, flag = "cmd", "/C"
	} else {
		shell, flag = "sh", "-c"
	}
	cmd := exec.Command(shell, flag, cmdStr)
	cmd.Dir = wd
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out) + "\nERROR: " + err.Error()
	}
	return string(out)
}

func doSearchCode(pattern, include, wd string) string {
	var b strings.Builder
	filepath.Walk(wd, func(p string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() { return nil }
		rel, _ := filepath.Rel(wd, p)
		if strings.HasPrefix(filepath.Base(p), ".") { return nil }
		if include != "**/*" {
			m, _ := filepath.Match(include, filepath.Base(p))
			if !m { return nil }
		}
		data, err := os.ReadFile(p)
		if err != nil { return nil }
		for i, line := range strings.Split(string(data), "\n") {
			if strings.Contains(line, pattern) {
				fmt.Fprintf(&b, "%s:%d: %s\n", rel, i+1, strings.TrimSpace(line)[:min(len(line), 200)])
			}
		}
		return nil
	})
	if b.Len() == 0 { return "Tidak ditemukan" }
	return strings.TrimSpace(b.String()[:min(b.Len(), 4000)])
}

func doWebSearch(query string) string {
	url := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", strings.ReplaceAll(query, " ", "+"))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return fmt.Sprintf("ERROR: %v", err) }
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil { return fmt.Sprintf("ERROR: %v", err) }

	var b strings.Builder
	doc.Find(".result").Each(func(i int, s *goquery.Selection) {
		if i >= 8 { return }
		title := strings.TrimSpace(s.Find(".result__title").Text())
		snippet := strings.TrimSpace(s.Find(".result__snippet").Text())
		if title != "" {
			fmt.Fprintf(&b, "%d. %s\n   %s\n\n", i+1, title, snippet[:min(len(snippet), 200)])
		}
	})
	if b.Len() == 0 { return "Tidak ada hasil" }
	return strings.TrimSpace(b.String())
}

func doWebFetch(url string) string {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return fmt.Sprintf("ERROR: %v", err) }
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil { return fmt.Sprintf("ERROR: %v", err) }
	doc.Find("script, style, nav, footer, header").Remove()
	text := strings.TrimSpace(doc.Text())
	if len(text) > 5000 { text = text[:5000] + "\n... (truncated)" }
	return text
}

func doHTTPRequest(method, url string, headers map[string]string, body string) string {
	var reqBody io.Reader
	if body != "" { reqBody = bytes.NewReader([]byte(body)) }
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil { return fmt.Sprintf("ERROR: %v", err) }
	req.Header.Set("User-Agent", "Mozilla/5.0")
	for k, v := range headers { req.Header.Set(k, v) }
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return fmt.Sprintf("ERROR: %v", err) }
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 3000))
	return fmt.Sprintf("Status: %d\nBody:\n%s", resp.StatusCode, string(b))
}

func escapeArg(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

func ShellIsWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

func min(a, b int) int {
	if a < b { return a }
	return b
}


