package llm

import (
	"context"
	"encoding/json"

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
