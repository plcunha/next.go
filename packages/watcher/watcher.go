package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches for file changes and triggers hot reload
type Watcher struct {
	AppDir     string
	OnChange   func()
	watcher    *fsnotify.Watcher
	debounce   map[string]time.Time
}

// New creates a new file watcher
func New(appDir string) (*Watcher, error) {
	w := &Watcher{
		AppDir:   appDir,
		debounce: make(map[string]time.Time),
	}

	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Watch starts watching for file changes
func Watch(appDir string) {
	w, err := New(appDir)
	if err != nil {
		fmt.Printf("Watcher error: %v\n", err)
		return
	}

	go w.watch()
}

// watch runs the file watcher
func (w *Watcher) watch() {
	// Add app directory to watcher
	err := w.addDirRecursive(w.AppDir)
	if err != nil {
		fmt.Printf("Error watching directory: %v\n", err)
		return
	}

	// Watch public directory too
	publicDir := filepath.Join(filepath.Dir(w.AppDir), "public")
	if _, err := os.Stat(publicDir); err == nil {
		w.addDirRecursive(publicDir)
	}

	fmt.Println("👀 Watching for file changes...")

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}

// addDirRecursive adds a directory and all subdirectories to the watcher
func (w *Watcher) addDirRecursive(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip node_modules and .next
		if info.Name() == "node_modules" || info.Name() == ".next" {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return w.watcher.Add(path)
		}

		return nil
	})
}

// handleEvent handles a file system event
func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Debounce events (same file can trigger multiple events)
	now := time.Now()
	lastTime, exists := w.debounce[event.Name]
	if exists && now.Sub(lastTime) < 500*time.Millisecond {
		return
	}
	w.debounce[event.Name] = now

	// Ignore certain files
	if shouldIgnore(event.Name) {
		return
	}

	// Handle different event types
	switch {
	case event.Op&fsnotify.Write == fsnotify.Write:
		w.onFileChanged(event.Name, "modified")
	case event.Op&fsnotify.Create == fsnotify.Create:
		w.onFileChanged(event.Name, "created")
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		w.onFileChanged(event.Name, "deleted")
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		w.onFileChanged(event.Name, "renamed")
	}
}

// onFileChanged handles file change events
func (w *Watcher) onFileChanged(path string, action string) {
	fmt.Printf("  File %s: %s\n", action, path)

	// Trigger reload callback
	if w.OnChange != nil {
		w.OnChange()
	}

	// Send HMR signal (in a real implementation, this would use WebSocket)
	fmt.Println("  🔄 Hot reload triggered")
}

// shouldIgnore returns true if the file should be ignored
func shouldIgnore(path string) bool {
	ignored := []string{
		".next",
		"node_modules",
		".git",
		"go.sum",
		".tmp",
		".DS_Store",
	}

	for _, ignore := range ignored {
		if contains(path, ignore) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Close closes the watcher
func (w *Watcher) Close() {
	if w.watcher != nil {
		w.watcher.Close()
	}
}
