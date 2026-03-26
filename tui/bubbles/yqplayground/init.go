package yqplayground

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/arch-err/yqp/tui/utils"
)

// invalidInputMsg signals that the user's data is not valid YAML
type invalidInputMsg struct{}

type setupMsg struct{}

// initialQueryMsg represents a message containing an initial query string to execute when
// the app is loaded.
type initialQueryMsg struct {
	query string
}

func setupCmd() tea.Cmd {
	return func() tea.Msg {
		return setupMsg{}
	}
}

// initialQueryCmd creates a command that returns an initialQueryMsg with the provided query string.
func initialQueryCmd(query string) tea.Cmd {
	return func() tea.Msg {
		return initialQueryMsg{query: query}
	}
}

func (b Bubble) Init() tea.Cmd {
	var cmds []tea.Cmd

	// validate input data
	_, err := utils.IsValidInput(b.inputdata.GetInputData())
	if err != nil {
		return func() tea.Msg {
			return invalidInputMsg{}
		}
	}

	// initialize rest of app
	cmds = append(cmds, b.queryinput.Init(), b.inputdata.Init(), setupCmd())
	if b.queryinput.GetInputValue() != "" {
		cmds = append(cmds, initialQueryCmd(b.queryinput.GetInputValue()))
	}
	return tea.Sequence(cmds...)
}
