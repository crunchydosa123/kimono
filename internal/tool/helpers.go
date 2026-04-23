package tool

import "google.golang.org/genai"

func FormatTools(tools []Tool) []map[string]any {
	var out []map[string]any

	for _, t := range tools {
		out = append(out, map[string]any{
			"name":         t.Name(),
			"description":  t.Description(),
			"input_schema": t.InputSchema(),
		})
	}

	return out
}

func JSONSchemaToGemini(s map[string]any) *genai.Schema {
	schema := &genai.Schema{}

	// type
	if t, ok := s["type"].(string); ok {
		switch t {
		case "object":
			schema.Type = genai.TypeObject
		case "string":
			schema.Type = genai.TypeString
		case "number":
			schema.Type = genai.TypeNumber
		case "integer":
			schema.Type = genai.TypeInteger
		case "boolean":
			schema.Type = genai.TypeBoolean
		case "array":
			schema.Type = genai.TypeArray
		}
	}

	// properties
	if props, ok := s["properties"].(map[string]any); ok {
		schema.Properties = make(map[string]*genai.Schema)
		for k, v := range props {
			if child, ok := v.(map[string]any); ok {
				schema.Properties[k] = JSONSchemaToGemini(child)
			}
		}
	}

	// required
	if req, ok := s["required"].([]string); ok {
		schema.Required = req
	} else if reqAny, ok := s["required"].([]any); ok {
		for _, r := range reqAny {
			if str, ok := r.(string); ok {
				schema.Required = append(schema.Required, str)
			}
		}
	}

	// description
	if desc, ok := s["description"].(string); ok {
		schema.Description = desc
	}

	return schema
}

func ToGeminiTools(tools []Tool) []*genai.Tool {
	var out []*genai.Tool

	for _, t := range tools {
		out = append(out, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        t.Name(),
					Description: t.Description(),
					Parameters:  JSONSchemaToGemini(t.InputSchema()),
				},
			},
		})
	}

	return out
}
