package tui

import (
	"fmt"
	backend "launcher/backend"
	C "launcher/globalConstants"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list        list.Model
	stdout      []byte
	currentPath string
	chain       backend.Chain
	currentView string
	inputModel  inputModel
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return generateSelectedItemViewMsg(true) }
}

func loadCustomChain(m model, name string) (tea.Model, tea.Cmd) {
	path := viper.GetString(C.PathConfig.LauncherDir.Name) + "/custom/"
	m.chain = backend.LoadCustomChain(path, name)
	backend.MaybeAutoSaveChain(m.chain)
	return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case inputFinishedMsg:
		switch m.inputModel.returnCommand {
		case C.RUN_SCRIPT:
			command := m.inputModel.textInput.Value()
			if command != "" {
				m.stdout = backend.RunKnownScript(command, m.stdout)
			}
		case C.ADD_ARGS_TO_SCRIPT_AND_RUN:
			scriptArgs := m.inputModel.textInput.Value()
			script := backend.AddArgsToScript(m.list.SelectedItem().(item).script, scriptArgs)

			m.stdout = backend.RunScript(script, m.stdout)
		case C.ADD_SCRIPT_TO_CHAIN:
			command := m.inputModel.textInput.Value()
			if command != "" {
				script := backend.GetScriptFromCommand(command)
				m.chain = backend.AddScriptToChain(script, m.chain)
			}
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
		case C.ADD_ARGS_TO_SCRIPT_THEN_ADD_TO_CHAIN:
			script := backend.AddArgsToScript(m.list.SelectedItem().(item).script, m.inputModel.textInput.Value())
			m.chain = backend.AddScriptToChain(script, m.chain)
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
		case C.SAVE_CUSTOM_CHAIN:
			name := m.inputModel.textInput.Value()
			backend.SaveCustomChain(m.chain, viper.GetString(C.PathConfig.LauncherDir.Name)+"/custom/", name)
			return m, nil

		case C.LOAD_CUSTOM_CHAIN:
			name := m.inputModel.textInput.Value()
			return loadCustomChain(m, name)
		}

	case inputRejectedMsg:
		m.currentView = "list"
		return m, nil
	}

	if m.currentView == "list" || m.inputModel.Selected {
		m.inputModel.Selected = false
		m.currentView = "list"
		return listUpdate(msg, m)
	}

	if m.currentView == "chains" {
		return chainsUpdate(msg, m)
	}
	inputModel, cmd := inputUpdate(m.inputModel, msg)
	m.inputModel = inputModel

	return m, cmd
}

func (m model) View() string {
	if m.currentView == "list" || m.currentView == "chains" || m.inputModel.Selected {
		return docStyle.Render(m.list.View())
	}

	return inputView(m.inputModel)
}

func Start(path string) {
	path = backend.ResolvePath(path)

	m := model{
		currentPath: path,
		currentView: "list",
		chain:       backend.ReadChainConfig(),
		stdout:      backend.ReadStdin(),
		list:        createScriptList(path),
	}
	// fmt.Printf("Loaded chain: %+v\n", m.chain)

	// backend.SaveChain(m.chain)

	if !viper.GetBool(C.Autosave.Name) {
		backend.ClearAutoSave()
	}

	//TODO  The help screen could do with updating to show the keymaps
	//TODO  Set the keymaps as config options
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
