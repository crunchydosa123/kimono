package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type EditFile struct{}

func (t *EditFile) Name() string { return "edit_file" }

func (t *EditFile) Description() string {
	return "Edit a file by replacing old content with new content"
}

func (t *EditFile) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{"type": "string"},
			"old":  map[string]any{"type": "string"},
			"new":  map[string]any{"type": "string"},
		},
		"required": []string{"path", "old", "new"},
	}
}

func (t *EditFile) Run(ctx context.Context, input json.RawMessage) (string, error) {
	var args struct {
		Path string `json:"path"`
		Old  string `json:"old"`
		New  string `json:"new"`
	}

	if err := json.Unmarshal(input, &args); err != nil {
		return "", err
	}

	data, err := os.ReadFile(args.Path)
	if err != nil {
		return "", err
	}

	if !strings.Contains(string(data), args.Old) {
		return "", fmt.Errorf("old content not found in file")
	}

	updated := strings.Replace(string(data), args.Old, args.New, 1)

	err = os.WriteFile(args.Path, []byte(updated), 0644)
	if err != nil {
		return "", err
	}

	return "file updated", nil
}
