package tui

import (
	"bytes"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	backend "launcher/backend"
	C "launcher/globalConstants"
	"os"
	"os/exec"
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

func listUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case updateStructureMsg:
		m = createNewModelList(m)
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
			m.chain = backend.Chain{}
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
				script := m.list.SelectedItem().(item).script
				script.Selected = true
				m.chain = backend.AddScriptToChain(script, m.chain)
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
		if msg.String() == "r" {
			// refresh view
			m.list.ResetSelected()
			return m, func() tea.Msg { return tea.ClearScreen() }

		}
		if msg.String() == "R" {
			// run chain
			stdout := backend.RunChain(m.stdout, m.chain)
			m.stdout = stdout

			if C.CLEAR_CHAIN_AFTER_RUN {
				m.chain = backend.Chain{}
			}

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
		if msg.String() == "L" {
			//load chain
			m.inputModel = initialInputModel("Name:", C.LOAD_CUSTOM_CHAIN)
			m.currentView = "input"
			return m, nil

		}

		if msg.String() == "W" {
			//write chain
			// TODO  might be fun to allow people to share these chains once created

			m.inputModel = initialInputModel("Name:", C.SAVE_CUSTOM_CHAIN)
			m.currentView = "input"
			return m, nil

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
