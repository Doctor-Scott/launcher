package tui_list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	backend "launcher/backend"
	tui_input "launcher/tui/input"
	"os"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list        list.Model
	stdout      []byte
	currentPath string
	chain       []backend.Script
	currentView string
	inputModel  tui_input.InputModel
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.currentView == "list" {
		return ListUpdate(msg, m)
	}
	return ListUpdate(msg, m)
}

func (m model) View() string {
	if m.currentView == "list" {
		return docStyle.Render(m.list.View())
	}
	return docStyle.Render(m.list.View())
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
