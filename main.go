package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/crunchydosa123/kimono/agent"
	"github.com/crunchydosa123/kimono/internal/llm"
	"github.com/crunchydosa123/kimono/internal/tool"
)

func main() {
	var model llm.LLM

	gemini, err := llm.NewGemini(os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	Content :=
		`
You are a coding agent.

Available tools:
- write_file: create or overwrite files
- edit_file: modify files

Rules:
- Only use the tools listed above
- Do NOT call any other tools
`

	model = gemini

	registry := tool.NewRegistry()
	registry.Register(&tool.WriteFile{})
	registry.Register(&tool.EditFile{})

	messages := []llm.Message{
		{
			Role:    "system",
			Content: "You are a coding assistant that can use tools." + Content,
		},
		{
			Role:    "user",
			Content: "Edit hello.go, understand what it does and without editing its main functionality fix any errors you see pertaining to syntax",
		},
	}

	// make sure your interface supports tools
	ctx := context.Background()

	agent := agent.New(model, registry)

	res, err := agent.Run(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Final response:", res)
}
