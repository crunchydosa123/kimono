package tool

import (
	"context"
	"encoding/json"
	"os/exec"
)

type SearchCode struct{}

func (t *SearchCode) Name() string { return "search_code" }

func (t *SearchCode) Description() string {
	return "Search for a string in the codebase and return matching file paths and lines"
}

func (t *SearchCode) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query": map[string]any{"type": "string"},
		},
		"required": []string{"query"},
	}
}

func (t *SearchCode) Run(ctx context.Context, input json.RawMessage) (string, error) {
	var args struct {
		Query string `json:"query"`
	}

	if err := json.Unmarshal(input, &args); err != nil {
		return "", err
	}

	cmd := exec.CommandContext(ctx, "grep", "-rn", args.Query, ".")
	out, err := cmd.CombinedOutput()

	// grep returns error if no matches → handle that
	if err != nil && len(out) == 0 {
		return "no matches found", nil
	}

	return string(out), nil
}
