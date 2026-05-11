package router

import (
	"testing"
)

func TestNewRouter(t *testing.T) {
	// Create a temp directory for testing
	router := New("/tmp/test-app")
	if router == nil {
		t.Error("Router should not be nil")
	}
}

func TestDirToURLPath(t *testing.T) {
	router := &Router{}
	
	tests := []struct {
		input    string
		expected string
	}{
		{".", "/"},
		{"about", "/about"},
		{"blog/post", "/blog/post"},
	}
	
	for _, tt := range tests {
		result := router.dirToURLPath(tt.input)
		if result != tt.expected {
			t.Errorf("dirToURLPath(%s) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestProcessDynamicSegments(t *testing.T) {
	router := &Router{}
	
	tests := []struct {
		input    string
		expected string
	}{
		{"/blog/[id]", "/blog/:id"},
		{"/posts/[...slug]", "/posts/*slug"},
		{"/users/[userId]/posts", "/users/:userId/posts"},
	}
	
	for _, tt := range tests {
		result := router.processDynamicSegments(tt.input)
		if result != tt.expected {
			t.Errorf("processDynamicSegments(%s) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}
