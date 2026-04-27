package tool

import (
	"context"
	"encoding/json"
)

type SuggestCommand struct{}

func (t *SuggestCommand) Name() string {
	return "suggest_command"
}

func (t *SuggestCommand) Description() string {
	return "Suggest a terminal command to run. The user will approve before execution"
}

func (t *SuggestCommand) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The terminal command to run",
			},
		},
		"required": []string{"command"},
	}
}

func (t *SuggestCommand) Run(ctx context.Context, input json.RawMessage) (string, error) {
	var args struct {
		Command string `json:"command"`
	}

	if err := json.Unmarshal(input, &args); err != nil {
		return "", err
	}

	return args.Command, nil
}
