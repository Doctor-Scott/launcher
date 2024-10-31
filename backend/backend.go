package backend

import (
	// "bytes"
	"fmt"
	// "io"
	"log"
	"os"
	"os/exec"
)

// const DEFAULT_SCRIPT_PATH string = "/Users/felix/.scripts/"

func getFiles(scriptPath string) []string {

	if scriptPath == "" {
		scriptPath = os.Getenv("DEFAULT_SCRIPT_PATH")
	}

	entries, err := os.ReadDir(scriptPath)
	if err != nil {
		log.Fatal(err)
	}
	files := []string{}

	for _, e := range entries {
		// fmt.Println(e.Name())
		if !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	return files
}

type Script struct {
	Name string
	Path string
	Args []string
}

func GetStructure(path string) []Script {
	files := getFiles(path)
	scripts := []Script{}
	for _, file := range files {
		scripts = append(scripts, Script{Name: file, Path: path + file})
	}
	return scripts
}

func RunScript(script Script) []byte {
	//https://www.youtube.com/watch?v=jYqFUdFUej4&t=258s
	fmt.Println("Running", script.Name)
	cmd := exec.Command(script.Path)

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(output))
	// new_output, err := exec.Command("/bin/zsh", string(output)+" | vim -").Output()

	// fs, err := os.Create("file.txt")
	// defer fs.Close()
	// fs.Write(stdout)
	// fs.Sync()
	// fs.Close()

	// cmd := exec.Command("nvim", "-u", "NONE", "-", "+set noswapfile")
	// buffer := bytes.NewBuffer(output)
	// cmd.Stdin = buffer
	// fmt.Println("we are running nvim")
	// // cmd.Stdin = os.Stdin
	//
	// // var outputBuffer bytes.Buffer
	// outputBuffer := os.Stdout
	// cmd.Stdout = outputBuffer
	//
	// // cmd.Stdout =
	//
	// // err = cmd.Run()
	// fmt.Println(outputBuffer)
	// fmt.Println("we are done running nvim")

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if _, err := io.Copy(os.Stdout, outputBuffer); err != nil {
	// 	log.Fatalf("Error writing to stdout: %v", err)
	// }

	return stdout
}

func RunKnownScript(scriptName string) []byte {
	cmd := exec.Command(scriptName)
	stdout, err := cmd.CombinedOutput()
	// fs, err := os.Create("file.txt")
	// defer fs.Close()
	// fs.Write(output)
	// fs.Sync()
	// fs.Close()

	// cmd := exec.Command("nvim", "-u", "NONE", "-", "+set noswapfile")
	// buffer := bytes.NewBuffer(output)
	// cmd.Stdin = buffer
	// fmt.Println("we are running nvim")
	// // cmd.Stdin = os.Stdin
	//
	// // var outputBuffer bytes.Buffer
	// outputBuffer := os.Stdout
	// cmd.Stdout = outputBuffer
	//
	// // cmd.Stdout =
	//
	// err = cmd.Run()
	// fmt.Println(outputBuffer)
	// fmt.Println("we are done running nvim")
	//
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if _, err := io.Copy(os.Stdout, outputBuffer); err != nil {
	// 	log.Fatalf("Error writing to stdout: %v", err)
	// }

	// RunScript("", scriptName)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("\n", string(output))
	return stdout
}

func PrintStructure(path string) {
	for _, script := range GetStructure(path) {
		fmt.Println(script.Name)
	}
}

// NOTE
// you could have 2 options for running a workflow
// 1 that runs sync, manually taking the full output of one and putting it into another
// and one that builds the scipt as a chain
// need to make sure launcher can take stdin and pass it to the script
func joinScripts(stdin []byte, scripts []Script) {

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
