package tui

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	backend "launcher/backend"
	C "launcher/globalConstants"
	"os"
	"os/exec"
	"strconv"
)

type item struct {
	title, titlePretty, desc string
	script                   backend.Script
	focused                  bool
}

func (i item) Title() string       { return i.titlePretty }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type vimFinishedMsg []byte
type updateStructureMsg bool
type generateSelectedItemViewMsg bool

const USE_AND_IN_DESC bool = false

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

func findScriptIndexes(chain []backend.Script, script backend.Script) []int {
	indexes := []int{}
	for i, chainScript := range chain {
		if chainScript.Name == script.Name {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func generatePositionString(indexes []int, chainLength int) string {
	desc := "Position: " + strconv.Itoa(indexes[0]+1)
	if len(indexes) != 1 {

		for i, index := range indexes {
			if i == 0 {
				continue
			}
			if USE_AND_IN_DESC {
				if i != len(indexes)-1 {
					desc += ", "
				} else {
					desc += " and "
				}
			} else {
				desc += ", "
			}

			desc += strconv.Itoa(index + 1)
		}
	}
	desc += " of " + strconv.Itoa(chainLength)
	return desc

}

func generateSelectedItemView(m model) model {
	if len(m.chain) == 0 {
		return emptySelected(m)
	}
	for i, listItem := range m.list.Items() {
		if item, ok := listItem.(item); ok {
			for _, chainScript := range m.chain {
				if item.script.Name == chainScript.Name {
					item.script.Selected = true
					indexes := findScriptIndexes(m.chain, item.script)
					desc := generatePositionString(indexes, len(m.chain))
					item.desc = desc

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
	fmt.Println("")
	fmt.Println(m.chain)
	// fmt.Println(string(m.stdout))
	fmt.Println("")
	// fmt.Print(m)

	// fmt.Println(m.Items())
	return m
}

func addArgsToScript(m model, scriptArgs string) model {
	script := m.list.SelectedItem().(item).script
	script.Args = append(script.Args, scriptArgs)
	m.list.SetItem(m.list.Index(), item{title: script.Name, script: script})
	return m
}

func listUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
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
			// standard run of known script or input command
			if m.list.SelectedItem().(item).title == "Input" {
				m.inputModel = initialInputModel("Script:", C.RUN_SCRIPT)
				m.currentView = "input"
				return m, nil
			} else {
				stdout := backend.RunScript(m.list.SelectedItem().(item).script, m.stdout)
				m.stdout = stdout
				cmd = func() tea.Msg {
					return tea.ClearScreen()
				}
				return m, cmd
			}
		}
		if msg.String() == tea.KeySpace.String() {
			// run script with args
			if m.list.SelectedItem().(item).title != "Input" {
				m.inputModel = initialInputModel("Args:", C.ADD_ARGS_TO_SCRIPT_AND_RUN)
				m.currentView = "input"
				cmd = func() tea.Msg {
					return tea.ClearScreen()
				}
				return m, cmd
			}
		}
		if msg.String() == "c" {
			// clear state
			m.chain = []backend.Script{}
			m.stdout = []byte{}
			m = generateSelectedItemView(m)
			return m, func() tea.Msg { return updateStructureMsg(true) }

		}
		if msg.String() == "d" {
			debug(m)
			return m, tea.Println(m)
		}
		if msg.String() == "a" {
			// add script
			if m.list.SelectedItem().(item).title == "Input" {
				m.inputModel = initialInputModel("Script:", C.ADD_SCRIPT_TO_CHAIN)
				m.currentView = "input"
				return m, nil
			} else {
				m.chain = backend.AddScriptToChain(m.list.SelectedItem().(item).script, m.chain)
				return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
			}

		}
		if msg.String() == "A" {
			// add script with args
			if m.list.SelectedItem().(item).title != "Input" {
				m.inputModel = initialInputModel("Args:", C.ADD_ARGS_TO_SCRIPT_THEN_ADD_TO_CHAIN)
				m.currentView = "input"
				return m, nil
			}

		}
		if msg.String() == "s" {
			// remove script from chain
			m.chain = backend.RemoveScriptFromChain(m.list.SelectedItem().(item).script, m.chain)
			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }

		}
		if msg.String() == "S" {
			//save chain
			// TODO
			// might be fun to allow people to share these chains once created
		}

		if msg.String() == "r" {
			// refresh view
			m.list.ResetSelected()
			return m, func() tea.Msg { return tea.ClearScreen() }

		}
		if msg.String() == "R" {
			// run chain
			stdout := backend.RunChain(m.stdout, m.chain)
			m.stdout = stdout
			cmd = func() tea.Msg {
				return tea.ClearScreen()
			}

			m.chain = []backend.Script{}

			return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }

		}
		if msg.String() == "e" {
			//edit file under cursor
			if m.list.SelectedItem().(item).title != "Input" {
				cmd := exec.Command("nvim", m.list.SelectedItem().(item).script.Path)
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					if err != nil {
						return fmt.Errorf("failed to run : %w", err)
					}
					return updateStructureMsg(true)
				})

			}
		}
		if msg.String() == "n" {
			// open nvim in launcher directory
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
			// open stdout from last script in editor
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
