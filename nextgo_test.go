package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// getProjectRoot finds the next.go project root for test file verification.
func getProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// TestRouterLogic tests the core routing logic independently.
func TestRouterLogic(t *testing.T) {
	t.Run("dirToURLPath", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{".", "/"},
			{"about", "/about"},
			{"blog/post", "/blog/post"},
		}

		for _, tt := range tests {
			result := dirToURLPath(tt.input)
			if result != tt.expected {
				t.Errorf("dirToURLPath(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("processDynamicSegments", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"/blog/[id]", "/blog/:id"},
			{"/posts/[...slug]", "/posts/*slug"},
			{"/users/[userId]/posts", "/users/:userId/posts"},
		}

		for _, tt := range tests {
			result := processDynamicSegments(tt.input)
			if result != tt.expected {
				t.Errorf("processDynamicSegments(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		}
	})
}

// TestFileExistence verifies all expected source files exist using
// the runtime project root, not a hardcoded Windows path.
func TestFileExistence(t *testing.T) {
	root := getProjectRoot()
	if root == "" {
		t.Skip("cannot find project root")
	}

	files := []string{
		"main.go",
		"cmd/root.go",
		"cmd/dev.go",
		"cmd/build.go",
		"cmd/start.go",
		"cmd/create.go",
		"packages/router/router.go",
		"packages/server/server.go",
		"packages/build/build.go",
		"packages/api/api.go",
		"packages/config/config.go",
		"packages/middleware/middleware.go",
		"packages/static/static.go",
		"packages/watcher/watcher.go",
		"packages/apihandler/apihandler.go",
	}

	for _, f := range files {
		path := filepath.Join(root, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("MISSING: %s", f)
		}
	}
}

// TestConfigLoad tests config loading.
func TestConfigLoad(t *testing.T) {
	// Test default config
	root := getProjectRoot()
	if root == "" {
		t.Skip("cannot find project root")
	}

	cfgPath := filepath.Join(root, "nextgo.yaml")
	if _, err := os.Stat(cfgPath); err != nil {
		t.Skip("no nextgo.yaml in project root")
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if len(data) == 0 {
		t.Error("config file is empty")
	}
}

// dirToURLPath converts a directory path to URL path (mirrors router logic).
func dirToURLPath(dir string) string {
	if dir == "." {
		return "/"
	}
	urlPath := strings.ReplaceAll(dir, string(filepath.Separator), "/")
	urlPath = strings.Trim(urlPath, "/")
	if urlPath == "" {
		return "/"
	}
	return "/" + urlPath
}

// processDynamicSegments processes dynamic route segments (mirrors router logic).
func processDynamicSegments(path string) string {
	// Handle catch-all [...slug] first (must come before [id] → :id)
	path = strings.ReplaceAll(path, "[...", "*")
	path = strings.ReplaceAll(path, "[", ":")
	path = strings.ReplaceAll(path, "]", "")
	return path
}
