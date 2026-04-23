package tool

import "context"

type Tool interface {
	Name() string
	Description() string
	InputSchema() map[string]any
	Run(ctx context.Context, input string) (string, error)
}
