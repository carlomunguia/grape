// Package worker provides file search functionality for concurrent text searching.
package worker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type Result struct {
	Line    string
	LineNum int
	Path    string
}

type Results struct {
	Inner []Result
}

func NewResult(line string, lineNum int, path string) Result {
	return Result{Line: line, LineNum: lineNum, Path: path}
}

func NewResults() *Results {
	return &Results{
		Inner: make([]Result, 0, 10), // Pre-allocate capacity of 10
	}
}

const maxFileSize = 10 * 1024 * 1024 // 10 MB
const maxLineLength = 10000          // Maximum line length to process

// FindInFile searches for a string in a file and returns all matching lines.
// If caseInsensitive is true, the search is case-insensitive.
func FindInFile(path string, find string, caseInsensitive bool) *Results {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil
	}

	// Skip files larger than 10MB
	if fileInfo.Size() > maxFileSize {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Skip binary files
	if isBinary, err := isBinaryFile(file); err != nil || isBinary {
		return nil
	}

	results := NewResults()

	// Prepare search term for case-insensitive comparison
	searchTerm := find
	if caseInsensitive {
		searchTerm = strings.ToLower(find)
	}

	scanner := bufio.NewScanner(file)
	// Increase buffer size for long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 1MB max token size

	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()

		// Skip extremely long lines to prevent memory issues
		if len(line) > maxLineLength {
			lineNum++
			continue
		}

		// Skip lines with invalid UTF-8
		if !utf8.ValidString(line) {
			lineNum++
			continue
		}

		lineToSearch := line

		if caseInsensitive {
			lineToSearch = strings.ToLower(line)
		}

		if strings.Contains(lineToSearch, searchTerm) {
			r := NewResult(line, lineNum, path)
			results.Inner = append(results.Inner, r)
		}
		lineNum++
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
		return nil
	}

	if len(results.Inner) == 0 {
		return nil
	}

	return results
}

// isBinaryFile checks if a file appears to be binary by reading the first chunk
func isBinaryFile(file *os.File) (bool, error) {
	// Read first 512 bytes
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return false, err
	}

	// Check for null bytes and control characters
	controlCharCount := 0
	for i := 0; i < n; i++ {
		// Null byte is strong indicator of binary
		if buf[i] == 0 {
			return true, nil
		}
		// Count suspicious control characters (excluding common ones like tab, newline, carriage return)
		if buf[i] < 0x20 && buf[i] != '\t' && buf[i] != '\n' && buf[i] != '\r' {
			controlCharCount++
		}
	}

	// If more than 30% are control characters, likely binary
	if n > 0 && float64(controlCharCount)/float64(n) > 0.3 {
		return true, nil
	}

	return false, nil
}
