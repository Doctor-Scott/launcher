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

func debug(m model) (tea.Model, tea.Cmd) {
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
	return m, nil
}

func swapView(m model) (tea.Model, tea.Cmd) {
	// swap view between script and chains
	if m.currentView == "list" {
		m.list.ResetSelected()
		m.list = createChainList()
		m.currentView = "chains"
	} else {
		m.list.ResetSelected()
		m.list = createScriptList(m.currentPath)
		m.currentView = "list"
	}
	return m, func() tea.Msg { return updateStructureMsg(true) }
}

func addItemToChain(m model, itemType string) (tea.Model, tea.Cmd) {
	if m.list.SelectedItem().(item).title == "Input" {
		m.inputModel = initialInputModel("Script:", C.ADD_SCRIPT_TO_CHAIN)
		m.currentView = "input"
		return m, nil
	} else {
		if itemType == "chain" {
			m.chain = backend.AddChainToChain(m.list.SelectedItem().(item).chainItem.Chain, m.chain)
			return swapView(m)
		}
		script := m.list.SelectedItem().(item).script
		m.chain = backend.AddScriptToChain(script, m.chain)
		return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
	}
}

func refreshView(m model) (tea.Model, tea.Cmd) {
	m.list.ResetSelected()
	return m, func() tea.Msg { return tea.ClearScreen() }
}

func runChain(m model) (tea.Model, tea.Cmd) {
	stdout := backend.RunChain(m.stdout, m.chain)
	m.stdout = stdout

	if C.CLEAR_CHAIN_AFTER_RUN {
		m.chain = backend.Chain{}
	}

	return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
}

func editItemUnderCursor(m model, itemType string) (tea.Model, tea.Cmd) {
	if m.list.SelectedItem().(item).title != "Input" {
		var pathToChain string
		if itemType == "chain" {
			pathToChain = backend.ResolvePath("~") + ".launcher/custom/" + m.list.SelectedItem().(item).chainItem.Name + ".json"
		} else {
			pathToChain = m.list.SelectedItem().(item).script.Path
		}

		cmd := exec.Command("nvim", pathToChain)
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return fmt.Errorf("failed to run : %w", err)
			}
			return updateStructureMsg(true)
		})

	}
	return m, nil

}

func openNvimInLauncherDirectory(m model) (tea.Model, tea.Cmd) {
	cmd := exec.Command("nvim", "--cmd", "cd"+m.currentPath+" | enew")
	m.list.ResetSelected()

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return fmt.Errorf("failed to run : %w", err)
		}
		return updateStructureMsg(true)
	})
}

func openWithVipe(m model) (tea.Model, tea.Cmd) {
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

func loadChain(m model) (tea.Model, tea.Cmd) {
	m.inputModel = initialInputModel("Name:", C.LOAD_CUSTOM_CHAIN)
	m.currentView = "input"
	return m, nil
}

func writeChain(m model) (tea.Model, tea.Cmd) {
	m.inputModel = initialInputModel("Name:", C.SAVE_CUSTOM_CHAIN)
	m.currentView = "input"
	return m, nil
}

func clearState(m model) (tea.Model, tea.Cmd) {
	m.chain = backend.Chain{}
	m.stdout = []byte{}
	m = generateSelectedItemView(m)
	return m, func() tea.Msg { return updateStructureMsg(true) }
}

func runItemUnderCursor(m model, itemType string) (tea.Model, tea.Cmd) {
	if itemType == "chain" {
		stdout := backend.RunChain(m.stdout, m.list.SelectedItem().(item).chainItem.Chain)
		m.stdout = stdout
		cmd := func() tea.Msg {
			return tea.ClearScreen()
		}
		return m, cmd
	}

	// standard run of known script or input command
	if m.list.SelectedItem().(item).title == "Input" {
		m.inputModel = initialInputModel("Script:", C.RUN_SCRIPT)
		m.currentView = "input"
		return m, nil
	} else {
		stdout := backend.RunScript(m.list.SelectedItem().(item).script, m.stdout)
		m.stdout = stdout
		cmd := func() tea.Msg {
			return tea.ClearScreen()
		}
		return m, cmd
	}

}

func runScriptWithArgs(m model) (tea.Model, tea.Cmd) {
	if m.list.SelectedItem().(item).title != "Input" {
		m.inputModel = initialInputModel("Args:", C.ADD_ARGS_TO_SCRIPT_AND_RUN)
		m.currentView = "input"
		cmd := func() tea.Msg {
			return tea.ClearScreen()
		}
		return m, cmd
	}
	return m, nil
}

func addScriptWithArgs(m model) (tea.Model, tea.Cmd) {
	if m.list.SelectedItem().(item).title != "Input" {
		m.inputModel = initialInputModel("Args:", C.ADD_ARGS_TO_SCRIPT_THEN_ADD_TO_CHAIN)
		m.currentView = "input"
	}
	return m, nil
}

func removeScriptFromChain(m model) (tea.Model, tea.Cmd) {
	m.chain = backend.RemoveScriptFromChain(m.list.SelectedItem().(item).script, m.chain)
	return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
}
