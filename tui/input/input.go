package tui_input

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func Input(prompt string) string {
	p := tea.NewProgram(initialModel(prompt))
	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	finalModel := m.(model)
	return finalModel.textInput.Value()
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
	prompt    string
}

func initialModel(prompt string) model {
	ti := textinput.New()
	// ti.Placeholder = prompt
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
		prompt:    prompt,
	}
}

func (m model) Init() tea.Cmd {
	cmd := func() tea.Msg {
		return tea.ClearScreen()
	}
	return tea.Sequence(cmd, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			cmd = func() tea.Msg {
				return tea.ClearScreen()
			}
			// quit := func() tea.Msg {
			// 	return tea.Quit
			// }

			return m, tea.Sequence(cmd, tea.Quit)
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		m.prompt+"\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
