package backend

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	C "github.com/Doctor-Scott/launcher/globalConstants"
	"github.com/spf13/viper"
)

func ResolvePath(path string) string {
	path = os.ExpandEnv(path)
	if path == "" {
		return viper.GetString(C.PathConfig.ScriptDir.Name)
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

func RunLauncherCommand(commandInput string) {
	fmt.Printf(commandInput)
	time.Sleep(1000)
	splitCommand := strings.Split(commandInput, " ")
	command := splitCommand[0]
	args := splitCommand[1:]
	if command == "genFreshConf" {
		if len(args) > 1 {
			// TODO  handle errors
			// maybe show it back in the command view
			return
		}
		if len(args) == 0 {
			// default to current directory
			args = append(args, ".")
		}
		path := args[0]

		configName := "launcher.toml"
		pathAndFileName := path + "/" + configName
		currentConfigPath := viper.ConfigFileUsed()
		fmt.Printf(currentConfigPath)
		resetDefaults()
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// standard directory perms
			os.Mkdir(path, 0755)
		}
		error := viper.WriteConfigAs(pathAndFileName)

		if error != nil {
			panic(error)

		}

		viper.SetConfigFile(currentConfigPath)
		viper.ReadInConfig()
	}

}

func resetDefaults() {
	C.SetConfigValue(viper.Set, C.DefaultItems...)
}

func SetDefaults() {
	C.SetConfigValue(viper.SetDefault, C.DefaultItems...)
}
