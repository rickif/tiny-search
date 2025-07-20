# Tiny Research

This a deep research agent built with Golang, designed to perform comprehensive research tasks with intelligent information gathering and analysis capabilities. This project is mainly inspired by [DeerFlow](https://github.com/bytedance/deer-flow), thanks for their contributions.

## Project Structure

```
tiny-research/
├── .env                     # Environment variables
├── .gitignore              # Git ignore file
├── README.md               # Project documentation
├── go.mod                  # Go module definition
├── go.sum                  # Go dependencies checksum
├── main.go                 # Application entry point
├── internal/               # Private application code
│   ├── agent/             # Core agent implementation
│   │   ├── agent.go       # Main agent orchestrator
│   │   ├── coder.go       # Code generation agent
│   │   ├── coordinator.go # Workflow coordinator
│   │   ├── executor.go    # Task execution engine
│   │   ├── planner.go     # Research planning strategies
│   │   ├── reporter.go    # Report generation agent
│   │   ├── research_team.go # Research team coordination
│   │   ├── researcher.go  # Individual researcher agent
│   │   └── state.go       # Agent state management
│   ├── config/            # Configuration management
│   │   └── config.go      # Environment-based configuration
│   ├── llm/               # Language model integrations
│   │   └── llm.go         # LLM client wrapper
│   ├── prompts/           # Prompt templates
│   │   ├── coder.md       # Code generation prompts
│   │   ├── coordinator.md # Coordination prompts
│   │   ├── planner.md     # Planning prompts
│   │   ├── reporter.md    # Reporting prompts
│   │   └── researcher.md  # Research prompts
│   └── tool/              # Research tools
│       ├── crawl.go       # Web crawling (Jina AI)
│       ├── python.go      # Python code execution
│       └── search.go      # Web search (Tavily API)
└── util/                  # Utility functions
    └── json.go            # JSON processing utilities
```

## Core Components

### Multi-Agent System (`internal/agent/`)
A sophisticated workflow orchestration system with specialized agents:
- **Coordinator**: Orchestrates the overall research workflow and determines next steps
- **Planner**: Creates detailed research plans with structured steps and strategies
- **Research Team**: Manages team-based research coordination and task distribution
- **Researcher**: Executes individual research tasks using available tools
- **Coder**: Handles code generation and programming-related research tasks
- **Reporter**: Synthesizes findings into comprehensive reports
- **Agent State**: Manages conversation history, plans, and workflow state

### Research Tools (`internal/tool/`)
Integrated tools for information gathering and processing:
- **Tavily Search**: Web search capabilities using Tavily API for real-time information
- **Web Crawling**: Content extraction from URLs using Jina AI's reader service
- **Bash Execution**: Command-line tool execution for system operations
- **Python Execution**: Python code execution for data processing and analysis

### LLM Integration (`internal/llm/`)
LangChain Go integration for language model operations:
- **OpenAI Integration**: Primary LLM provider with configurable models and endpoints
- **Prompt Templates**: Specialized prompts for each agent type stored in markdown files
- **Tool Integration**: Seamless integration between LLM and research tools

### Configuration Management (`internal/config/`)
Environment-based configuration system:
- **Environment Variables**: Secure API key and configuration management
- **LLM Configuration**: Configurable model, base URL, and authentication
- **Tool Configuration**: API keys for external services (Tavily, etc.)

## Features

- 🔍 **Intelligent Research Planning**: Automatically breaks down complex research queries
- 🌐 **Multi-Source Integration**: Gathers information from diverse sources
- 🧠 **LLM-Powered Analysis**: Leverages advanced language models for understanding
- 📊 **Comprehensive Reporting**: Generates detailed research reports

## Getting Started

### Prerequisites
- Go 1.23 or higher
- OpenAI API key (or compatible LLM provider)
- Tavily API key for web search
- Internet connection for web-based research

### Installation

```bash
# Clone the repository
git clone https://github.com/rickif/tiny-research.git
cd tiny-research

go mod download

# Build the application
go build -o tiny-research .
```

### Configuration

1. Create a `.env` file in the project root:
```bash
cp .env.example .env  # If example exists, or create manually
```

2. Edit the `.env` file with your API keys and configuration:
```env
# LLM Configuration
LLM_MODEL=gpt-4o-mini
LLM_BASE_URL=https://api.openai.com/v1
LLM_TOKEN=your_openai_api_key_here

# Search Configuration
TAVILY_KEY=your_tavily_api_key_here
```

3. Run the research agent:
```bash
./tiny-research
```

The application will execute the hardcoded research query. To customize the query, modify the `main.go` file.

## Usage Examples

Currently, the research query is hardcoded in `main.go`. The default example query is:

```go
result, err := agent.Research(context.Background(), "What's the weather like in Chengdu today?")
```

### Customizing Research Queries

To research different topics, modify the query in `main.go`:

```go
// Example queries you can try:
result, err := agent.Research(context.Background(), "What are the latest developments in quantum computing?")
result, err := agent.Research(context.Background(), "Explain the impact of AI on healthcare")
result, err := agent.Research(context.Background(), "Write a Python script to analyze stock market data")
```

### Multi-Agent Workflow

The system automatically:
1. **Coordinates** the research workflow
2. **Plans** the research strategy with structured steps
3. **Research** for information using Tavily API and Jina AI search
4. **Code** code when programming is needed
5. **Reports** comprehensive findings

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [deerflow](https://github.com/bytedance/deer-flow)
- [langchaingo](https://github.com/tmc/langchaingo)