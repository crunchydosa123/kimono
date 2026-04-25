package agent

import (
	"context"
	"fmt"

	"github.com/crunchydosa123/kimono/internal/llm"
	"github.com/crunchydosa123/kimono/internal/tool"
)

type Agent struct {
	model    llm.LLM
	registry *tool.Registry
	maxSteps int
}

func New(model llm.LLM, registry *tool.Registry) *Agent {
	return &Agent{
		model:    model,
		registry: registry,
		maxSteps: 10, // safety limit
	}
}

func (a *Agent) Run(
	ctx context.Context,
	messages []llm.Message,
) (llm.LLMResponse, error) {

	tools := a.registry.All()

	for step := 0; step < a.maxSteps; step++ {
		res, err := a.model.ChatCompletion(ctx, messages, tools)
		if err != nil {
			return llm.LLMResponse{}, err
		}

		c := res.Candidates[0]

		var toolExecuted bool

		for _, p := range c.Parts {
			// print text
			if p.Text != nil {
				fmt.Println(*p.Text)
			}

			// execute tool
			if p.ToolCall != nil {
				toolExecuted = true

				fmt.Printf("🔧 Executing: %s %s\n",
					p.ToolCall.Name,
					string(p.ToolCall.Args),
				)

				result, err := a.registry.Execute(ctx, p.ToolCall.Name, p.ToolCall.Args)
				if err != nil {
					fmt.Println("Tool error:", err)
					continue
				}

				// append tool result as user message (Gemini-friendly)
				messages = append(messages, llm.Message{
					Role: "user",
					Content: fmt.Sprintf(
						"Tool %s result:\n%s",
						p.ToolCall.Name,
						result,
					),
				})
			}
		}

		// stop condition
		if !toolExecuted || c.FinishReason == "stop" {
			return res, nil
		}
	}

	return llm.LLMResponse{}, fmt.Errorf("max steps exceeded")
}
