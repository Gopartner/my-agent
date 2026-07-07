package internal

import (
	"encoding/json"
	"fmt"
)

type Agent struct {
	Client   *HFClient
	Messages []Message
	WD       string
}

func NewAgent(wd string) *Agent {
	return &Agent{
		Client: NewHFClient(),
		Messages: []Message{
			{
				Role: "system",
				Content: `Kamu adalah AI coding agent yang sangat cerdas dan membantu.

Kemampuanmu:
1. Membaca/menulis/mengedit/menghapus file
2. Menjalankan command terminal (npm, pip, git, dll)
3. Git operations (status, diff, commit)
4. Mencari kode dan menganalisa struktur project
5. Search web dan fetch dokumentasi
6. HTTP requests ke API

Aturan:
- Gunakan tools yang tersedia untuk menyelesaikan task
- Kerjakan langkah demi langkah
- Analisa project structure dulu sebelum memulai task besar
- Gunakan Bahasa Indonesia
- Jika ada error, cari solusi dari web
- Tanyakan jika ada yang kurang jelas`,
			},
		},
		WD: wd,
	}
}

type ToolResult struct {
	Role       string `json:"role"`
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
}

type StreamState struct {
	Content   string
	ToolCalls []ToolCall
	Done      bool
	Finish    string
}

type ToolCallAccum struct {
	toolCalls map[int]*ToolCall
}

func NewToolCallAccum() *ToolCallAccum {
	return &ToolCallAccum{toolCalls: make(map[int]*ToolCall)}
}

func (a *ToolCallAccum) Add(tc ToolCall) {
	idx := tc.Index
	if existing, ok := a.toolCalls[idx]; ok {
		if tc.ID != "" {
			existing.ID += tc.ID
		}
		if tc.Function.Name != "" {
			existing.Function.Name += tc.Function.Name
		}
		if tc.Function.Arguments != "" {
			existing.Function.Arguments += tc.Function.Arguments
		}
	} else {
		call := ToolCall{
			ID:   tc.ID,
			Type: "function",
		}
		call.Function.Name = tc.Function.Name
		call.Function.Arguments = tc.Function.Arguments
		a.toolCalls[idx] = &call
	}
}

func (a *ToolCallAccum) GetList() []ToolCall {
	var result []ToolCall
	for i := 0; i < len(a.toolCalls); i++ {
		if tc, ok := a.toolCalls[i]; ok {
			result = append(result, *tc)
		}
	}
	return result
}

// ProcessRun menjalankan satu iterasi: send messages -> stream response -> handle tools
// Returns: content, toolResults, done, error
func (ag *Agent) ProcessRun() (string, []ToolCall, []ToolResult, bool, error) {
	var streamState StreamState
	tcAccum := NewToolCallAccum()

	req := ChatRequest{
		Model:      "deepseek-ai/DeepSeek-V3.1",
		Messages:   ag.Messages,
		Tools:      GetToolDefs(),
		ToolChoice: "auto",
		MaxTokens:  4096,
	}

	var toolCallsFinal []ToolCall

	err := ag.Client.ChatStream(req,
		func(content string) {
			streamState.Content += content
		},
		func(tc ToolCall) {
			tcAccum.Add(tc)
		},
		func(finish string) {
			streamState.Done = true
			streamState.Finish = finish
			toolCallsFinal = tcAccum.GetList()
		},
	)
	if err != nil {
		return "", nil, nil, false, fmt.Errorf("API error: %w", err)
	}

	streamState.ToolCalls = toolCallsFinal

	// Add assistant message to history
	assistantMsg := Message{Role: "assistant", Content: streamState.Content}
	if len(streamState.ToolCalls) > 0 {
		jsonTC, _ := json.Marshal(streamState.ToolCalls)
		var tcList []ToolCall
		json.Unmarshal(jsonTC, &tcList)
		assistantMsg.ToolCalls = tcList
	}
	ag.Messages = append(ag.Messages, assistantMsg)

	// Execute tools
	var toolResults []ToolResult
	for _, tc := range streamState.ToolCalls {
		var args json.RawMessage
		if tc.Function.Arguments != "" {
			args = json.RawMessage(tc.Function.Arguments)
		}
		result := ExecuteTool(tc.Function.Name, args, ag.WD)

		tr := ToolResult{
			Role:       "tool",
			ToolCallID: tc.ID,
			Content:    result,
		}
		toolResults = append(toolResults, tr)

		ag.Messages = append(ag.Messages, Message{
			Role:       "tool",
			ToolCallID: tc.ID,
			Content:    result,
		})
	}

	hasTools := len(streamState.ToolCalls) > 0

	if streamState.Content != "" {
		return streamState.Content, streamState.ToolCalls, toolResults, !hasTools, nil
	}
	return streamState.Content, streamState.ToolCalls, toolResults, !hasTools, nil
}

// ProcessFull menjalankan loop hingga tidak ada tool calls lagi (max 15 iterasi)
func (ag *Agent) ProcessFull(onContent func(string), onToolCall func(string, string, string)) error {
	for i := 0; i < 15; i++ {
		content, tcs, results, done, err := ag.ProcessRun()
		if err != nil {
			onContent(fmt.Sprintf("\n❌ %v\n", err))
			return err
		}
		if content != "" {
			onContent(content)
		}
		for i, tc := range tcs {
			onToolCall(tc.Function.Name, tc.Function.Arguments, results[i].Content)
		}
		if done {
			break
		}
	}
	return nil
}
