package tool

import (
	"context"
	"encoding/json"
	"os"
)

type WriteFile struct{}

func (w *WriteFile) Name() string {
	return "write_file"
}

func (w *WriteFile) Description() string {
	return "Writes content to a file. Input: JSON {\"path\": \"\", \"content\": \"\"}"
}

func (t *WriteFile) Run(ctx context.Context, input string) (string, error) {
	var args struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return "", err
	}

	err := os.WriteFile(args.Path, []byte(args.Content), 0644)
	if err != nil {
		return "", err
	}

	return "file written", nil
}

func (w *WriteFile) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "File path",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "File content",
			},
		},
		"required": []string{"path", "content"},
	}
}
