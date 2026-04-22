package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"google.golang.org/genai"
)

var geminiClient *genai.Client

func init() {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		panic(err)
	}
	geminiClient = client
}

//
// ===== EVENTS =====
//

type EventType string

const (
	EventLog  EventType = "log"
	EventDiff EventType = "diff"
	EventDone EventType = "done"
)

type Event struct {
	Type EventType
	Data string
}

//
// ===== AGENT RESPONSE =====
//

type AgentResponse struct {
	Action string `json:"action"` // tool | edit | final
	Tool   string `json:"tool"`
	Input  string `json:"input"`

	File    string `json:"file"`
	Find    string `json:"find"`
	Replace string `json:"replace"`
}

//
// ===== TOOLS =====
//

func runTool(name, input string) string {
	switch name {
	case "grep":
		return "found 'x' in main.go"
	case "read_file":
		return "let x = a + b"
	default:
		return "unknown tool"
	}
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 {
		return s
	}
	return s[start : end+1]
}

//
// ===== MOCK LLM =====
// Replace this later with real API
//

func callLLM(prompt string) (AgentResponse, error) {
	ctx := context.Background()

	systemPrompt := `
You are a coding agent.

Respond ONLY in JSON.

Schema:
{
  "action": "tool" | "edit" | "final",
  "tool": "grep" | "read_file",
  "input": "string",
  "file": "string",
  "find": "string",
  "replace": "string"
}
`

	fullPrompt := systemPrompt + "\n\n" + prompt

	resp, err := geminiClient.Models.GenerateContent(
		ctx,
		"gemini-1.5-flash",
		genai.Text(fullPrompt),
		nil,
	)
	if err != nil {
		return AgentResponse{}, err
	}

	// Extract text
	var content string
	for _, cand := range resp.Candidates {
		if cand.Content == nil {
			continue
		}
		for _, part := range cand.Content.Parts {
			content += part.Text
		}
	}

	// Clean (Gemini sometimes adds markdown)
	content = extractJSON(content)

	var parsed AgentResponse
	err = json.Unmarshal([]byte(content), &parsed)
	if err != nil {
		return AgentResponse{}, err
	}

	return parsed, nil
}

//
// ===== AGENT LOOP =====
//

func runAgent(events chan Event, prompt string) {
	state := ""
	maxSteps := 5

	for step := 0; step < maxSteps; step++ {

		resp, err := callLLM(prompt + state)
		if err != nil {
			events <- Event{Type: EventLog, Data: "LLM error: " + err.Error()}
			return
		}

		switch resp.Action {

		case "tool":
			events <- Event{
				Type: EventLog,
				Data: fmt.Sprintf("%s(%s)", resp.Tool, resp.Input),
			}

			output := runTool(resp.Tool, resp.Input)
			state += "\n" + output

			time.Sleep(800 * time.Millisecond)

		case "edit":
			diff := generateDiff(resp.Find, resp.Replace)

			events <- Event{
				Type: EventDiff,
				Data: diff,
			}

			return

		case "final":
			events <- Event{Type: EventDone}
			return
		}
	}
}

//
// ===== DIFF =====
//

func generateDiff(before, after string) string {
	var b strings.Builder
	b.WriteString("--- before\n")
	b.WriteString(before + "\n")
	b.WriteString("+++ after\n")
	b.WriteString(after + "\n")
	return b.String()
}

//
// ===== TUI MODEL =====
//

type model struct {
	prompt string
	steps  []string
	diff   string
	done   bool

	events chan Event
}

func initialModel(prompt string) model {
	return model{
		prompt: prompt,
		steps:  []string{},
		diff:   "",
		done:   false,
		events: make(chan Event),
	}
}

func (m model) Init() tea.Cmd {
	return waitForEvent(m.events)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case Event:
		switch msg.Type {
		case EventLog:
			m.steps = append(m.steps, msg.Data)
		case EventDiff:
			m.diff = msg.Data
		case EventDone:
			m.done = true
		}
		return m, waitForEvent(m.events)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "a":
			m.steps = append(m.steps, "✅ Applied changes")
			m.done = true

		case "s":
			m.steps = append(m.steps, "⏭ Skipped changes")
			m.done = true
		}
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Prompt: %s\n\n", m.prompt)

	s += "Steps:\n"
	for _, step := range m.steps {
		s += "• " + step + "\n"
	}

	s += "\nDiff Preview:\n"
	if m.diff == "" {
		s += "(waiting...)\n"
	} else {
		s += m.diff + "\n"
	}

	s += "\n[a] Apply   [s] Skip   [q] Quit\n"

	if m.done {
		s += "\n✔ Done\n"
	}

	return s
}

func waitForEvent(events chan Event) tea.Cmd {
	return func() tea.Msg {
		return <-events
	}
}

//
// ===== MAIN =====
//

func main() {
	prompt := "rename x to total"

	m := initialModel(prompt)

	go runAgent(m.events, prompt)

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
	}
}
