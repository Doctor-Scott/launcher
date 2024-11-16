package backend

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	C "launcher/globalConstants"
)

func ResolvePath(path string) string {
	path = os.ExpandEnv(path)
	if path == "" {
		return viper.GetString(C.ScriptDir.Name)
	}
	if path == "~" {
		path = os.Getenv("HOME")
	}

	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		log.Println(err)
	}
	return path + "/"
}

func PrintStructure(path string) {
	for _, script := range GetStructure(path) {
		fmt.Println(script.Name)
	}
}
