/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	viper.SetDefault(C.ClearChainAfterRun.Name, C.ClearChainAfterRun.DefaultValue)
	viper.SetDefault(C.Autosave.Name, C.Autosave.DefaultValue)
	viper.SetDefault(C.ColorConfig.ScriptTitle.Name, C.ColorConfig.ScriptTitle.DefaultValue)
	viper.SetDefault(C.ColorConfig.ChainTitle.Name, C.ColorConfig.ChainTitle.DefaultValue)
	viper.SetDefault(C.ColorConfig.InputTitle.Name, C.ColorConfig.InputTitle.DefaultValue)
	viper.SetDefault(C.ColorConfig.Cursor.Name, C.ColorConfig.Cursor.DefaultValue)
	viper.SetDefault(C.ColorConfig.SelectedScript.Name, C.ColorConfig.SelectedScript.DefaultValue)
	viper.SetDefault(C.ItemDescriptionConfig.UseAnd.Name, C.ItemDescriptionConfig.UseAnd.DefaultValue)
	viper.SetDefault(C.ItemDescriptionConfig.ChainSeparator.Name, C.ItemDescriptionConfig.ChainSeparator.DefaultValue)
	viper.SetDefault(C.ItemDescriptionConfig.ChainTotalSeparator.Name, C.ItemDescriptionConfig.ChainTotalSeparator.DefaultValue)

	viper.SetDefault(C.PathConfig.LauncherDir.Name, C.PathConfig.LauncherDir.DefaultValue)
	viper.SetDefault(C.PathConfig.ScriptDir.Name, C.PathConfig.ScriptDir.DefaultValue)

	viper.SetConfigName("launcher")
	viper.SetConfigType("toml")
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

	// viper.OnConfigChange(func(e fsnotify.Event) {
	// fmt.Println("Config file changed:", e.Name)
	// })
	viper.WatchConfig()
}
