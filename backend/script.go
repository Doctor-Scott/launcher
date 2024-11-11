package backend

import (
	"bytes"
	"fmt"
	C "launcher/globalConstants"
	"os"
	"os/exec"
	"strings"
)

type Script struct {
	Name string
	Path string
	Args []string
}

func GetStructure(path string) []Script {
	path = ResolvePath(path)

	files := getFiles(path)
	scripts := []Script{}
	for _, file := range files {
		scripts = append(scripts, Script{Name: file, Path: path + file, Args: []string{}})
	}
	return scripts
}

func GetScriptFromCommand(command string) Script {
	scriptName, argsString, found := strings.Cut(command, " ")
	name := C.INPUT_SCRIPT_NAME

	if found && len(argsString) != 0 {
		return Script{Name: name, Path: scriptName, Args: resolveArgsString(argsString)}
	}
	// command with no args
	return Script{Name: name, Path: scriptName, Args: []string{}}
}

func resolveArgsString(argsString string) []string {
	args := strings.Split(argsString, " ")
	//rejoin quoted args and expand environment variables
	for i, arg := range args {
		// BUG  I think this needs to be adjusted
		// a string with spaces in the quotes would not work here I dont think
		//TODO  good oppertunity for a test
		if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") {
			args[i] = strings.Trim(arg, "\"")
		}
		if strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'") {
			args[i] = strings.Trim(arg, "'")
		}
		args[i] = os.ExpandEnv(arg)
	}
	return args
}

func AddArgsToScript(script Script, argsString string) Script {
	args := resolveArgsString(argsString)
	script.Args = append(script.Args, args...)
	// script.Args = append(script.Args, scriptArgs)
	return script
}

func RunScript(script Script, stdin []byte) []byte {
	cmd := exec.Command(script.Path, script.Args...)

	if len(stdin) > 0 {
		stdinBuffer := bytes.NewBuffer(stdin)
		cmd.Stdin = stdinBuffer
	}

	stdout, err := cmd.CombinedOutput()
	saveStdout(stdout)
	if err != nil {
		fmt.Println(err)
	}
	return stdout
}

func RunKnownScript(command string, stdin []byte) []byte {
	return RunScript(GetScriptFromCommand(command), stdin)
}
