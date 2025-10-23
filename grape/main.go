package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/carlomunguia/grape/worker"
	"github.com/carlomunguia/grape/worklist"

	"github.com/alexflint/go-arg"
)

// discoverDirs recursively traverses directories and adds files to the worklist.
// It respects context cancellation and skips hidden directories and common excluded paths.
func discoverDirs(ctx context.Context, wl *worklist.Worklist, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", path, err)
	}
	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden directories and common excluded paths
		if len(name) > 0 && (name[0] == '.' || name == "node_modules" || name == "vendor") {
			continue
		}
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())
			if err := discoverDirs(ctx, wl, nextPath); err != nil {
				// Don't log context cancellation errors (user interrupted)
				if err != context.Canceled {
					fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
				}
			}
		} else {
			wl.Add(worklist.NewEntry(filepath.Join(path, entry.Name())))
		}
	}
	return nil
}

var args struct {
	SearchTerm      string `arg:"positional, required" help:"text to search for"`
	SearchDir       string `arg:"positional" help:"directory to search (default: current directory)"`
	Workers         int    `arg:"-w,--workers" default:"10" help:"number of concurrent workers"`
	CaseInsensitive bool   `arg:"-i,--ignore-case" help:"case-insensitive search"`
	Verbose         bool   `arg:"-v,--verbose" help:"show verbose output"`
	ColorOutput     bool   `arg:"-c,--color" help:"colorize output" default:"true"`
}

// main is the entry point for the grape text search tool.
// It performs concurrent file searching across a directory tree.
func main() {
	arg.MustParse(&args)

	if args.SearchDir == "" {
		args.SearchDir = "."
	}

	// Validate workers count
	if args.Workers < 1 {
		fmt.Fprintf(os.Stderr, "Error: workers must be at least 1\n")
		os.Exit(1)
	}

	// Validate search directory exists
	if _, err := os.Stat(args.SearchDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directory '%s' does not exist\n", args.SearchDir)
		os.Exit(1)
	}

	if args.Verbose {
		fmt.Fprintf(os.Stderr, "Searching for '%s' in '%s' with %d workers\n",
			args.SearchTerm, args.SearchDir, args.Workers)
	}

	var workersWg sync.WaitGroup

	wl := worklist.New(100)

	results := make(chan worker.Result, 100)

	numWorkers := args.Workers

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		fmt.Fprintln(os.Stderr, "\nInterrupt received, shutting down...")
		cancel()
	}()

	// Start directory discovery
	workersWg.Add(1)
	go func() {
		defer workersWg.Done()
		if err := discoverDirs(ctx, wl, args.SearchDir); err != nil && err != context.Canceled {
			fmt.Fprintf(os.Stderr, "Error during discovery: %v\n", err)
		}
		wl.Close()
		if args.Verbose {
			fmt.Fprintf(os.Stderr, "Discovery complete\n")
		}
	}()

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		workerID := i
		go func() {
			defer workersWg.Done()
			if args.Verbose {
				fmt.Fprintf(os.Stderr, "Worker %d started\n", workerID)
			}
			for {
				workEntry, ok := wl.NextWithContext(ctx)
				if !ok {
					if args.Verbose {
						fmt.Fprintf(os.Stderr, "Worker %d finished\n", workerID)
					}
					return
				}
				workerResult := worker.FindInFile(workEntry.Path, args.SearchTerm, args.CaseInsensitive)
				if workerResult != nil {
					for _, r := range workerResult.Inner {
						select {
						case results <- r:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}()
	}

	var displayWg sync.WaitGroup

	// Start display goroutine
	displayWg.Add(1)
	go func() {
		defer displayWg.Done()
		for r := range results {
			if args.ColorOutput {
				// Green path, yellow line number, white content
				fmt.Printf("\033[32m%v\033[0m[\033[33m%v\033[0m]:%v\n", r.Path, r.LineNum, r.Line)
			} else {
				fmt.Printf("%v[%v]:%v\n", r.Path, r.LineNum, r.Line)
			}
		}
	}()

	// Close results channel when all workers are done
	go func() {
		workersWg.Wait()
		close(results)
		if args.Verbose {
			fmt.Fprintf(os.Stderr, "All workers finished\n")
		}
	}()

	displayWg.Wait()
}
