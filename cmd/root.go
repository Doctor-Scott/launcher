/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "fmt"
	// "fmt"
	"launcher/backend"
	C "launcher/globalConstants"
	"launcher/tui"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "launcher",
	Short: "A script launch pad",
	Long:  `This app works like a homepage for your scripts`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		// for flag := range flags {
		//
		// }

		tui.Start(path)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.launcher/launcher.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// rootCmd.Flags().StringP("path", "p", "", "optional path")
	rootCmd.PersistentFlags().StringP("path", "p", "", "optional path")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("useAndInDesc", C.USE_AND_IN_DESC)
	viper.SetDefault("clearChainAfterRun", C.CLEAR_CHAIN_AFTER_RUN)
	viper.SetDefault("autosave", C.AUTO_SAVE)
	viper.SetDefault("scriptTitleColor", C.SCRIPT_TITLE_COLOR)
	viper.SetDefault("chainTitleColor", C.CHAIN_TITLE_COLOR)
	viper.SetDefault("inputTitleColor", C.INPUT_TITLE_COLOR)
	viper.SetDefault("cursorColor", C.CURSOR_COLOR)
	viper.SetDefault("selectedScriptColor", C.SELECTED_SCRIPT_COLOR)
	viper.SetDefault("chainSeparator", C.CHAIN_SEPARATOR)
	viper.SetDefault("chainTotalSeparator", C.CHAIN_TOTAL_SEPARATOR)

	defaultLauncherDir := backend.ResolvePath("~") + ".launcher"
	defaultScriptDir := backend.ResolvePath("~") + ".scripts/launcherScripts/"
	viper.SetDefault("launcherDir", defaultLauncherDir)
	viper.SetDefault("scriptDir", defaultScriptDir)

	viper.SetConfigName("launcher")
	viper.SetConfigType("json")
	viper.AutomaticEnv() // read in environment variables that match

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if viperEnvConfigPath := os.Getenv("VIPER_CONFIG_PATH"); viperEnvConfigPath == "" {
		// Use config file from $ENV

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".launcher" (without extension).
		viper.AddConfigPath(home + "/.config")

	} else {
		viper.AddConfigPath(viperEnvConfigPath)

	}

	viper.ReadInConfig()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
