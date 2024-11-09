package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	backend "launcher/backend"
	"os"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list        list.Model
	stdout      []byte
	currentPath string
	chain       []backend.Script
	currentView string
	inputModel  inputModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg.(type) {
	case inputFinishedMsg:
		switch m.inputModel.inputType {
		case "runScript":
			command := m.inputModel.textInput.Value()
			if command != "" {
				stdout := backend.RunKnownScript(command, m.stdout)
				m.stdout = stdout
			}
		case "addArgsToScriptAndRun":
			scriptArgs := m.inputModel.textInput.Value()
			m = addArgsToScript(m, scriptArgs)

			stdout := backend.RunScript(m.list.SelectedItem().(item).script, m.stdout)
			m.stdout = stdout
		case "addScriptToChain":
			command := m.inputModel.textInput.Value()
			if command != "" {
				scriptName, args := backend.GetScriptNameAndArgs(command)
				script := backend.Script{Name: scriptName, Path: scriptName, Args: args}
				m.chain = backend.AddScriptToChain(script, m.chain)
			}
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
		case "addArgsToScriptThenAddToChain":
			scriptArgs := m.inputModel.textInput.Value()
			m = addArgsToScript(m, scriptArgs)
			m.chain = backend.AddScriptToChain(m.list.SelectedItem().(item).script, m.chain)
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
		}

	case inputRejectedMsg:
		m.currentView = "list"
		return m, nil
	}

	if m.currentView == "list" || m.inputModel.Selected {
		return listUpdate(msg, m)
	}
	inputModel, cmd := inputUpdate(m.inputModel, msg)
	m.inputModel = inputModel

	return m, cmd
}

func (m model) View() string {
	if m.currentView == "list" || m.inputModel.Selected {
		return docStyle.Render(m.list.View())
	}
	return inputView(m.inputModel)
}

func Main(path string) {
	path = backend.ResolvePath(path)

	var m model
	m.currentPath = path
	m.currentView = "list"
	m = updateModelList(m)

	m.stdout = backend.ReadStdin()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func main() {
	Main("")
}