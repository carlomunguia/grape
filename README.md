# 🍇 grape

A fast, concurrent file search tool written in Go - like `grep`, but concurrent and easier to use.

## Features

✨ **Concurrent Search** - Uses worker pools for fast parallel file searching  
🎯 **Case-Insensitive Search** - Optional `-i` flag for case-insensitive matching  
🎨 **Colorized Output** - Beautiful colored output (can be disabled)  
🚀 **Configurable Workers** - Adjust concurrency with `-w` flag  
🛡️ **Smart Filtering** - Automatically skips binary files, hidden directories, and common paths like `node_modules`  
⚡ **Fast & Efficient** - Handles large directory trees with configurable file size limits  
🔄 **Graceful Shutdown** - Handles Ctrl+C interrupts cleanly

## Installation

### From Source

```bash
git clone https://github.com/carlomunguia/grape.git
cd grape
go build -o grape ./grape
```

### Install Globally

```bash
go install github.com/carlomunguia/grape/grape@latest
```

## Usage

### Basic Search

```bash
# Search for "TODO" in current directory
./grape TODO

# Search for "error" in specific directory
./grape error ./src
```

### Advanced Options

```bash
# Case-insensitive search
./grape -i "hello" .

# Use 20 concurrent workers
./grape -w 20 "func main" .

# Disable colored output
./grape --color=false "import" .

# Verbose mode
./grape -v "package" .

# Combine flags
./grape -i -w 20 -v "error" ./logs
```

### Command-Line Options

```
Usage: grape [--workers WORKERS] [--ignore-case] [--verbose] [--color] SEARCHTERM [SEARCHDIR]

Positional arguments:
  SEARCHTERM             text to search for
  SEARCHDIR              directory to search (default: current directory)

Options:
  --workers WORKERS, -w WORKERS
                         number of concurrent workers [default: 10]
  --ignore-case, -i      case-insensitive search
  --verbose, -v          show verbose output
  --color, -c            colorize output [default: true]
  --help, -h             display this help and exit
```

## How It Works

grape uses a **concurrent worker pool pattern**:

1. **Discovery Phase** - Recursively traverses directories and adds files to a work queue
2. **Worker Pool** - Multiple workers pull files from the queue and search concurrently
3. **Smart Filtering** - Skips binary files, large files (>10MB), and common excluded paths
4. **Results Display** - Matches are displayed as they're found

### Architecture

```
Directory Discovery → Work Queue → Worker Pool → Results Channel → Display
```

## Performance

- **Concurrent workers** can be tuned based on your system (default: 10)
- **Binary file detection** prevents wasting time on non-text files
- **UTF-8 validation** ensures clean output
- **Line length limits** prevent memory issues with malformed files
- **Smart directory filtering** skips `.git`, `node_modules`, and hidden directories

## Requirements

- Go 1.25.3 or higher

## Examples

### Search for function definitions

```bash
./grape "func " ./myproject
```

### Find all TODO comments (case-insensitive)

```bash
./grape -i todo .
```

### Search with maximum concurrency

```bash
./grape -w 50 "error" /var/logs
```

### Debugging with verbose output

```bash
./grape -v "import" .
```

## Output Format

```
path/to/file.go[42]:matching line content
```

With colors enabled:

- 🟢 **Green** - File paths
- 🟡 **Yellow** - Line numbers
- ⚪ **White** - Line content

## Development

### Project Structure

```
grape/
├── grape/          # Main application
│   └── main.go     # Entry point and orchestration
├── worker/         # File searching logic
│   └── worker.go   # Search implementation
├── worklist/       # Thread-safe work queue
│   └── worklist.go # Queue implementation
├── go.mod          # Module dependencies
└── README.md       # This file
```

### Building

```bash
go build -o grape ./grape
```

### Running Tests

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See [LICENSE](LICENSE) file for details.

## Author

Carlo Munguia

---

**grape** - Picking the ripest results from your codebase, one search at a time! Cultivating concurrent searches with vine-tastic performance. 🍇
