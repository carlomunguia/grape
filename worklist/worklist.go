// Package worklist provides a thread-safe work queue for distributing
// file paths to concurrent workers.
package worklist

import "context"

// Entry represents a single work item containing a file path to process.
type Entry struct {
	Path string
}

// Worklist is a thread-safe queue for distributing work entries to multiple workers.
// It is safe for concurrent use by multiple goroutines.
type Worklist struct {
	jobs chan Entry
}

// New creates a new Worklist with the specified buffer size.
func New(bufferSize int) *Worklist {
	return &Worklist{
		jobs: make(chan Entry, bufferSize),
	}
}

// NewEntry creates a new work entry with the given file path.
func NewEntry(path string) Entry {
	return Entry{Path: path}
}

// Add adds a work entry to the worklist.
// This will block if the buffer is full.
func (w *Worklist) Add(work Entry) {
	w.jobs <- work
}

// Next retrieves the next work entry from the worklist.
// Returns false if the channel is closed and empty.
func (w *Worklist) Next() (Entry, bool) {
	j, ok := <-w.jobs
	return j, ok
}

// NextWithContext retrieves the next work entry or returns early if context is cancelled.
func (w *Worklist) NextWithContext(ctx context.Context) (Entry, bool) {
	select {
	case <-ctx.Done():
		return Entry{}, false
	case j, ok := <-w.jobs:
		return j, ok
	}
}

// Close closes the worklist channel, signaling that no more work will be added.
// Workers should drain remaining items before exiting.
func (w *Worklist) Close() {
	close(w.jobs)
}

// Len returns the approximate number of pending jobs in the worklist.
// This value may change immediately after calling.
func (w *Worklist) Len() int {
	return len(w.jobs)
}

// Cap returns the buffer capacity of the worklist.
func (w *Worklist) Cap() int {
	return cap(w.jobs)
}
