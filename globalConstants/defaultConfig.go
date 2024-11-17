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

type itemBindingsConfig struct {
	RunUnderCursor stringConfigItem
	AddToChain     stringConfigItem
}

type scriptBindingsConfig struct {
	AddArgsAndRun         stringConfigItem
	AddArgsThenAddToChain stringConfigItem
	RemoveFromChain       stringConfigItem
}

type chainBindingsConfig struct {
	RunChain          stringConfigItem
	LoadKnown         stringConfigItem
	LoadUnderCursor   stringConfigItem
	DeleteUnderCursor stringConfigItem
	Write             stringConfigItem
}

type editBindingsConfig struct {
	OpenStdout          stringConfigItem
	OpenNvim            stringConfigItem
	OpenConfig          stringConfigItem
	OpenItemUnderCursor stringConfigItem
}

type keybindingConfig struct {
	// both
	Item itemBindingsConfig

	// script
	Script scriptBindingsConfig

	//chain
	Chain chainBindingsConfig

	// Edit
	Edit editBindingsConfig

	//other
	ClearState  stringConfigItem
	WriteConfig stringConfigItem
	RefreshView stringConfigItem

	Debug stringConfigItem
}

var KeybindingConfig = keybindingConfig{
	itemBindingsConfig{
		stringConfigItem{"keybindings.item.runUnderCursor", "enter"},
		stringConfigItem{"keybindings.item.addToChain", "a"},
	},
	scriptBindingsConfig{
		stringConfigItem{"keybindings.script.addArgsAndRun", "space"},
		stringConfigItem{"keybindings.script.addArgsThenAddToChain", "A"},
		stringConfigItem{"keybindings.script.removeFromChain", "s"},
	},
	chainBindingsConfig{
		stringConfigItem{"keybindings.chain.runChain", "R"},
		stringConfigItem{"keybindings.chain.loadKnown", "L"},
		stringConfigItem{"keybindings.chain.loadUnderCursor", "l"},
		stringConfigItem{"keybindings.chain.deleteUnderCursor", "D"},
		stringConfigItem{"keybindings.chain.write", "W"},
	},
	editBindingsConfig{
		stringConfigItem{"keybindings.edit.openStdout", "v"},
		stringConfigItem{"keybindings.edit.openNvim", "n"},
		stringConfigItem{"keybindings.edit.openConfig", "C"},
		stringConfigItem{"keybindings.edit.openItemUnderCursor", "e"},
	},
	stringConfigItem{"keybindings.clearState", "c"},
	stringConfigItem{"keybindings.writeConfig", "U"},
	stringConfigItem{"keybindings.refreshView", "r"},
	stringConfigItem{"keybindings.debug", "d"},
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
