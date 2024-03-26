package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error
type model struct {
	textInput textinput.Model
	err       error
	prompt    string
}

func (m model) GetTextInput() string {
	return m.textInput.Value()
}

func initialModel(prompt string) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 50

	return model{
		textInput: ti,
		err:       nil,
		prompt:    prompt,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
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
		"(ctrl-c or esc to quit)",
	) + "\n"
}
