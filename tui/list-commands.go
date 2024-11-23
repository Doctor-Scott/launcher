package tui

import (
	"bytes"
	"fmt"
	backend "github.com/Doctor-Scott/launcher/backend"
	C "github.com/Doctor-Scott/launcher/globalConstants"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

func debugModel(m model) {
	for _, item := range m.lists.scripts.Items() {
		fmt.Printf("%+v", item)
		fmt.Printf("\n")

	}

}
func debug(m model) (tea.Model, tea.Cmd) {
	// fmt.Println(m.currentPath)
	// fmt.Println("")
	// fmt.Println(m.list.Items())
	// fmt.Println("")
	debugModel(m)

	// fmt.Println(m.chain)
	// fmt.Println(string(m.stdout))
	// fmt.Print(m)

	// fmt.Println(m.Items())

	// fmt.Print(viper.AllSettings())
	return m, nil
}

func inputLauncherCommand(m model) (tea.Model, tea.Cmd) {
	m.newView(InputView)
	m.inputModel = createNewInputModel("> ", C.RUN_LAUNCHER_COMMAND)
	return m, nil

}

func workflowView(m model) (tea.Model, tea.Cmd) {
	// swap view between scripts and workflows
	m.newView(WorkflowsView)
	return m, func() tea.Msg { return updateStructureMsg(true) }
}

func scriptView(m model) (tea.Model, tea.Cmd) {
	// swap view between workflows and scripts
	m.newView(ScriptsView)
	return m, func() tea.Msg { return updateStructureMsg(true) }
}

func addScriptToChain(m model, itemType string) (tea.Model, tea.Cmd) {
	if m.lists.scripts.SelectedItem().(item).title == "Input" && m.views.currentView == ScriptsView {
		m.inputModel = createNewInputModel("Script:", C.ADD_SCRIPT_TO_CHAIN)
		m.newView(InputView)
		return m, nil
	} else {
		if itemType == "chain" {
			m.chain = backend.AddChainToChain(m.lists.workflows.SelectedItem().(item).chainItem.Chain, m.chain)
			return scriptView(m)
		}
		script := m.lists.scripts.SelectedItem().(item).script
		m.chain = backend.AddScriptToChain(script, m.chain)
		return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
	}
}

func refreshView(m model) (tea.Model, tea.Cmd) {
	m.lists.scripts.ResetSelected()
	m.lists.workflows.ResetSelected()
	return m, func() tea.Msg { return tea.ClearScreen() }
}

func runChain(m model) (tea.Model, tea.Cmd) {
	chainResult := backend.RunChain(m.stdout, m.chain)
	lastScriptResult := chainResult[len(chainResult)-1]
	//TODO  here we should check if the last script failed

	m.stdout = lastScriptResult.Stdout
	m = maybeSetLastFailedScript(m, lastScriptResult)
	//TODO  maybeSetLastFailedChain(m, lastScriptResult)

	if !lastScriptResult.Success {
		fmt.Println(lastScriptResult.Stderr)
		return m, nil
	}

	if viper.GetBool(C.ClearChainAfterRun.Name) {
		m.chain = backend.Chain{}
	}

	return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
}

func editItemUnderCursor(m model, itemType string) (tea.Model, tea.Cmd) {
	//WARN  Now that we hold bot states, will the selected item still hold input
	// yes, this is the case
	if m.lists.scripts.SelectedItem().(item).title != "Input" || m.views.currentView != ScriptsView {
		var pathToChain string
		if itemType == "chain" {
			pathToChain = viper.GetString(C.PathConfig.LauncherDir.Name) + "/custom/" + m.lists.workflows.SelectedItem().(item).chainItem.Name + ".json"
		} else {
			pathToChain = m.lists.scripts.SelectedItem().(item).script.Command
		}

		editor := os.ExpandEnv("$EDITOR")
		cmd := exec.Command(editor, pathToChain)
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return fmt.Errorf("failed to run : %w", err)
			}
			return updateStructureMsg(true)
		})

	}
	return m, nil

}

func deleteChainUnderCursor(m model) (tea.Model, tea.Cmd) {
	backend.DeleteChainConfig(m.lists.workflows.SelectedItem().(item).chainItem.Name)
	return m, func() tea.Msg { return updateStructureMsg(true) }
}

func openEditorInLauncherDirectory(m model) (tea.Model, tea.Cmd) {
	editor := os.ExpandEnv("$EDITOR")
	//WARN  Not sure if this works for all editors
	cmd := exec.Command(editor, "--cmd", "cd"+m.currentPath+" | enew")
	m.lists.scripts.ResetSelected()
	m.lists.workflows.ResetSelected()

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return fmt.Errorf("failed to run : %w", err)
		}
		return updateStructureMsg(true)
	})
}

func openConfig(m model) (tea.Model, tea.Cmd) {
	configFile := viper.ConfigFileUsed()
	editor := os.ExpandEnv("$EDITOR")
	cmd := exec.Command(editor, configFile)
	// m.list.ResetSelected()
	// fmt.Println(configFile)
	// return m, nil

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return fmt.Errorf("failed to run : %w", err)
		}
		return updateStructureMsg(true)
	})
}

func writeConfig(m model) (tea.Model, tea.Cmd) {
	viper.WriteConfig()
	return m, nil
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
	m.inputModel = createNewInputModel("Name:", C.LOAD_CUSTOM_CHAIN)
	m.newView(InputView)
	return m, nil
}

func writeChain(m model) (tea.Model, tea.Cmd) {
	m.inputModel = createNewInputModel("Name:", C.SAVE_CUSTOM_CHAIN)
	m.newView(InputView)
	return m, nil
}

func clearState(m model) (tea.Model, tea.Cmd) {
	m.chain = backend.Chain{}
	m.stdout = []byte{}
	m = createNewScriptModelList(createNewChainModelList(m))
	return m, func() tea.Msg { return updateStructureMsg(true) }
}

func runItemUnderCursor(m model, itemType string) (tea.Model, tea.Cmd) {
	if itemType == "chain" {
		chainResult := backend.RunChain(m.stdout, m.lists.workflows.SelectedItem().(item).chainItem.Chain)
		lastScriptResult := chainResult[len(chainResult)-1]

		m.stdout = lastScriptResult.Stdout
		if !lastScriptResult.Success {
			fmt.Println(lastScriptResult.Stderr)
			return m, nil
		}
		m = maybeSetLastFailedScript(m, lastScriptResult)

		cmd := func() tea.Msg {
			return tea.ClearScreen()
		}
		return m, cmd
	}

	// standard run of known script or input command
	if m.lists.scripts.SelectedItem().(item).title == "Input" {
		m.inputModel = createNewInputModel("Script:", C.RUN_SCRIPT)
		m.newView(InputView)
		return m, nil
	} else {
		scriptResult := backend.RunScript(m.lists.scripts.SelectedItem().(item).script, m.stdout)
		m.stdout = scriptResult.Stdout
		m = maybeSetLastFailedScript(m, scriptResult)
		cmd := func() tea.Msg {
			return tea.ClearScreen()
		}
		return m, cmd
	}

}

func runScriptWithArgs(m model) (tea.Model, tea.Cmd) {
	if m.lists.scripts.SelectedItem().(item).title != "Input" {
		m.inputModel = createNewInputModel("Args:", C.ADD_ARGS_TO_SCRIPT_AND_RUN)
		m.newView(InputView)
		cmd := func() tea.Msg {
			return tea.ClearScreen()
		}
		return m, cmd
	}
	return m, nil
}

func addScriptWithArgs(m model) (tea.Model, tea.Cmd) {
	if m.lists.scripts.SelectedItem().(item).title != "Input" {
		m.inputModel = createNewInputModel("Args:", C.ADD_ARGS_TO_SCRIPT_THEN_ADD_TO_CHAIN)
		m.newView(InputView)
	}
	return m, nil
}

func removeScriptFromChain(m model) (tea.Model, tea.Cmd) {
	m.chain = backend.RemoveScriptFromChain(m.lists.scripts.SelectedItem().(item).script, m.chain)
	return m, func() tea.Msg { return generateSelectedItemViewMsg(true) }
}
