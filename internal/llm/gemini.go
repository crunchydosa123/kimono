package llm

import (
	"context"
	"encoding/json"
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
		model:  "gemini-2.5-flash",
	}, nil
}

func (g *Gemini) ChatCompletion(
	ctx context.Context,
	messages []Message,
	tools []tool.Tool,
) (LLMResponse, error) {

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

	// attach tools
	config := &genai.GenerateContentConfig{}
	if len(tools) > 0 {
		config.Tools = tool.ToGeminiTools(tools)
	}

	resp, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		contents,
		config,
	)
	if err != nil {
		return LLMResponse{}, err
	}

	var llmRes LLMResponse

	for _, c := range resp.Candidates {
		var parts []Part

		for _, p := range c.Content.Parts {
			part := Part{}

			if p.Text != "" {
				part.Text = &p.Text
			}

			if p.FunctionCall != nil {
				argsBytes, _ := json.Marshal(p.FunctionCall.Args)

				part.ToolCall = &ToolCall{
					Name: p.FunctionCall.Name,
					Args: argsBytes,
				}
			}

			parts = append(parts, part)
		}

		llmRes.Candidates = append(llmRes.Candidates, Candidate{
			Role:         c.Content.Role,
			Parts:        parts,
			FinishReason: string(c.FinishReason),
		})
	}

	return llmRes, nil

}

func (g *Gemini) GeneratePlan(ctx context.Context, message Message) (string, error) {
	planPrompt := `
You are a planning agent.

Given a user request, generate a clear execution plan.

Rules:
- Output ONLY numbered steps
- Each step must be actionable
- Avoid redundancy
- Choose ONE approach (do not mix strategies)
- Prefer CLI commands when applicable
`

	messages := []Message{
		{Role: "system", Content: planPrompt},
		{Role: "user", Content: message.Content},
	}

	resp, err := g.ChatCompletion(ctx, messages, nil)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Parts) == 0 {
		return "", fmt.Errorf("no plan generated")
	}

	part := resp.Candidates[0].Parts[0]
	if part.Text == nil {
		return "", fmt.Errorf("plan response not text")
	}

	return *part.Text, nil
}
