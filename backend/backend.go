package backend

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	C "launcher/globalConstants"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ReadStdin() []byte {
	var stdin = []byte{}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdin = append(stdin, scanner.Bytes()...)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return stdin
}

func getFiles(scriptPath string) []string {
	scriptPath = ResolvePath(scriptPath)

	entries, err := os.ReadDir(scriptPath)
	if err != nil {
		log.Fatal(err)
	}
	files := []string{}

	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	return files
}

type Script struct {
	Name     string
	Path     string
	Args     []string
	Selected bool
}

func ResolvePath(path string) string {
	if path == "" {
		return os.Getenv("DEFAULT_SCRIPT_PATH")
	}

	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return path + "/"
}

func GetStructure(path string) []Script {
	path = ResolvePath(path)

	files := getFiles(path)
	scripts := []Script{}
	for _, file := range files {
		scripts = append(scripts, Script{Name: file, Path: path + file, Args: []string{}, Selected: false})
	}
	return scripts
}

func RunScript(script Script, stdin []byte) []byte {
	fmt.Println("Running", script.Name)
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

func saveStdout(stdout []byte) {
	fs, err := os.Create("/tmp/launcher.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer fs.Close()
	fs.Write(stdout)
	fs.Sync()
	fs.Close()
}

func RunInVim() []byte {
	cmd := exec.Command("nvim", "/tmp/launcher.txt")

	stdout, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
	}

	return stdout
}

func PrintStructure(path string) {
	for _, script := range GetStructure(path) {
		fmt.Println(script.Name)
	}
}

func RunChain(stdin []byte, chain []Script) []byte {
	if len(chain) == 0 {
		if C.CLEAR_CHAIN_AFTER_RUN {
			SaveChain(chain)
		}
		return stdin
	}
	stdout := RunScript(chain[0], stdin)
	return RunChain(stdout, chain[1:])

}

func AddScriptToChain(scriptToAdd Script, chain []Script) []Script {
	return SaveChain(append(chain, scriptToAdd))
}

func RemoveScriptFromChain(scriptToRemove Script, chain []Script) []Script {
	for i := len(chain) - 1; i >= 0; i-- {
		shouldRemoveScript := chain[i].Name == scriptToRemove.Name && chain[i].Path == scriptToRemove.Path
		shouldRemoveInput := chain[i].Name == C.INPUT_SCRIPT_NAME && scriptToRemove.Name == C.INPUT_SCRIPT_NAME
		if shouldRemoveScript || shouldRemoveInput {
			//pop the item
			return SaveChain(append(chain[0:i], chain[i+1:]...))
		}
	}
	// item not found in chain, so just return the chain
	return SaveChain(chain)
}

func SaveChain(chain []Script) []Script {
	viper.Set("chain", chain)
	viper.WriteConfig()
	return chain
}

func makeScriptFromItem(item interface{}) Script {
	var script Script
	if scriptMap, ok := item.(map[string]interface{}); ok {
		// Handle Args as a slice of strings
		args := []string{}
		if argsInterface, exists := scriptMap["args"]; exists && argsInterface != nil {
			if argsSlice, ok := argsInterface.([]interface{}); ok {
				for _, arg := range argsSlice {
					args = append(args, arg.(string))
				}
			}
		}

		script = Script{
			Path:     scriptMap["path"].(string),
			Name:     scriptMap["name"].(string),
			Args:     args,
			Selected: true,
		}
	}
	return script
}

func ReadChainConfig() []Script {
	// Safely handle chain configuration
	if chain := viper.Get("chain"); chain != nil {
		// Convert the interface{} slice to []backend.Script
		if chainSlice, ok := chain.([]interface{}); ok {

			scripts := make([]Script, 0, len(chainSlice))
			for _, item := range chainSlice {
				scripts = append(scripts, makeScriptFromItem(item))
			}
			return scripts
		}
	}
	return []Script{}

}

func main() {
	// getFiles()
	// for _, script := range GetStructure() {
	// 	fmt.Println(script.Name)
	// }

}
