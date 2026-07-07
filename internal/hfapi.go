package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const HFBase = "https://router.huggingface.co/v1"

type Message struct {
	Role       string      `json:"role"`
	Content    string      `json:"content,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

type ToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type ToolDef struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
	Index int `json:"index,omitempty"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Tools       []ToolDef `json:"tools,omitempty"`
	ToolChoice  string    `json:"tool_choice,omitempty"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
}

type ChatChoice struct {
	Delta struct {
		Content    string     `json:"content,omitempty"`
		ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	} `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

type ChatStreamResponse struct {
	Choices []ChatChoice `json:"choices"`
}

type HFClient struct {
	apiKey string
	client *http.Client
}

func NewHFClient() *HFClient {
	return &HFClient{
		apiKey: os.Getenv("MY_AGENT_TOKEN"),
		client: &http.Client{},
	}
}

func (c *HFClient) ChatStream(req ChatRequest, onContent func(string), onToolCall func(ToolCall), onDone func(string)) error {
	if req.MaxTokens == 0 {
		req.MaxTokens = 4096
	}
	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	httpReq, err := http.NewRequest("POST", HFBase+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API %d: %s", resp.StatusCode, string(b))
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024*64), 1024*1024)
	var finish string

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var sresp ChatStreamResponse
		if err := json.Unmarshal([]byte(data), &sresp); err != nil {
			continue
		}
		if len(sresp.Choices) == 0 {
			continue
		}
		ch := sresp.Choices[0]
		if ch.Delta.Content != "" {
			onContent(ch.Delta.Content)
		}
		for _, tc := range ch.Delta.ToolCalls {
			onToolCall(tc)
		}
		if ch.FinishReason != "" {
			finish = ch.FinishReason
		}
	}

	onDone(finish)
	return scanner.Err()
}
