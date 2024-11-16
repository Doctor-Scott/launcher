package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	backend "launcher/backend"
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
		if msg.String() == "enter" {
			return runItemUnderCursor(m, "chain")
		}
		if msg.String() == "a" {
			return addItemToChain(m, "chain")
		}
		if msg.String() == "c" {
			return clearState(m)
		}
		if msg.String() == "C" {
			return openConfig(m)
		}
		if msg.String() == "d" {
			return debug(m)
		}
		if msg.String() == "D" {
			return deleteChainUnderCursor(m)
		}
		if msg.String() == "e" {
			return editItemUnderCursor(m, "chain")
		}
		if msg.String() == "l" {
			m, _ := loadCustomChain(m, m.list.SelectedItem().(item).chainItem.Name)
			return swapView(m.(model))
		}
		if msg.String() == "L" {
			return loadChain(m)
		}
		if msg.String() == "n" {
			return openNvimInLauncherDirectory(m)
		}
		if msg.String() == "r" {
			return refreshView(m)
		}
		if msg.String() == "R" {
			return runChain(m)
		}
		if msg.String() == "U" {
			return writeConfig(m)

		}
		if msg.String() == "v" {
			return openWithVipe(m)
		}
		if msg.String() == "W" {
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
		if msg.String() == "enter" {
			return runItemUnderCursor(m, "script")
		}
		if msg.String() == tea.KeySpace.String() {
			return runScriptWithArgs(m)
		}
		if msg.String() == "a" {
			return addItemToChain(m, "script")
		}
		if msg.String() == "A" {
			return addScriptWithArgs(m)
		}
		if msg.String() == "c" {
			return clearState(m)
		}
		if msg.String() == "C" {
			return openConfig(m)
		}
		if msg.String() == "d" {
			return debug(m)
		}
		if msg.String() == "e" {
			return editItemUnderCursor(m, "script")
		}
		if msg.String() == "L" {
			return loadChain(m)
		}
		if msg.String() == "n" {
			return openNvimInLauncherDirectory(m)
		}
		if msg.String() == "r" {
			return refreshView(m)
		}
		if msg.String() == "R" {
			return runChain(m)
		}
		if msg.String() == "s" {
			return removeScriptFromChain(m)
		}
		if msg.String() == "v" {
			return openWithVipe(m)
		}
		if msg.String() == "U" {
			return writeConfig(m)

		}
		if msg.String() == "W" {
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
