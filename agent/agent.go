package agent

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/crunchydosa123/kimono/internal/llm"
	"github.com/crunchydosa123/kimono/internal/tool"
)

type Agent struct {
	model      llm.LLM
	registry   *tool.Registry
	maxSteps   int
	confirmCmd func(string) (string, error)
}

func New(model llm.LLM, registry *tool.Registry) *Agent {
	return &Agent{
		model:      model,
		registry:   registry,
		maxSteps:   10, // safety limit
		confirmCmd: defaultConfirmCmd,
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

				var result string
				var err error

				if strings.ToLower(p.ToolCall.Name) == "suggest_command" {
					// get command from tool
					cmd, err := a.registry.Execute(ctx, p.ToolCall.Name, p.ToolCall.Args)
					if err != nil {
						fmt.Println("Tool error:", err)
						continue
					}

					fmt.Println("🚨 INTERCEPTING suggest_command")

					// 👇 THIS triggers your TUI
					result, err = a.confirmCmd(cmd)
					if err != nil {
						fmt.Println("Command error:", err)
						continue
					}
				} else {
					result, err = a.registry.Execute(ctx, p.ToolCall.Name, p.ToolCall.Args)
					if err != nil {
						fmt.Println("Tool error:", err)
						continue
					}
				}
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

func defaultConfirmCmd(cmd string) (string, error) {
	fmt.Printf("\n⚡ Suggested command:\n%s\n", cmd)
	fmt.Print("Run this? (y/n): ")

	var input string
	fmt.Scanln(&input)

	if input != "y" {
		return "command skipped by user", nil
	}

	parts := strings.Fields(cmd)
	c := exec.Command(parts[0], parts[1:]...)
	out, err := c.CombinedOutput()

	return string(out), err
}
