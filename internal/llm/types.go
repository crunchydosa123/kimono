package llm

type Response struct {
	Text     string
	ToolCall *ToolCall
}

type ToolCall struct {
	Name string
	Args map[string]any
}
