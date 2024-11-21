package tui

import (
	backend "github.com/Doctor-Scott/launcher/backend"
	C "github.com/Doctor-Scott/launcher/globalConstants"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type item struct {
	title, titlePretty, desc string
	script                   backend.Script
	selected                 bool
	failed                   bool
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

func workflowsUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case updateStructureMsg:
		m = m.setSelectedScriptsInView()
		backend.MaybeAutoSaveChain(m.chain)
		return m, tea.WindowSize()
	case generateSelectedItemViewMsg:
		m = m.setSelectedScriptsInView()
		return m, tea.WindowSize()
	case vimFinishedMsg:
		m.stdout = []byte(msg)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == tea.KeyTab.String() {
			return scriptView(m)
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
			m, _ := loadCustomChain(m, m.lists.workflows.SelectedItem().(item).chainItem.Name)
			return scriptView(m.(model))
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Chain.LoadKnown.Name) {
			return loadChain(m)
		}
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenEditor.Name) {
			return openEditorInLauncherDirectory(m)
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
		m.lists.workflows.SetSize(msg.Width-h, msg.Height-v)
	}

	m.lists.workflows, cmd = m.lists.workflows.Update(msg)
	return m, cmd

}

func scriptsUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case updateStructureMsg:
		m = m.setSelectedScriptsInView()
		backend.MaybeAutoSaveChain(m.chain)
		return m, tea.WindowSize()
	case generateSelectedItemViewMsg:
		m = m.setSelectedScriptsInView()
		return m, tea.WindowSize()
	case vimFinishedMsg:
		m.stdout = []byte(msg)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == tea.KeyTab.String() {
			return workflowView(m)
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
		if msg.String() == viper.GetString(C.KeybindingConfig.Edit.OpenEditor.Name) {
			return openEditorInLauncherDirectory(m)
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
		m.lists.scripts.SetSize(msg.Width-h, msg.Height-v)
	}

	m.lists.scripts, cmd = m.lists.scripts.Update(msg)
	return m, cmd
}
