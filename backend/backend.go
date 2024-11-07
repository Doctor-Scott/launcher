package backend

import (
	"bufio"
	"bytes"
	"fmt"
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
		scripts = append(scripts, Script{Name: file, Path: path + file})
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
	scriptName, args := GetScriptNameAndArgs(command)

	cmd := exec.Command(scriptName, args...)

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

func GetScriptNameAndArgs(command string) (string, []string) {
	splitCommand := strings.Split(command, " ")

	scriptName := splitCommand[0]
	args := splitCommand[1:]
	//rejoin quoted args
	for i, arg := range args {
		if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") {
			args[i] = strings.Trim(arg, "\"")
		}
	}

	return scriptName, args
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

func removeScript(slice []Script, s int) []Script {
	return append(slice[:s], slice[s+1:]...)
}

func RunChain(stdin []byte, scripts []Script) []byte {
	if len(scripts) == 0 {
		return stdin
	}
	stdout := RunScript(scripts[0], stdin)
	return RunChain(stdout, removeScript(scripts, 0))

}

func AddScriptToChain(scriptToAdd Script, scripts []Script) []Script {
	scripts = append(scripts, scriptToAdd)
	return scripts

}

func RemoveScriptFromChain(scriptToRemove Script, scripts []Script) []Script {
	for i, script := range scripts {
		if script.Name == scriptToRemove.Name && script.Path == scriptToRemove.Path {
			return append(scripts[:i], scripts[i+1:]...)
		}
	}
	return scripts
}

func runScriptWithStdin(stdin []byte, script Script) []byte {
	cmd := exec.Command(script.Path)
	cmd.Stdin.Read(stdin)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(err)
	}
	return stdout

}

func main() {
	// getFiles()
	// for _, script := range GetStructure() {
	// 	fmt.Println(script.Name)
	// }

}
