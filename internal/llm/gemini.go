package llm

import (
	"context"
	"fmt"

	"github.com/crunchydosa123/kimono/internal/tool"
	"google.golang.org/genai"
)

type Gemini struct {
	client *genai.Client
	model  string
}

func NewGemini(apiKey string) (*Gemini, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})

	if err != nil {
		return nil, err
	}

	return &Gemini{
		client: client,
		model:  "gemini-flash-latest",
	}, nil
}

func (g *Gemini) ChatCompletion(
	ctx context.Context,
	messages []Message,
	tools []tool.Tool,
) (string, error) {

	var contents []*genai.Content

	for _, m := range messages {
		role := "user"
		if m.Role == "assistant" {
			role = "model"
		}

		contents = append(contents, &genai.Content{
			Role: role,
			Parts: []*genai.Part{
				{Text: m.Content},
			},
		})
	}

	// ✅ attach tools
	config := &genai.GenerateContentConfig{}
	if len(tools) > 0 {
		config.Tools = tool.ToGeminiTools(tools)
	}

	resp, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		contents,
		config, // 👈 pass config instead of nil
	)
	if err != nil {
		return "", err
	}

	fmt.Println("---- GEMINI RESPONSE ----")

	for i, c := range resp.Candidates {
		fmt.Printf("Candidate %d:\n", i)

		if c.Content == nil {
			fmt.Println("  ❌ No content")
			continue
		}

		fmt.Println("  Role:", c.Content.Role)

		for j, p := range c.Content.Parts {
			fmt.Printf("  Part %d:\n", j)

			if p.Text != "" {
				fmt.Println("    Text:", p.Text)
			}

			if p.FunctionCall != nil {
				fmt.Println("    🔥 Function Call:")
				fmt.Println("      Name:", p.FunctionCall.Name)
				fmt.Println("      Args:", p.FunctionCall.Args)
			}

			if p.FunctionResponse != nil {
				fmt.Println("    Function Response:", p.FunctionResponse)
			}
		}
	}

	fmt.Println("--------------------------")
	// ⚠️ still only returns text (no tool handling yet)
	for _, c := range resp.Candidates {
		if c.Content == nil {
			continue
		}
		for _, p := range c.Content.Parts {
			if p.Text != "" {
				return p.Text, nil
			}
		}
	}

	return "", fmt.Errorf("empty response")
}
