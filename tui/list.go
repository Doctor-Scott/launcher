package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	backend "launcher/backend"
	C "launcher/globalConstants"
)

type item struct {
	title, titlePretty, desc string
	script                   backend.Script
	selected                 bool
	focused                  bool
	chainItem                backend.ChainItem
}

func (i item) Title() string       { return i.titlePretty }
func (i item) Description() string { return i.desc }

// BUG  Filtering is not working, it is still interpreting the key commands
func (i item) FilterValue() string { return i.title }

type vimFinishedMsg []byte
type updateStructureMsg bool
type generateSelectedItemViewMsg bool

func chainsUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case updateStructureMsg:
		m = generateSelectedItemView(createNewChainModelList(m))
		backend.MaybeAutoSaveChain(m.chain)
		return m, tea.WindowSize()
	case generateSelectedItemViewMsg:
		m = generateSelectedItemView(m)
		return m, tea.WindowSize()
	case vimFinishedMsg:
		m.stdout = []byte(msg)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == tea.KeyTab.String() {
			return swapView(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Item.RunUnderCursor.Name) {
			return runItemUnderCursor(m, "chain")
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Item.AddToChain.Name) {
			return addScriptToChain(m, "chain")
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.ClearState.Name) {
			return clearState(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenConfig.Name) {
			return openConfig(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Debug.Name) {
			return debug(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Chain.DeleteUnderCursor.Name) {
			return deleteChainUnderCursor(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenItemUnderCursor.Name) {
			return editItemUnderCursor(m, "chain")
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Chain.LoadUnderCursor.Name) {
			m, _ := loadCustomChain(m, m.list.SelectedItem().(item).chainItem.Name)
			return swapView(m.(model))
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Chain.LoadKnown.Name) {
			return loadChain(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenNvim.Name) {
			return openNvimInLauncherDirectory(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.RefreshView.Name) {
			return refreshView(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Chain.RunChain.Name) {
			return runChain(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.WriteConfig.Name) {
			return writeConfig(m)

		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenStdout.Name) {
			return openWithVipe(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Chain.Write.Name) {
			return writeChain(m)
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func listUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case updateStructureMsg:
		m = generateSelectedItemView(createNewScriptModelList(m))
		backend.MaybeAutoSaveChain(m.chain)
		return m, tea.WindowSize()
	case generateSelectedItemViewMsg:
		m = generateSelectedItemView(m)
		return m, tea.WindowSize()
	case vimFinishedMsg:
		m.stdout = []byte(msg)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == tea.KeyTab.String() {
			return swapView(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Item.RunUnderCursor.Name) {
			return runItemUnderCursor(m, "script")
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Script.AddArgsAndRun.Name) {
			return runScriptWithArgs(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Item.AddToChain.Name) {
			return addScriptToChain(m, "script")
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Script.AddArgsThenAddToChain.Name) {
			return addScriptWithArgs(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.ClearState.Name) {
			return clearState(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenConfig.Name) {
			return openConfig(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Debug.Name) {
			return debug(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenItemUnderCursor.Name) {
			return editItemUnderCursor(m, "script")
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Chain.LoadKnown.Name) {
			return loadChain(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenNvim.Name) {
			return openNvimInLauncherDirectory(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.RefreshView.Name) {
			return refreshView(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Chain.RunChain.Name) {
			return runChain(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Script.RemoveFromChain.Name) {
			return removeScriptFromChain(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenStdout.Name) {
			return openWithVipe(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.WriteConfig.Name) {
			return writeConfig(m)

		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Chain.Write.Name) {
			return writeChain(m)
		}

		// if msg.String() == "x" {
		// 	// set chmod +x on script
		//
		// }
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
