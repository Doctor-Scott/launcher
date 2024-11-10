package backend

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ResolvePath(path string) string {
	if path == "" {
		return os.Getenv("DEFAULT_SCRIPT_PATH")
	}
	if path == "~" {
		path = os.Getenv("HOME")
	}

	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return path + "/"
}

func PrintStructure(path string) {
	for _, script := range GetStructure(path) {
		fmt.Println(script.Name)
	}
}
