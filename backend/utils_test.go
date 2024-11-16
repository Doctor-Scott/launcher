package backend

import (
	"os"
	"reflect"
	"testing"
)

func TestResolvePath(t *testing.T) {
	env := os.Getenv("HOME")

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/test",
			expected: "/test/",
		},
		{
			name:     "path with tilde",
			path:     "~",
			expected: env + "/",
		},
		{
			name:     "path with env var",
			path:     "$HOME",
			expected: env + "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolvePath(tt.path)
			if result != tt.expected {
				t.Errorf("ResolvePath(%v) = %v, want %v",
					tt.path, result, tt.expected)
			}
		})
	}
}

func TestPrintStructure(t *testing.T) {
	cwd, _ := os.Getwd()
	testDir := cwd + "/test/script"
	testScriptFileName := "test-script.sh"
	testFilePath := testDir + "/test-script.sh"
	tests := []struct {
		name     string
		path     string
		expected []Script
	}{
		{
			name:     "test in testDir",
			path:     testDir,
			expected: []Script{{Name: testScriptFileName, Command: testFilePath, Args: []string{}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStructure(tt.path)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetStructure(%v) = %v, want %v",
					tt.path, result, tt.expected)
			}
		})
	}

}
