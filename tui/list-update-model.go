package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	backend "launcher/backend"
	C "launcher/globalConstants"
)

func getCustomDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		// make each selected item a different color
		for i, listItem := range m.Items() {
			item := listItem.(item)

			if item.script.Selected == true {
				// TODO  Add this to the config
				// or maybe from the terminal colour scheme?
				item.titlePretty = lipgloss.NewStyle().Foreground(lipgloss.Color("#6fe600")).Render(item.title)
				m.SetItem(i, item)
			} else if item.title == "Input" {
				item.titlePretty = lipgloss.NewStyle().Foreground(lipgloss.Color("#e64d00")).Render(item.title)
				m.SetItem(i, item)

			} else {
				item.titlePretty = lipgloss.NewStyle().Foreground(lipgloss.NoColor{}).Render(item.title)
				m.SetItem(i, item)

			}
		}

		return nil
	}
	// TODO  Add this to the config
	// or maybe from the terminal colour scheme?
	c := lipgloss.Color("#6fe6fc")
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(c).BorderLeftForeground(c)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle // reuse the title style here
	return delegate
}

func createList(path string) list.Model {
	structure := backend.GetStructure(path)
	items := []list.Item{}

	items = append(items, item{title: "Input", desc: C.INPUT_SCRIPT_DESC, script: backend.Script{Name: C.INPUT_SCRIPT_NAME}})
	for _, script := range structure {
		items = append(items, item{title: script.Name, script: script})
	}
	delegate := getCustomDelegate()
	list := list.New(items, delegate, 0, 0)
	list.Title = "Running a script are we???"
	return list
}

func createNewModelList(m model) model {
	m.list = createList(m.currentPath)
	// backend.SaveChain(m.chain)
	return m
}
