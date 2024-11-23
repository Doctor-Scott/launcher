package tui

import (
	"fmt"
	backend "github.com/Doctor-Scott/launcher/backend"
	C "github.com/Doctor-Scott/launcher/globalConstants"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	lists               lists
	stdout              []byte
	currentPath         string
	chain               backend.Chain
	views               views
	inputModel          inputModel
	lastFaildScriptName string
}

type ViewState int

const (
	ScriptsView ViewState = iota
	WorkflowsView
	InputView
)

var stateName = map[ViewState]string{
	ScriptsView:   "scripts",
	WorkflowsView: "workflows",
	InputView:     "input",
}

func (ss ViewState) String() string {
	return stateName[ss]
}

type views struct {
	currentView  ViewState
	previousView ViewState
}

type lists struct {
	scripts   list.Model
	workflows list.Model
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

func runLauncherCommandFromInput(m model, command string) (tea.Model, tea.Cmd) {
	backend.RunLauncherCommand(command)
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case inputFinishedMsg:
		switch m.inputModel.returnCommand {
		case C.RUN_SCRIPT:
			command := m.inputModel.textInput.Value()
			if command != "" {
				scriptResult := backend.RunKnownScript(command, m.stdout)
				//TODO  maybe have m hold a ChainResult?
				m.stdout = scriptResult.Stdout
			}
		case C.ADD_ARGS_TO_SCRIPT_AND_RUN:
			scriptArgs := m.inputModel.textInput.Value()
			script := backend.AddArgsToScript(m.lists.scripts.SelectedItem().(item).script, scriptArgs)

			scriptResult := backend.RunScript(script, m.stdout)
			m.stdout = scriptResult.Stdout
			m = maybeSetLastFailedScript(m, scriptResult)
		case C.ADD_SCRIPT_TO_CHAIN:
			command := m.inputModel.textInput.Value()
			if command != "" {
				script := backend.GetScriptFromCommand(command)
				m.chain = backend.AddScriptToChain(script, m.chain)
			}
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
		case C.ADD_ARGS_TO_SCRIPT_THEN_ADD_TO_CHAIN:
			script := backend.AddArgsToScript(m.lists.scripts.SelectedItem().(item).script, m.inputModel.textInput.Value())
			m.chain = backend.AddScriptToChain(script, m.chain)
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
		case C.SAVE_CUSTOM_CHAIN:
			name := m.inputModel.textInput.Value()
			backend.SaveCustomChain(m.chain, viper.GetString(C.PathConfig.LauncherDir.Name)+"/custom/", name)
			return m, nil

		case C.LOAD_CUSTOM_CHAIN:
			name := m.inputModel.textInput.Value()
			return loadCustomChain(m, name)
		case C.RUN_LAUNCHER_COMMAND:
			command := m.inputModel.textInput.Value()
			return runLauncherCommandFromInput(m, command)
		}

	case inputRejectedMsg:
		m.swapViews()
		return m, nil
	}

	return selectAndDisplayView(m, msg)
}

func selectAndDisplayView(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.inputModel.Selected {
		m.inputModel.Selected = false
		m.swapViews()
	}

	if m.views.currentView == ScriptsView {
		return scriptsUpdate(msg, m)
	}

	if m.views.currentView == WorkflowsView {
		return workflowsUpdate(msg, m)
	}

	inputModel, cmd := inputUpdate(m.inputModel, msg)
	m.inputModel = inputModel

	return m, cmd

}

// trying out a method receiver
func (m *model) swapViews() {
	state := m.views.currentView
	m.views.currentView = m.views.previousView
	m.views.previousView = state
}

func (m *model) newView(view ViewState) {
	alreadyInView := m.views.currentView == view
	if alreadyInView {
		return
	}
	m.views.previousView = m.views.currentView
	m.views.currentView = view
}

func (m model) View() string {
	if m.views.currentView == ScriptsView {
		return docStyle.Render(m.lists.scripts.View())
	}

	if m.views.currentView == WorkflowsView {
		return docStyle.Render(m.lists.workflows.View())
	}

	return inputView(m.inputModel)
}

func Start(path string) {
	path = backend.ResolvePath(path)

	m := model{
		currentPath: path,
		views:       views{currentView: ScriptsView},
		chain:       backend.ReadChainConfig(),
		stdout:      backend.ReadStdin(),
		lists:       lists{scripts: createScriptList(path), workflows: createChainList("")},
	}
	// fmt.Printf("Loaded chain: %+v\n", m.chain)

	// backend.SaveChain(m.chain)

	if !viper.GetBool(C.Autosave.Name) {
		backend.ClearAutoSave()
	}

	//TODO  The help screen could do with updating to show the keymaps
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
