package watch

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// ChangeEvent is emitted when a watched file changes.
type ChangeEvent struct {
	Path    string
	OldHash string
	NewHash string
	At      time.Time
}

// Watcher polls a set of file paths and emits ChangeEvents on change.
type Watcher struct {
	paths    []string
	hashes   map[string]string
	interval time.Duration
	Events   chan ChangeEvent
	Errors   chan error
	stop     chan struct{}
}

// New creates a Watcher that polls the given paths every interval.
func New(paths []string, interval time.Duration) (*Watcher, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("watch: interval must be positive")
	}
	w := &Watcher{
		paths:    paths,
		hashes:   make(map[string]string),
		interval: interval,
		Events:   make(chan ChangeEvent, 16),
		Errors:   make(chan error, 4),
		stop:     make(chan struct{}),
	}
	for _, p := range paths {
		h, err := hashFile(p)
		if err != nil {
			return nil, fmt.Errorf("watch: initial hash for %q: %w", p, err)
		}
		w.hashes[p] = h
	}
	return w, nil
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-w.stop:
				return
			case <-ticker.C:
				w.poll()
			}
		}
	}()
}

// Stop shuts down the background poller.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) poll() {
	for _, p := range w.paths {
		newHash, err := hashFile(p)
		if err != nil {
			w.Errors <- err
			continue
		}
		oldHash := w.hashes[p]
		if newHash != oldHash {
			w.Events <- ChangeEvent{
				Path:    p,
				OldHash: oldHash,
				NewHash: newHash,
				At:      time.Now(),
			}
			w.hashes[p] = newHash
		}
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
