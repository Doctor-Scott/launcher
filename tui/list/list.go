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
	"strconv"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, titlePretty, desc string
	script                   backend.Script
	focused                  bool
}

func (i item) Title() string       { return i.titlePretty }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list        list.Model
	stdout      []byte
	currentPath string
	chain       []backend.Script
}

type vimFinishedMsg []byte
type updateStructureMsg bool
type generateSelectedItemViewMsg bool

func (m model) Init() tea.Cmd {
	return nil
}

func getCustomDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		// make each selected item a different color
		for i, listItem := range m.Items() {
			item := listItem.(item)

			if item.script.Selected == true {
				item.titlePretty = lipgloss.NewStyle().Foreground(lipgloss.Color("#6fe600")).Render(item.title)
				m.SetItem(i, item)
			} else {
				item.titlePretty = lipgloss.NewStyle().Foreground(lipgloss.NoColor{}).Render(item.title)
				m.SetItem(i, item)
			}
		}

		return nil
	}
	c := lipgloss.Color("#6fe6fc")
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(c).BorderLeftForeground(c)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle // reuse the title style here
	return delegate
}

func updateModelList(m model) model {
	structure := backend.GetStructure(m.currentPath)
	items := []list.Item{}

	items = append(items, item{title: "Input", desc: "Enter a script"})
	for _, script := range structure {
		items = append(items, item{title: script.Name, script: script})
	}
	delegate := getCustomDelegate()
	m.list = list.New(items, delegate, 0, 0)
	m.list.Title = "Running a script are we???"
	return m
}

func emptySelected(m model) model {
	for i, listItem := range m.list.Items() {
		if item, ok := listItem.(item); ok {
			item.script.Selected = false
			if item.title != "Input" {
				item.desc = ""
			}
			m.list.SetItem(i, item)
		}
	}
	return m
}

func generateSelectedItemView(m model) model {
	if len(m.chain) == 0 {
		return emptySelected(m)
	}
	for i, listItem := range m.list.Items() {
		if item, ok := listItem.(item); ok {
			for scriptIndex, chainScript := range m.chain {
				if item.script.Name == chainScript.Name {
					item.script.Selected = true

					item.desc = "selected " + strconv.Itoa(scriptIndex+1) + " of " + strconv.Itoa(len(m.chain))
					m.list.SetItem(i, item)
					break
				} else {
					item.script.Selected = false
					if item.title != "Input" {
						item.desc = ""
					}
					m.list.SetItem(i, item)
				}
			}
		}
	}

	return m
}

func debug(m model) model {
	// fmt.Println(m.currentPath)
	// fmt.Println("")
	// fmt.Println(m.list.Items())
	// fmt.Println("")
	// fmt.Println(string(m.stdout))
	fmt.Println("")
	fmt.Println(m.chain)
	fmt.Println("")
	// fmt.Print(m)

	// fmt.Println(m.Items())
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case updateStructureMsg:
		m = updateModelList(m)
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
			m.chain = []backend.Script{}
			m.stdout = []byte{}
			m = generateSelectedItemView(m)
			return m, func() tea.Msg { return updateStructureMsg(true) }

		}
		if msg.String() == "d" {
			debug(m)
			return m, nil
		}
		if msg.String() == "a" {
			m.chain = backend.AddScriptToChain(m.list.SelectedItem().(item).script, m.chain)
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }

		}
		if msg.String() == "s" {
			m.chain = backend.RemoveScriptFromChain(m.list.SelectedItem().(item).script, m.chain)
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }

		}

		if msg.String() == "r" {
			return m, func() tea.Msg { return tea.ClearScreen() }

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
					return updateStructureMsg(true)
				})

			}
		}
		if msg.String() == "n" {
			// fmt.Printf(m.currentPath)
			cmd := exec.Command("nvim", "--cmd", "cd"+m.currentPath+" | enew")
			m.list.ResetSelected()

			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				if err != nil {
					return fmt.Errorf("failed to run : %w", err)
				}
				return updateStructureMsg(true)
			})

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
	path = backend.ResolvePath(path)

	var m model
	m.currentPath = path
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
