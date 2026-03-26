package yqplayground

import (
	"os"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikefarah/yq/v4/pkg/yqlib"

	"github.com/arch-err/yqp/tui/utils"
)

type errorMsg struct {
	error error
}

type queryResultMsg struct {
	rawResults         string
	highlightedResults string
}

type writeToFileMsg struct{}

type copyQueryToClipboardMsg struct{}

type copyResultsToClipboardMsg struct{}

func (b *Bubble) executeQueryOnInput() (string, error) {
	expression := b.queryinput.GetInputValue()
	if expression == "" {
		expression = "."
	}

	yamlInput := string(b.inputdata.GetInputData())

	encoder := yqlib.NewYamlEncoder(yqlib.NewDefaultYamlPreferences())
	decoder := yqlib.NewYamlDecoder(yqlib.NewDefaultYamlPreferences())

	eval := yqlib.NewStringEvaluator()
	result, err := eval.EvaluateAll(expression, yamlInput, encoder, decoder)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(result, "\n"), nil
}

func (b *Bubble) executeQueryCommand() tea.Cmd {
	return func() tea.Msg {
		results, err := b.executeQueryOnInput()
		if err != nil {
			return errorMsg{error: err}
		}
		highlightedOutput, err := utils.Prettify([]byte(results), b.theme.ChromaStyle)
		if err != nil {
			return errorMsg{error: err}
		}
		return queryResultMsg{
			rawResults:         results,
			highlightedResults: highlightedOutput.String(),
		}
	}
}

func (b Bubble) saveOutput() tea.Cmd {
	if b.fileselector.GetInput() == "" {
		return b.copyOutputToClipboard()
	}
	return b.writeOutputToFile()
}

func (b Bubble) copyOutputToClipboard() tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(b.results)
		if err != nil {
			return errorMsg{
				error: err,
			}
		}
		return copyResultsToClipboardMsg{}
	}
}

func (b Bubble) writeOutputToFile() tea.Cmd {
	return func() tea.Msg {
		err := os.WriteFile(b.fileselector.GetInput(), []byte(b.results), 0o600)
		if err != nil {
			return errorMsg{
				error: err,
			}
		}
		return writeToFileMsg{}
	}
}

func (b Bubble) copyQueryToClipboard() tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(b.queryinput.GetInputValue())
		if err != nil {
			return errorMsg{
				error: err,
			}
		}
		return copyQueryToClipboardMsg{}
	}
}
