# Deep Research Agent

A sophisticated research agent built with Go and LangChain Go, designed to perform comprehensive research tasks with intelligent information gathering and analysis capabilities.

## Overview

This project implements a deep research agent that can autonomously conduct research on various topics, gather information from multiple sources, analyze data, and provide comprehensive reports. The agent leverages the power of large language models through LangChain Go to understand context, plan research strategies, and synthesize findings.

## Technology Stack

- **Go**: Core programming language for performance and concurrency
- **LangChain Go**: Framework for building LLM-powered applications
- **Additional technologies**: To be determined based on specific requirements

## Project Structure

```
tiny-research/
├── README.md                 # Project documentation
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies checksum
├── main.go                  # Application entry point
├── cmd/                     # Command-line interfaces
│   └── research/            # Research agent CLI
├── internal/                # Private application code
│   ├── agent/              # Core agent implementation
│   │   ├── research.go     # Main research agent logic
│   │   ├── planner.go      # Research planning strategies
│   │   └── executor.go     # Task execution engine
│   ├── sources/            # Information source integrations
│   │   ├── web/            # Web scraping and search
│   │   ├── academic/       # Academic database access
│   │   └── documents/      # Document processing
│   ├── llm/                # Language model integrations
│   │   ├── client.go       # LLM client wrapper
│   │   ├── prompts/        # Prompt templates
│   │   └── chains/         # LangChain Go chains
│   ├── storage/            # Data persistence layer
│   │   ├── memory/         # In-memory storage
│   │   ├── database/       # Database operations
│   │   └── cache/          # Caching mechanisms
│   ├── analysis/           # Data analysis and processing
│   │   ├── summarizer.go   # Content summarization
│   │   ├── classifier.go   # Information classification
│   │   └── synthesizer.go  # Research synthesis
│   └── utils/              # Utility functions
│       ├── http.go         # HTTP utilities
│       ├── text.go         # Text processing
│       └── config.go       # Configuration management
├── pkg/                    # Public library code
│   ├── types/              # Shared data types
│   └── interfaces/         # Public interfaces
├── configs/                # Configuration files
│   ├── config.yaml         # Main configuration
│   └── prompts.yaml        # Prompt configurations
├── docs/                   # Documentation
│   ├── architecture.md     # System architecture
│   ├── api.md             # API documentation
│   └── examples/           # Usage examples
├── scripts/                # Build and deployment scripts
│   ├── build.sh           # Build script
│   └── deploy.sh          # Deployment script
└── tests/                  # Test files
    ├── integration/        # Integration tests
    └── unit/              # Unit tests
```

## Core Components

### Research Agent (`internal/agent/`)
The heart of the system that orchestrates the research process:
- **Research Planner**: Breaks down research queries into actionable tasks
- **Task Executor**: Executes individual research tasks
- **Result Synthesizer**: Combines findings into coherent reports

### Information Sources (`internal/sources/`)
Modular integrations for various information sources:
- **Web Sources**: Search engines, websites, APIs
- **Academic Sources**: Research papers, journals, databases
- **Document Processing**: PDF, Word, text file analysis

### LLM Integration (`internal/llm/`)
LangChain Go integration for language model operations:
- **Chain Management**: Custom chains for different research tasks
- **Prompt Engineering**: Optimized prompts for research scenarios
- **Model Abstraction**: Support for multiple LLM providers

### Storage Layer (`internal/storage/`)
Flexible storage solutions for research data:
- **Memory Storage**: Fast in-memory caching
- **Persistent Storage**: Database for long-term data
- **Cache Management**: Intelligent caching strategies

### Analysis Engine (`internal/analysis/`)
Advanced analysis capabilities:
- **Content Summarization**: Extract key insights from sources
- **Information Classification**: Categorize and tag findings
- **Research Synthesis**: Generate comprehensive reports

## Features

- 🔍 **Intelligent Research Planning**: Automatically breaks down complex research queries
- 🌐 **Multi-Source Integration**: Gathers information from diverse sources
- 🧠 **LLM-Powered Analysis**: Leverages advanced language models for understanding
- 📊 **Comprehensive Reporting**: Generates detailed research reports
- ⚡ **Concurrent Processing**: Efficient parallel information gathering
- 🔄 **Iterative Refinement**: Continuously improves research quality
- 💾 **Persistent Memory**: Maintains context across research sessions

## Getting Started

### Prerequisites
- Go 1.21 or higher
- API keys for LLM providers (OpenAI, Anthropic, etc.)
- Internet connection for web-based research

### Installation

```bash
# Clone the repository
git clone https://github.com/rickif/tiny-research.git
cd tiny-research

# Install dependencies
go mod download

# Build the application
go build -o research ./cmd/research
```

### Configuration

1. Copy the example configuration:
```bash
cp configs/config.yaml.example configs/config.yaml
```

2. Edit the configuration file with your API keys and preferences

3. Run the research agent:
```bash
./research "Your research query here"
```

## Usage Examples

```bash
# Basic research query
./research "What are the latest developments in quantum computing?"

# Research with specific sources
./research --sources=academic,web "Climate change impact on agriculture"

# Generate detailed report
./research --format=detailed "Artificial intelligence in healthcare"
```

## Development

### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests
go test ./tests/integration/...
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## Roadmap

- [ ] Basic research agent implementation
- [ ] Web source integration
- [ ] LangChain Go integration
- [ ] Academic database connectors
- [ ] Advanced analysis features
- [ ] Web interface
- [ ] API endpoints
- [ ] Plugin system for custom sources
- [ ] Multi-language support
- [ ] Cloud deployment options

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [LangChain Go](https://github.com/tmc/langchaingo) for the excellent LLM framework
- The Go community for amazing tools and libraries
- Contributors and researchers who inspire this project