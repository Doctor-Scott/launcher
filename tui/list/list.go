package tui_list

import (
	"fmt"
	"os"

	backend "launcher/backend"
	tui_input "launcher/tui/input"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
	script      backend.Script
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if m.list.SelectedItem().(item).title == "Input" {
				command := tui_input.Input("Script:")

				if command != "" {
					stdout := backend.RunKnownScript(command)
					// fmt.Scanln()
					fmt.Println(string(stdout))
				}
			} else {
				stdout := backend.RunScript(m.list.SelectedItem().(item).script)
				fmt.Println(string(stdout))
			}
			duration := 2 * time.Second
			time.Sleep(duration)
			cmd = func() tea.Msg {
				return tea.ClearScreen()
			}
			m.list.ResetSelected()
			return m, cmd
			// m.list, cmd = m.list.Update(msg)
			// return m, cmd
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Main(path string) {
	structure := backend.GetStructure(path)
	items := []list.Item{}

	items = append(items, item{title: "Input", desc: "Enter a script"})
	for _, script := range structure {
		items = append(items, item{title: script.Name, script: script})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Running a script are we???"

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
