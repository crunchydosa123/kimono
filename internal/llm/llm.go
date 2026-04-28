package llm

import (
	"context"

	"github.com/crunchydosa123/kimono/internal/tool"
)

type Message struct {
	Role    string // "system", "user", "assistant"
	Content string
}

type LLM interface {
	ChatCompletion(ctx context.Context, messages []Message, tools []tool.Tool) (LLMResponse, error)
	GeneratePlan(ctx context.Context, message Message) (LLMResponse, error)
}
