package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type inputFinishedMsg bool
type inputRejectedMsg bool

type inputModel struct {
	textInput     textinput.Model
	err           error
	prompt        string
	Selected      bool
	returnCommand int
}

func initialInputModel(prompt string, returnCommand int) inputModel {
	ti := textinput.New()
	// ti.Placeholder = prompt
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return inputModel{
		textInput:     ti,
		err:           nil,
		prompt:        prompt,
		Selected:      false,
		returnCommand: returnCommand,
	}
}

func inputUpdate(m inputModel, msg tea.Msg) (inputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmd = func() tea.Msg {
				return tea.ClearScreen()
			}
			m.Selected = true

			return m, tea.Sequence(cmd, func() tea.Msg { return inputFinishedMsg(true) })

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Sequence(cmd, func() tea.Msg { return inputRejectedMsg(true) })
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd

}

func inputView(m inputModel) string {
	return fmt.Sprintf(
		m.prompt+"\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
