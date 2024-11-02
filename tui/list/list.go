package tui_list

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	backend "launcher/backend"
	tui_input "launcher/tui/input"
	"os"
	"os/exec"
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
	list   list.Model
	stdout []byte
}

type vimFinishedMsg []byte

func (m model) Init() tea.Cmd {
	return nil
}

func debug(m model) model {
	fmt.Println(m.currentPath)
	fmt.Println("")
	fmt.Println(m.list.Items())
	fmt.Println("")
	fmt.Println(string(m.stdout))
	fmt.Println("")
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case vimFinishedMsg:
		m.stdout = []byte(msg)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if m.list.SelectedItem().(item).title == "Input" {
				command := tui_input.Input("Script:")

				if command != "" {
					stdout := backend.RunKnownScript(command, m.stdout)
					m.stdout = stdout
				}
			} else {
				stdout := backend.RunScript(m.list.SelectedItem().(item).script, m.stdout)
				m.stdout = stdout
			}
			cmd = func() tea.Msg {
				return tea.ClearScreen()
			}
			m.list.ResetSelected()
			return m, cmd
		}
		if msg.String() == "c" {
			//clear stdout
			m.stdout = []byte{}

		}
		if msg.String() == "d" {
			m = debug(m)
			return m, nil
		}
		if msg.String() == "e" {
			//edit script
			if m.list.SelectedItem().(item).title != "Input" {
				cmd := exec.Command("nvim", m.list.SelectedItem().(item).script.Path)
				m.list.ResetSelected()
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					if err != nil {
						return fmt.Errorf("failed to run : %w", err)
					}
					return nil
				})

			}
		}
		if msg.String() == "v" {
			// run script in editor
			m.list.ResetSelected()
			cmd := exec.Command("vipe")
			cmd.Stdin = bytes.NewBuffer(m.stdout)
			var outBuf bytes.Buffer
			cmd.Stdout = &outBuf
			cmd.Stderr = os.Stderr

			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				if err != nil {
					return fmt.Errorf("failed to run vipe: %w", err)
				}
				return vimFinishedMsg(outBuf.Bytes())
			})
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
