package chatsession

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/crunchydosa123/kimono/agent"
	"github.com/crunchydosa123/kimono/internal/llm"
)

type ChatSession struct {
	agent    *agent.Agent
	messages []llm.Message
}

func NewChatSession(agent *agent.Agent, m []llm.Message) *ChatSession {
	msgs := make([]llm.Message, len(m))
	copy(msgs, m)

	return &ChatSession{
		agent:    agent,
		messages: msgs,
	}
}

func (c *ChatSession) Start(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		// -------------------------
		// 1. Generate plan
		// -------------------------
		plan, err := c.agent.GeneratePlan(ctx, llm.Message{
			Role:    "user",
			Content: input,
		})
		if err != nil {
			fmt.Println("Plan error:", err)
			continue
		}

		fmt.Println("\n🧠 Plan:")
		fmt.Println(plan)

		// -------------------------
		// 2. Enhance query with plan
		// -------------------------
		enhanced := fmt.Sprintf(`
User request:
%s

Execution plan:
%s

Follow this plan step-by-step. Do not repeat steps or mix strategies.
`, input, plan)

		// add enhanced message instead of raw input
		c.messages = append(c.messages, llm.Message{
			Role:    "user",
			Content: enhanced,
		})

		// -------------------------
		// 3. Run agent
		// -------------------------
		res, err := c.agent.Run(ctx, c.messages)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// extract assistant text
		var reply string
		for _, p := range res.Candidates[0].Parts {
			if p.Text != nil {
				reply += *p.Text + "\n"
			}
		}

		fmt.Println("\n🤖 Answer:")
		fmt.Println(reply)

		// -------------------------
		// 4. Store assistant reply
		// -------------------------
		c.messages = append(c.messages, llm.Message{
			Role:    "assistant",
			Content: reply,
		})
	}
}
