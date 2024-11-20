package backend

import (
	"os"
	"reflect"
	"testing"
)

func TestResolveArgsString(t *testing.T) {
	// Set up test environment variable
	os.Setenv("TEST_VAR", "test-value")
	env := os.Getenv("HOME")

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple args",
			input:    "arg1 arg2 arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "double quoted args",
			input:    `arg1 "quoted arg" arg3`,
			expected: []string{"arg1", "quoted arg", "arg3"},
		},
		{
			name:     "single quoted args",
			input:    "arg1 'quoted arg' arg3",
			expected: []string{"arg1", "quoted arg", "arg3"},
		},
		{
			name:     "environment variables",
			input:    "$TEST_VAR $HOME",
			expected: []string{"test-value", env},
		},
		{
			name:     "mixed quotes and env vars",
			input:    `"$TEST_VAR" '$HOME' regular`,
			expected: []string{"test-value", "$HOME", "regular"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveArgsString(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("resolveArgsString(%q) = %v, want %v",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestRunScriptWithArgsString(t *testing.T) {
	// Set up test environment variable
	os.Setenv("TEST_VAR", "test-value")

	tests := []struct {
		name            string
		script          Script
		argsString      string
		stdin           []byte
		expectedStdout  []byte
		expectedSuccess bool
	}{
		{
			name: "simple script",
			script: Script{
				Name:    "test-script",
				Command: "echo",
			},
			argsString:      "hello world",
			stdin:           []byte{},
			expectedStdout:  []byte("hello world\n"),
			expectedSuccess: true,
		},
		{
			name: "script with environment variable",
			script: Script{
				Name:    "test-script",
				Command: "echo",
			},
			argsString:      "$TEST_VAR",
			stdin:           []byte{},
			expectedStdout:  []byte("test-value\n"),
			expectedSuccess: true,
		},
		{
			name: "script with quoted environment variable",
			script: Script{
				Name:    "test-script",
				Command: "echo",
			},
			argsString:      `"$TEST_VAR"`,
			stdin:           []byte{},
			expectedStdout:  []byte("test-value\n"),
			expectedSuccess: true,
		},
		{
			name: "script with unquoted environment variable",
			script: Script{
				Name:    "test-script",
				Command: "echo",
			},
			argsString:      "$TEST_VAR",
			stdin:           []byte{},
			expectedStdout:  []byte("test-value\n"),
			expectedSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := AddArgsToScript(tt.script, tt.argsString)
			scriptResult := RunScript(script, tt.stdin)
			result := scriptResult.Stdout
			success := scriptResult.Success
			if !reflect.DeepEqual(result, tt.expectedStdout) {
				t.Errorf("RunScript(%v, %v) = %v, want %v",
					tt.script, tt.stdin, result, tt.expectedStdout)
			}
			if !reflect.DeepEqual(success, tt.expectedSuccess) {
				t.Errorf("RunScript(%v, %v) = %v, want %v",
					tt.script, tt.stdin, success, tt.expectedSuccess)
			}
		})
	}
}

func TestRunKnownScript(t *testing.T) {
	// Set up test environment variable
	os.Setenv("TEST_VAR", "test-value")
	tests := []struct {
		name            string
		command         string
		stdin           []byte
		expectedStdout  []byte
		expectedSuccess bool
	}{
		{
			name:            "script with args",
			command:         "echo hello world",
			stdin:           []byte{},
			expectedStdout:  []byte("hello world\n"),
			expectedSuccess: true,
		},
		{
			name:            "script with env var",
			command:         "echo $TEST_VAR",
			stdin:           []byte{},
			expectedStdout:  []byte("test-value\n"),
			expectedSuccess: true,
		},
		{
			name:            "script with quoted env var",
			command:         "echo \"$TEST_VAR\"",
			stdin:           []byte{},
			expectedStdout:  []byte("test-value\n"),
			expectedSuccess: true,
		},
		{
			name:            "script with unquoted env var",
			command:         "echo $TEST_VAR",
			stdin:           []byte{},
			expectedStdout:  []byte("test-value\n"),
			expectedSuccess: true,
		},
		{
			name:            "script with quoted args",
			command:         "echo 'hello world'",
			stdin:           []byte{},
			expectedStdout:  []byte("hello world\n"),
			expectedSuccess: true,
		},
		{
			name:            "script with unquoted args",
			command:         "echo hello world",
			stdin:           []byte{},
			expectedStdout:  []byte("hello world\n"),
			expectedSuccess: true,
		},
		{
			name:            "script with quoted and unquoted args",
			command:         "echo 'hello world' $TEST_VAR",
			stdin:           []byte{},
			expectedStdout:  []byte("hello world test-value\n"),
			expectedSuccess: true,
		},
		{
			name:            "script with unquoted and quoted args",
			command:         "echo hello world \"$TEST_VAR\"",
			stdin:           []byte{},
			expectedStdout:  []byte("hello world test-value\n"),
			expectedSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptResult := RunKnownScript(tt.command, tt.stdin)
			result := scriptResult.Stdout
			success := scriptResult.Success
			if !reflect.DeepEqual(result, tt.expectedStdout) {
				t.Errorf("RunKnownScript(%v, %v) = %v, want %v",
					tt.command, tt.stdin, result, tt.expectedStdout)
			}
			if !reflect.DeepEqual(success, tt.expectedSuccess) {
				t.Errorf("RunKnownScript(%v, %v) = %v, want %v",
					tt.command, tt.stdin, success, tt.expectedSuccess)
			}
		})
	}
}
