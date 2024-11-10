package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	backend "launcher/backend"
	C "launcher/globalConstants"
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
				script.Selected = true
				m.chain = backend.AddScriptToChain(script, m.chain)
			}
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
		case C.ADD_ARGS_TO_SCRIPT_THEN_ADD_TO_CHAIN:
			script := backend.AddArgsToScript(m.list.SelectedItem().(item).script, m.inputModel.textInput.Value())
			script.Selected = true
			m.chain = backend.AddScriptToChain(script, m.chain)
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

	m := model{
		currentPath: path,
		currentView: "list",
		stdout:      backend.ReadStdin(),
		list:        createList(path),
	}


	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func main() {
	Main("")
}
