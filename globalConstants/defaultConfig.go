package globalConstants

import "os"

type ConfigItem interface {
	GetName() string
	GetDefaultValue() interface{}
}

type stringConfigItem struct {
	Name         string
	DefaultValue string
}

func (s stringConfigItem) GetName() string {
	return s.Name
}

func (s stringConfigItem) GetDefaultValue() interface{} {
	return s.DefaultValue
}

type boolConfigItem struct {
	Name         string
	DefaultValue bool
}

func (b boolConfigItem) GetName() string {
	return b.Name
}

func (b boolConfigItem) GetDefaultValue() interface{} {
	return b.DefaultValue
}

// SetConfigValue sets a configuration value using the provided setter function
func SetConfigValue(setter func(string, interface{}), items ...ConfigItem) {
	for _, item := range items {
		setter(item.GetName(), item.GetDefaultValue())
	}
}

type pathConfig struct {
	LauncherDir stringConfigItem
	ScriptDir   stringConfigItem
}
type itemDescriptionConfig struct {
	UseAnd              boolConfigItem
	ChainSeparator      stringConfigItem
	ChainTotalSeparator stringConfigItem
	Prefix              stringConfigItem
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
	OpenEditor          stringConfigItem
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
		stringConfigItem{"keybindings.item.run_under_cursor", "enter"},
		stringConfigItem{"keybindings.item.add_to_chain", "a"},
	},
	scriptBindingsConfig{
		stringConfigItem{"keybindings.script.add_args_and_run", "space"},
		stringConfigItem{"keybindings.script.add_args_then_add_to_chain", "A"},
		stringConfigItem{"keybindings.script.remove_from_chain", "s"},
	},
	chainBindingsConfig{
		stringConfigItem{"keybindings.chain.run_chain", "R"},
		stringConfigItem{"keybindings.chain.load_known", "L"},
		stringConfigItem{"keybindings.chain.load_under_cursor", "l"},
		stringConfigItem{"keybindings.chain.delete_under_cursor", "D"},
		stringConfigItem{"keybindings.chain.write", "W"},
	},
	editBindingsConfig{
		stringConfigItem{"keybindings.edit.open_stdout", "v"},
		stringConfigItem{"keybindings.edit.open_editor", "n"},
		stringConfigItem{"keybindings.edit.open_config", "C"},
		stringConfigItem{"keybindings.edit.open_item_under_cursor", "e"},
	},
	stringConfigItem{"keybindings.clear_state", "c"},
	stringConfigItem{"keybindings.write_config", "U"},
	stringConfigItem{"keybindings.refresh_view", "r"},
	stringConfigItem{"keybindings.debug", "d"},
}

var ClearChainAfterRun = boolConfigItem{"clear_chain_after_run", false}
var Autosave = boolConfigItem{"autosave", true}

var path string = os.Getenv("HOME")
var PathConfig = pathConfig{
	stringConfigItem{"paths.launcher_dir", path + "/.launcher"},
	stringConfigItem{"paths.script_dir", path + "/.scripts/launcherScripts"},
}

var SelectedScriptDescriptionConfig = itemDescriptionConfig{
	boolConfigItem{"selected_script_description.use_and_in_description", false},
	stringConfigItem{"selected_script_description.chain_separator", ", "},
	stringConfigItem{"selected_script_description.chain_total_separator", " of "},
	stringConfigItem{"selected_script_description.prefix", "Position: "},
}

var ColorConfig = colorConfig{
	stringConfigItem{"colors.script_title", "#3300cc"},
	stringConfigItem{"colors.chain_title", "#c60062"},
	stringConfigItem{"colors.input_title", "#e64d00"},
	stringConfigItem{"colors.cursor", "#6fe6fc"},
	stringConfigItem{"colors.selected_script", "#6fe600"},
}

var DefaultItems = []ConfigItem{
	ClearChainAfterRun,
	Autosave,
	// Colors
	ColorConfig.ScriptTitle,
	ColorConfig.ChainTitle,
	ColorConfig.InputTitle,
	ColorConfig.Cursor,
	ColorConfig.SelectedScript,
	// Script Description
	SelectedScriptDescriptionConfig.UseAnd,
	SelectedScriptDescriptionConfig.ChainSeparator,
	SelectedScriptDescriptionConfig.ChainTotalSeparator,
	SelectedScriptDescriptionConfig.Prefix,
	// Paths
	PathConfig.LauncherDir,
	PathConfig.ScriptDir,
	// Keybindings
	KeybindingConfig.Item.RunUnderCursor,
	KeybindingConfig.Item.AddToChain,
	KeybindingConfig.Script.AddArgsAndRun,
	KeybindingConfig.Script.AddArgsThenAddToChain,
	KeybindingConfig.Script.RemoveFromChain,
	KeybindingConfig.Chain.RunChain,
	KeybindingConfig.Chain.LoadKnown,
	KeybindingConfig.Chain.LoadUnderCursor,
	KeybindingConfig.Chain.DeleteUnderCursor,
	KeybindingConfig.Chain.Write,
	KeybindingConfig.Edit.OpenStdout,
	KeybindingConfig.Edit.OpenEditor,
	KeybindingConfig.Edit.OpenConfig,
	KeybindingConfig.Edit.OpenItemUnderCursor,
	KeybindingConfig.ClearState,
	KeybindingConfig.WriteConfig,
	KeybindingConfig.RefreshView,
	KeybindingConfig.Debug,
}
