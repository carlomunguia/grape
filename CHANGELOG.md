# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-23

### Added

- Concurrent file search with configurable worker pool
- Case-insensitive search option (`-i` flag)
- Colorized output with option to disable (`-c` flag)
- Verbose mode for debugging (`-v` flag)
- Automatic binary file detection and skipping
- UTF-8 validation for clean output
- Smart directory filtering (skips `.git`, `node_modules`, `vendor`)
- File size limit (10MB) to prevent memory issues
- Line length limit (10,000 chars) for safety
- Graceful shutdown on interrupt (Ctrl+C)
- Context-aware cancellation support

### Features

- Thread-safe work queue for distributing tasks
- Configurable number of concurrent workers (default: 10)
- Colored output: green paths, yellow line numbers
- Multiple WaitGroups for proper synchronization
- Zero external dependencies (except CLI argument parsing)

### Performance

- Concurrent directory discovery
- Parallel file processing
- Efficient channel-based communication
- Smart filtering to avoid unnecessary work
