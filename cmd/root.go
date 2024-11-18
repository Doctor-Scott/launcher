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
	Short: "Script launcher",
	Long:  `A homepage TUI and runner for your scripts, build complex workflows with ease!`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		// for flag := range flags {
		//
		// }
		//TODO Add a create default config command

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

func setDefaults() {
	viper.SetDefault(C.ClearChainAfterRun.Name, C.ClearChainAfterRun.DefaultValue)
	viper.SetDefault(C.Autosave.Name, C.Autosave.DefaultValue)

	//colors
	viper.SetDefault(C.ColorConfig.ScriptTitle.Name, C.ColorConfig.ScriptTitle.DefaultValue)
	viper.SetDefault(C.ColorConfig.ChainTitle.Name, C.ColorConfig.ChainTitle.DefaultValue)
	viper.SetDefault(C.ColorConfig.InputTitle.Name, C.ColorConfig.InputTitle.DefaultValue)
	viper.SetDefault(C.ColorConfig.Cursor.Name, C.ColorConfig.Cursor.DefaultValue)
	viper.SetDefault(C.ColorConfig.SelectedScript.Name, C.ColorConfig.SelectedScript.DefaultValue)

	//selected script description
	viper.SetDefault(C.SelectedScriptDescriptionConfig.UseAnd.Name, C.SelectedScriptDescriptionConfig.UseAnd.DefaultValue)
	viper.SetDefault(C.SelectedScriptDescriptionConfig.ChainSeparator.Name, C.SelectedScriptDescriptionConfig.ChainSeparator.DefaultValue)
	viper.SetDefault(C.SelectedScriptDescriptionConfig.ChainTotalSeparator.Name, C.SelectedScriptDescriptionConfig.ChainTotalSeparator.DefaultValue)
	viper.SetDefault(C.SelectedScriptDescriptionConfig.Prefix.Name, C.SelectedScriptDescriptionConfig.Prefix.DefaultValue)

	//paths
	viper.SetDefault(C.PathConfig.LauncherDir.Name, C.PathConfig.LauncherDir.DefaultValue)
	viper.SetDefault(C.PathConfig.ScriptDir.Name, C.PathConfig.ScriptDir.DefaultValue)

	//keybindings
	viper.SetDefault(C.KeybindingConfig.Item.RunUnderCursor.Name, C.KeybindingConfig.Item.RunUnderCursor.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Item.AddToChain.Name, C.KeybindingConfig.Item.AddToChain.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Script.AddArgsAndRun.Name, C.KeybindingConfig.Script.AddArgsAndRun.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Script.AddArgsThenAddToChain.Name, C.KeybindingConfig.Script.AddArgsThenAddToChain.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Script.RemoveFromChain.Name, C.KeybindingConfig.Script.RemoveFromChain.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Chain.RunChain.Name, C.KeybindingConfig.Chain.RunChain.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Chain.LoadKnown.Name, C.KeybindingConfig.Chain.LoadKnown.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Chain.LoadUnderCursor.Name, C.KeybindingConfig.Chain.LoadUnderCursor.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Chain.Write.Name, C.KeybindingConfig.Chain.Write.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Edit.OpenStdout.Name, C.KeybindingConfig.Edit.OpenStdout.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Edit.OpenEditor.Name, C.KeybindingConfig.Edit.OpenEditor.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Edit.OpenConfig.Name, C.KeybindingConfig.Edit.OpenConfig.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Edit.OpenItemUnderCursor.Name, C.KeybindingConfig.Edit.OpenItemUnderCursor.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.ClearState.Name, C.KeybindingConfig.ClearState.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.WriteConfig.Name, C.KeybindingConfig.WriteConfig.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.RefreshView.Name, C.KeybindingConfig.RefreshView.DefaultValue)
	viper.SetDefault(C.KeybindingConfig.Debug.Name, C.KeybindingConfig.Debug.DefaultValue)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	setDefaults()

	viper.SetConfigName("launcher")
	viper.SetConfigType("toml")
	viper.AutomaticEnv() // read in environment variables that match

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if viperEnvConfigPath := os.Getenv("LAUNCHER_CONFIG_PATH"); viperEnvConfigPath == "" {
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
