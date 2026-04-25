package tool

import (
	"context"
	"encoding/json"
	"fmt"
)

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(t Tool) {
	r.tools[t.Name()] = t
}

func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

func (r *Registry) All() []Tool {
	out := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		out = append(out, t)
	}
	return out
}

func (r *Registry) Execute(ctx context.Context, toolName string, toolArgs json.RawMessage) (string, error) {
	t, ok := r.tools[toolName]
	if !ok {
		return "", fmt.Errorf("tool not found: %s", toolName)
	}

	return t.Run(ctx, toolArgs)
}
