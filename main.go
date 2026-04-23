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
			Content: "Create a Go file hello.go with a hello world program.",
		},
	}

	// ⚠️ make sure your interface supports tools
	res, err := model.ChatCompletion(
		context.Background(),
		messages,
		tools, // 👈 pass tools
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
