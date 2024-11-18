package tui

import (
	backend "github.com/Doctor-Scott/launcher/backend"
	C "github.com/Doctor-Scott/launcher/globalConstants"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

func getCustomDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		// make each selected item a different color
		for i, listItem := range m.Items() {
			item := listItem.(item)

			if item.selected == true {
				item.titlePretty = lipgloss.NewStyle().Foreground(lipgloss.Color(viper.GetString(C.ColorConfig.SelectedScript.Name))).Render(item.title)
				m.SetItem(i, item)
			} else if item.title == "Input" {
				item.titlePretty = lipgloss.NewStyle().Foreground(lipgloss.Color(viper.GetString(C.ColorConfig.InputTitle.Name))).Render(item.title)
				m.SetItem(i, item)

			} else {
				item.titlePretty = lipgloss.NewStyle().Foreground(lipgloss.NoColor{}).Render(item.title)
				m.SetItem(i, item)

			}
		}

		return nil
	}
	c := lipgloss.Color(viper.GetString(C.ColorConfig.Cursor.Name))
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(c).BorderLeftForeground(c)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle // reuse the title style here
	return delegate
}

func createScriptList(path string) list.Model {
	if path == "" {
		path = viper.GetString(C.PathConfig.ScriptDir.Name)
	}
	structure := backend.GetStructure(path)
	items := []list.Item{}

	items = append(items, item{title: "Input", desc: C.INPUT_COMMAND_DESC, script: backend.Script{Name: C.INPUT_COMMAND_NAME}})
	for _, script := range structure {
		items = append(items, item{title: script.Name, script: script})
	}
	delegate := getCustomDelegate()
	list := list.New(items, delegate, 0, 0)
	list.Title = "Scripts"
	list.Styles.Title = getTitleStyle("script")
	return list
}

func getTitleStyle(view string) lipgloss.Style {
	style := lipgloss.NewStyle().
		Bold(true).
		Width(15).
		Align(lipgloss.Center)

	if view == "chain" {
		return style.Background(lipgloss.Color(viper.GetString(C.ColorConfig.ChainTitle.Name)))
	}
	return style.Background(lipgloss.Color(viper.GetString(C.ColorConfig.ScriptTitle.Name)))

}

func createChainList(path string) list.Model {
	if path == "" {
		path = viper.GetString(C.PathConfig.LauncherDir.Name) + "/custom/"
	}
	structure := backend.GetChainStructure(path)
	items := []list.Item{}

	for _, chainItem := range structure {
		items = append(items, item{title: chainItem.Name, chainItem: chainItem})
	}
	delegate := getCustomDelegate()
	list := list.New(items, delegate, 0, 0)

	list.Title = "Workflows"
	list.Styles.Title = getTitleStyle("chain")

	return list
}

func createNewScriptModelList(m model) model {
	m.list = createScriptList(m.currentPath)
	return m
}

func createNewChainModelList(m model) model {
	m.list = createChainList("")
	return m
}
