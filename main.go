package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/crunchydosa123/kimono/internal/llm"
	"github.com/crunchydosa123/kimono/internal/tool"
)

func main() {
	var model llm.LLM

	gemini, err := llm.NewGemini(os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	model = gemini

	registry := tool.NewRegistry()
	registry.Register(&tool.WriteFile{})

	tools := registry.All()

	messages := []llm.Message{
		{
			Role:    "system",
			Content: "You are a coding assistant that can use tools.",
		},
		{
			Role:    "user",
			Content: "Create a Go file hello.go with a 2 go routines. the package will be test and the go routines will print something",
		},
	}

	// make sure your interface supports tools
	ctx := context.Background()

	for {
		res, err := model.ChatCompletion(ctx, messages, tools)
		if err != nil {
			log.Fatal(err)
		}

		// assume single candidate for now
		c := res.Candidates[0]

		var toolExecuted bool

		for _, p := range c.Parts {
			// print text if present
			if p.Text != nil {
				fmt.Println(*p.Text)
			}

			// execute tool if present
			if p.ToolCall != nil {
				toolExecuted = true

				fmt.Println("🔧 Executing:", p.ToolCall.Name)

				result, err := registry.Execute(p.ToolCall.Name, p.ToolCall.Args)
				if err != nil {
					log.Fatal(err)
				}

				// append assistant tool call
				messages = append(messages, llm.Message{
					Role:    "assistant",
					Content: "", // optional
				})

				// append tool result back
				messages = append(messages, llm.Message{
					Role:    "tool",
					Content: result, // string or JSON
				})
			}
		}

		// stop if no tool calls
		if !toolExecuted || c.FinishReason == "stop" {
			break
		}
	}
}
