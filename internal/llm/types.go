package llm

import "encoding/json"

type FinishReason string

const (
	FinishStop     FinishReason = "stop"
	FinishToolCall FinishReason = "tool_call"
	FinishLength   FinishReason = "length"
)

type ToolCall struct {
	Name string
	Args json.RawMessage
}

type Part struct {
	Text     *string
	ToolCall *ToolCall
}

type Candidate struct {
	Role         string
	Parts        []Part
	FinishReason string // "stop", "tool_call", "length", etc.
}

type LLMResponse struct {
	Candidates []Candidate
}
