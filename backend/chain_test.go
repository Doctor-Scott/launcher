package backend

import (
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
		name     string
		chain    Chain
		stdin    []byte
		expected []byte
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
			stdin:    []byte{},
			expected: []byte("hello world\n"),
		}, {}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RunChain(tt.stdin, tt.chain)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RunChain(%v, %v) = %v, want %v",
					tt.stdin, tt.chain, result, tt.expected)
			}
		})
	}
}

func TestLoadChainThenRun(t *testing.T) {
	cwd, _ := os.Getwd()
	testDir := cwd + "/test/chain/"
	testChainFileName := "test-chain"

	tests := []struct {
		name     string
		expected []byte
		stdin    []byte
	}{
		{
			name:     "test we get the chain structure",
			expected: []byte("hello world\n"),
			stdin:    []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := LoadCustomChain(testDir, testChainFileName)
			result := RunChain(tt.stdin, chain)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RunChain(%v, %v) = %v, want %v",
					tt.stdin, chain, result, tt.expected)
			}
		})
	}

}
