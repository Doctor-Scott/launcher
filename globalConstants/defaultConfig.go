package globalConstants

import "os"

type configItem struct {
	Name         string
	DefaultValue any
}

var ClearChainAfterRun = configItem{"clearChainAfterRun", false}

var Autosave = configItem{"autosave", true}
var UseAndInDescription = configItem{"useAndInDescription", false}
var ChainSeparator = configItem{"chainSeparator", ", "}
var ChainTotalSeparator = configItem{"chainTotalSeparator", " of "}
var ScriptTitleColor = configItem{"scriptTitleColor", "#3300cc"}
var ChainTitleColor = configItem{"chainTitleColor", "#c60062"}
var InputTitleColor = configItem{"inputTitleColor", "#e64d00"}
var CursorColor = configItem{"cursorColor", "#6fe6fc"}
var SelectedScriptColor = configItem{"selectedScriptColor", "#6fe600"}

var path string = os.Getenv("HOME")
var LauncherDir = configItem{"launcherDir", path + "/.launcher"}
var ScriptDir = configItem{"scriptDir", path + "/.scripts/launcherScripts/"}
