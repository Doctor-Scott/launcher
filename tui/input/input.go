package tui_input

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type InputFinishedMsg bool
type InputRejectedMsg bool

type InputModel struct {
	TextInput textinput.Model
	err       error
	prompt    string
	Selected  bool
	InputType string
}

func InitialInputModel(prompt string, inputType string) InputModel {
	ti := textinput.New()
	// ti.Placeholder = prompt
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return InputModel{
		TextInput: ti,
		err:       nil,
		prompt:    prompt,
		Selected:  false,
		InputType: inputType,
	}
}

func InputUpdate(m InputModel, msg tea.Msg) (InputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmd = func() tea.Msg {
				return tea.ClearScreen()
			}
			m.Selected = true

			return m, tea.Sequence(cmd, func() tea.Msg { return InputFinishedMsg(true) })

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Sequence(cmd, func() tea.Msg { return InputRejectedMsg(true) })
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd

}

func InputView(m InputModel) string {
	return fmt.Sprintf(
		m.prompt+"\n\n%s\n\n%s",
		m.TextInput.View(),
		"(esc to quit)",
	) + "\n"
}
