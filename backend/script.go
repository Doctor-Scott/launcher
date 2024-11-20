package backend

import (
	"bytes"
	"fmt"
	C "github.com/Doctor-Scott/launcher/globalConstants"
	"os"
	"os/exec"
	"strings"
)

type Script struct {
	Name    string
	Command string
	Args    []string
}

type ScriptResult struct {
	Script Script
	//TODO  use the exit code to determine success
	Success bool
	Stdout  []byte
	Stderr  error
}

func GetStructure(path string) []Script {
	path = ResolvePath(path)

	files := getFiles(path)
	scripts := []Script{}
	for _, file := range files {
		scripts = append(scripts, Script{Name: file, Command: path + file, Args: []string{}})
	}
	return scripts
}

func GetScriptFromCommand(command string) Script {
	scriptName, argsString, found := strings.Cut(command, " ")
	name := C.INPUT_COMMAND_NAME

	if found && len(argsString) != 0 {
		return Script{Name: name, Command: scriptName, Args: resolveArgsString(argsString)}
	}
	// command with no args
	return Script{Name: name, Command: scriptName, Args: []string{}}
}

func resolveArgsString(argsString string) []string {
	if len(argsString) == 0 {
		return []string{""}
	}

	var args []string
	var currentArg strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	// Parse character by character to handle quotes properly
	for _, char := range argsString {
		switch {
		case char == '"' || char == '\'':
			if inQuotes && char == quoteChar {
				// End of quoted section
				inQuotes = false
				arg := currentArg.String()
				if quoteChar == '"' {
					// Only expand env vars in double quotes
					arg = os.ExpandEnv(arg)
				}
				args = append(args, arg)
				currentArg.Reset()
				quoteChar = 0
			} else if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else {
				currentArg.WriteRune(char)
			}
		case char == ' ' && !inQuotes:
			if currentArg.Len() > 0 {
				// Unquoted argument
				arg := os.ExpandEnv(currentArg.String())
				args = append(args, arg)
				currentArg.Reset()
			}
		default:
			currentArg.WriteRune(char)
		}
	}

	// Add the last argument if there is one
	if currentArg.Len() > 0 {
		arg := currentArg.String()
		if !inQuotes {
			// Unquoted argument
			arg = os.ExpandEnv(arg)
		} else if quoteChar == '"' {
			// Double quoted argument
			arg = os.ExpandEnv(arg)
		}
		args = append(args, arg)
	}
	return args
}

func AddArgsToScript(script Script, argsString string) Script {
	args := resolveArgsString(argsString)
	script.Args = append(script.Args, args...)
	// script.Args = append(script.Args, scriptArgs)
	return script
}

func RunScript(script Script, stdin []byte) ScriptResult {
	cmd := exec.Command(script.Command, script.Args...)

	if len(stdin) > 0 {
		stdinBuffer := bytes.NewBuffer(stdin)
		cmd.Stdin = stdinBuffer
	}

	stdout, err := cmd.CombinedOutput()
	saveStdout(stdout)
	success := err == nil
	if !success {
		fmt.Println(err)
	}
	return ScriptResult{Script: script, Success: success, Stdout: stdout, Stderr: err}
}

func RunKnownScript(command string, stdin []byte) ScriptResult {
	return RunScript(GetScriptFromCommand(command), stdin)
}
