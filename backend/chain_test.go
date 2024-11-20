package backend

import (
	// "fmt"
	"os"
	"reflect"
	"testing"
)

func TestGetChainStructure(t *testing.T) {
	cwd, _ := os.Getwd()
	testDir := cwd + "/test/chain/"
	testChainFileName := "test-chain"

	tests := []struct {
		name     string
		path     string
		expected []ChainItem
	}{
		{
			name: "test we get the chain structure",
			path: testDir,
			expected: []ChainItem{{Name: testChainFileName, Chain: Chain{
				{Name: "command_one", Command: "echo", Args: []string{"hello", "world"}},
				{Name: "command_two", Command: "cat", Args: []string{}},
			}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetChainStructure(testDir)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetChainStructure(%v) = %v, want %v",
					tt.path, result, tt.expected)
			}
		})
	}

}

func TestRunChain(t *testing.T) {
	tests := []struct {
		name                   string
		chain                  Chain
		stdin                  []byte
		expectedStdout         []byte
		expectedSuccess        bool
		expectedLastScriptName string
	}{
		{
			name: "simple chain with env var",
			chain: Chain{
				Script{
					Name:    "test-script",
					Command: "echo",
					Args:    []string{"hello", "world"},
				},
				Script{
					Name:    "test-script2",
					Command: "cat",
				},
			},
			stdin:                  []byte{},
			expectedStdout:         []byte("hello world\n"),
			expectedSuccess:        true,
			expectedLastScriptName: "test-script2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chainResult := RunChain(tt.stdin, tt.chain)
			lastScriptResult := chainResult[(len(chainResult) - 1)]
			result := lastScriptResult.Stdout
			success := lastScriptResult.Success
			lastScriptName := lastScriptResult.Script.Name
			if !reflect.DeepEqual(result, tt.expectedStdout) {
				t.Errorf("RunChain(%v, %v) = %v, want %v",
					tt.stdin, tt.chain, result, tt.expectedStdout)
			}
			if success != tt.expectedSuccess {
				t.Errorf("RunChain(%v, %v) expected success = %v, got %v",
					tt.stdin, tt.chain, success, tt.expectedSuccess)
			}
			if lastScriptName != tt.expectedLastScriptName {
				t.Errorf("RunChain(%v, %v) expected lastScriptName = %v, got %v",
					tt.stdin, tt.chain, tt.expectedLastScriptName, lastScriptName)
			}
		})
	}
}

func TestLoadChainThenRun(t *testing.T) {
	cwd, _ := os.Getwd()
	testDir := cwd + "/test/chain/"
	testChainFileName := "test-chain"

	tests := []struct {
		name                   string
		expectedStdout         []byte
		expectedSuccess        bool
		stdin                  []byte
		expectedLastScriptName string
	}{
		{
			name:                   "test we get the chain structure",
			expectedStdout:         []byte("hello world\n"),
			stdin:                  []byte{},
			expectedSuccess:        true,
			expectedLastScriptName: "command_two",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := LoadCustomChain(testDir, testChainFileName)
			chainResult := RunChain(tt.stdin, chain)
			lastScriptResult := chainResult[len(chainResult)-1]

			result := lastScriptResult.Stdout
			success := lastScriptResult.Success
			lastScriptName := lastScriptResult.Script.Name

			if !reflect.DeepEqual(result, tt.expectedStdout) {
				t.Errorf("RunChain(%v, %v) = %v, want %v",
					tt.stdin, chain, result, tt.expectedStdout)
			}
			if success != tt.expectedSuccess {
				t.Errorf("RunChain(%v, %v) expected success = %v, got %v",
					tt.stdin, chain, success, tt.expectedSuccess)
			}
			if lastScriptName != tt.expectedLastScriptName {
				t.Errorf("RunChain(%v, %v) expected lastScriptName = %v, got %v",
					tt.stdin, chain, tt.expectedLastScriptName, lastScriptName)
			}
		})
	}

}
