package globalConstants

import "os"

type stringConfigItem struct {
	Name         string
	DefaultValue string
}
type boolConfigItem struct {
	Name         string
	DefaultValue bool
}

type pathConfig struct {
	LauncherDir stringConfigItem
	ScriptDir   stringConfigItem
}
type itemDescriptionConfig struct {
	UseAnd              boolConfigItem
	ChainSeparator      stringConfigItem
	ChainTotalSeparator stringConfigItem
}
type colorConfig struct {
	ScriptTitle    stringConfigItem
	ChainTitle     stringConfigItem
	InputTitle     stringConfigItem
	Cursor         stringConfigItem
	SelectedScript stringConfigItem
}

var ClearChainAfterRun = boolConfigItem{"clearChainAfterRun", false}
var Autosave = boolConfigItem{"autosave", true}

var path string = os.Getenv("HOME")
var PathConfig = pathConfig{
	stringConfigItem{"paths.launcherDir", path + "/.launcher/"},
	stringConfigItem{"paths.scriptDir", path + "/.scripts/launcherScripts/"},
}

var ItemDescriptionConfig = itemDescriptionConfig{
	boolConfigItem{"item_Description.useAndInDescription", false},
	stringConfigItem{"item_Description.chainSeparator", ", "},
	stringConfigItem{"item_Description.chainTotalSeparator", " of "},
}

var ColorConfig = colorConfig{
	stringConfigItem{"colors.scriptTitle", "#3300cc"},
	stringConfigItem{"colors.chainTitle", "#c60062"},
	stringConfigItem{"colors.inputTitle", "#e64d00"},
	stringConfigItem{"colors.cursor", "#6fe6fc"},
	stringConfigItem{"colors.selectedScript", "#6fe600"},
}
