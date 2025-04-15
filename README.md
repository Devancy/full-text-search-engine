# Full-Text Search Engine

A learning project to explore Go programming concepts by implementing a full-text search engine. This project demonstrates various Go features including concurrent programming, interfaces, and efficient data structures.

## Learning Objectives

This project was created to learn and demonstrate:
- Go's concurrency patterns (goroutines, channels, sync package)
- Interface design and implementation
- Efficient data structures in Go
- Memory management and pointer usage
- Package organization and project structure
- Performance optimization techniques
- Real-world algorithm implementation (TF-IDF)

## Features

- Fast full-text search using inverted index
- TF-IDF scoring for better search relevance
- Support for both simple and concurrent indexing
- Real-time search with interactive CLI
- Processes Wikipedia abstract dumps
- Memory-efficient document handling
- Detailed index statistics
- Comprehensive benchmarking suite
- Progress reporting and improved error handling
- **Robust interactive input with line editing, history, and arrow key support using [github.com/chzyer/readline](https://github.com/chzyer/readline)**

## Code Organization

The project is structured to demonstrate different Go concepts:
```
.
├── main.go                 # Entry point, CLI handling
├── utils/
│   ├── document.go         # Document handling and XML parsing
│   ├── index.go            # Simple indexing implementation
│   ├── index_concurrent.go # Concurrent indexing (advanced)
│   ├── index_interface.go  # Interface definitions
│   ├── concurrent_types.go # Thread-safe types
│   ├── tokenizer.go        # Text analysis
│   └── filter.go           # Text filtering utilities
```

## Prerequisites

- Go 1.18 or later
- Wikipedia abstract dump file (XML format, gzipped)

## Installation

```bash
# Clone the repository
git clone https://github.com/devancy/full-text-search-engine.git
cd full-text-search-engine

# Install dependencies
go mod download
```

## Data Preparation

1. Download the Wikipedia abstract dump:
   - Visit [Wikimedia Dumps](https://dumps.wikimedia.your.org/enwiki/latest/)
   - Download `enwiki-latest-abstract1.xml.gz`
   - Place it in your working directory or note its path

## Usage

### Basic Usage

```bash
# Run with default settings (simple indexing)
go run main.go

# Specify custom dump file path
go run main.go -p "path/to/enwiki-latest-abstract1.xml.gz"
```

### Advanced Options

```bash
# Use concurrent indexing for better performance
go run main.go -c

# Combine options
go run main.go -p "path/to/dump.xml.gz" -c
```

### Command Line Flags

- `-p`: Specify the path to the Wikipedia dump file (default: "enwiki-latest-abstract1.xml.gz")
- `-c`: Enable concurrent indexing for faster processing (default: false)
- `-n`: Maximum number of search results to display (default: 5)

### Interactive Search

After indexing completes:
1. Type your search query
2. Press Enter to search
3. Results show:
   - Document title
   - Relevance score
   - URL (if available)
   - Abstract text
   - Clear separation between results
4. Press Ctrl+C to exit
5. **Enjoy advanced line editing, history, and arrow key navigation in the search prompt thanks to the readline library!**

## Implementation Details

The project implements two indexing strategies to demonstrate different Go concepts:

### Simple Index (Learning Basics)
- Sequential document processing
- Basic Go data structures
- Easy to understand for beginners
- Good for learning memory management
- Demonstrates basic package organization
- Includes performance benchmarks

### Concurrent Index (Advanced Concepts)
- Demonstrates Go's concurrency features
- Uses goroutines and channels
- Shows sync package usage
- Implements thread-safe data structures
- Advanced error handling
- Parallel processing for better performance
- Includes comparative benchmarks

Both implementations provide:
- TF-IDF scoring for ranking results
- Index statistics (document count, term count, etc.)
- Memory usage information
- Performance metrics

Note:
- TF (Term Frequency): Measures word importance in a document
- IDF (Inverse Document Frequency): Measures word importance across all documents
- Final score = TF * IDF

## Benchmarking

The project includes comprehensive benchmarks to compare performance:

```bash
# Run all benchmarks
go test -bench=. ./utils

# Run specific benchmark groups
go test -bench=BenchmarkIndexAdd ./utils
go test -bench=BenchmarkIndexAddLarge ./utils
```

Benchmark scenarios include:
- Document indexing (1,000 documents)
- Large-scale indexing (1000,000 documents)
- Comparative analysis between simple and concurrent implementations

## Performance Considerations

- Simple indexing is recommended for:
  - Learning Go basics
  - Understanding the core algorithm
  - Debugging and testing
  - Small datasets (< 100,000 documents)
  - Memory-constrained environments

- Concurrent indexing is optimal for:
  - Production workloads
  - Large datasets (> 100,000 documents)
  - Multi-core systems
  - High-throughput requirements
  - Real-time search applications

Performance metrics show:
- Concurrent indexing typically 2-4x faster for large datasets
- Search performance scales well with document count
- Memory usage grows linearly with document count

## Contributing

This is a learning project, and contributions that help demonstrate Go concepts are welcome! Feel free to:
- Add comments explaining complex parts
- Improve documentation
- Add more Go features
- Optimize performance
- Fix bugs

## License

This project is licensed under the MIT License.

## Acknowledgments

This project was created as a learning exercise to understand Go programming concepts. It's meant to be educational and may not be production-ready.
