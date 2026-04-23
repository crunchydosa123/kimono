package tool

import (
	"context"
	"os"
)

type WriteFile struct{}

func (w *WriteFile) Name() string {
	return "write_file"
}

func (w *WriteFile) Description() string {
	return "Writes content to a file. Input: JSON {\"path\": \"\", \"content\": \"\"}"
}

func (w *WriteFile) Run(ctx context.Context, input string) (string, error) {
	// later: parse JSON properly
	path := "output.txt"
	content := input

	err := os.WriteFile(path, []byte(content), 0644)
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
