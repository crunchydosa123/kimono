# Kimono

A Go-based coding agent powered by Google's Gemini AI that provides intelligent coding assistance through a conversational interface.

## Features

- **AI-Powered Coding Assistant**: Uses Google's Gemini 2.5 Flash model for intelligent code analysis and generation
- **Tool-Based Architecture**: Extensible tool system for file operations and code manipulation
- **Interactive Chat Interface**: Conversational interface for natural coding assistance
- **File Operations**: Built-in tools for reading, writing, editing, and searching code files
- **Command Suggestions**: Intelligent terminal command suggestions for development tasks

## Architecture

The project is structured as follows:

- `agent/`: Core agent logic that orchestrates LLM interactions and tool execution
- `chat_session/`: Chat session management for conversational interactions
- `internal/llm/`: LLM interface and Gemini implementation
- `internal/tool/`: Tool registry and implementations for file operations

## Prerequisites

- Go 1.25.6 or later
- Google Gemini API key

## Installation

1. Clone the repository:
```bash
git clone https://github.com/crunchydosa123/kimono.git
cd kimono
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your environment variables:
```bash
export GEMINI_API_KEY="your-gemini-api-key-here"
```

## Usage

Run the coding agent:

```bash
go run main.go
```

The agent will start an interactive session where you can ask for coding assistance. For example:

- "Create a new Go file with a hello world function"
- "Search for all functions in the codebase"
- "Edit the main.go file to add error handling"
- "Suggest commands to build and run this project"

## Available Tools

- **write_file**: Create or overwrite files
- **edit_file**: Modify existing files
- **search_code**: Search for code patterns across the codebase
- **suggest_command**: Suggest terminal commands for development tasks

## Configuration

The agent is configured with a system prompt that defines its capabilities and rules. You can modify the system prompt in `main.go` to customize the agent's behavior.

## Development

### Building

```bash
go build -o kimono main.go
```

### Testing

```bash
go test ./...
```

### Adding New Tools

1. Implement the `Tool` interface in `internal/tool/tool.go`
2. Register the tool in the registry in `main.go`

## Dependencies

- [Google Generative AI Go SDK](https://github.com/google/generative-ai-go)
- [Charm Bracelet libraries](https://github.com/charmbracelet) for terminal UI

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- Powered by [Google Gemini](https://ai.google.dev/)
- Built with [Go](https://golang.org/)</content>
<parameter name="filePath">/Users/prathamgadkari/Projects/kimono/README.md