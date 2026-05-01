package main

import (
	"context"
	"log"
	"os"

	"github.com/crunchydosa123/kimono/agent"
	chatsession "github.com/crunchydosa123/kimono/chat_session"
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

You can suggest terminal commands using the suggest_command tool.

Rules:
- DO NOT assume commands are executed automatically
- Only suggest commands when needed
- Prefer safe commands like:
  - go build
  - go run
  - ls
`

	model = gemini

	registry := tool.NewRegistry()
	registry.Register(&tool.WriteFile{})
	registry.Register(&tool.EditFile{})
	registry.Register(&tool.SearchCode{})
	registry.Register(&tool.SuggestCommand{})

	messages := []llm.Message{
		{
			Role:    "system",
			Content: "You are a coding assistant that can use tools." + Content,
		},
		{
			Role:    "user",
			Content: "give me the list of all files in internal",
		},
	}

	// make sure your interface supports tools
	ctx := context.Background()

	agent := agent.New(model, registry)

	chat := chatsession.NewChatSession(agent, messages)

	chat.Start(ctx)
}
