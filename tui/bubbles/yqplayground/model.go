package yqplayground

import (
	"os"
	"time"

	"github.com/arch-err/yqp/tui/bubbles/fileselector"
	"github.com/arch-err/yqp/tui/bubbles/help"
	"github.com/arch-err/yqp/tui/bubbles/inputdata"
	"github.com/arch-err/yqp/tui/bubbles/output"
	"github.com/arch-err/yqp/tui/bubbles/queryinput"
	"github.com/arch-err/yqp/tui/bubbles/state"
	"github.com/arch-err/yqp/tui/bubbles/statusbar"
	"github.com/arch-err/yqp/tui/theme"
)

func (b Bubble) GetState() state.State {
	return b.state
}

type Bubble struct {
	width            int
	height           int
	workingDirectory string
	state            state.State
	queryinput       queryinput.Bubble
	inputdata        inputdata.Bubble
	output           output.Bubble
	help             help.Bubble
	statusbar        statusbar.Bubble
	fileselector     fileselector.Bubble
	results          string
	cancel           func()
	theme            theme.Theme
	ExitMessage      string
	showInputPanel   bool
}

func New(inputYAML []byte, filename string, query string, yqtheme theme.Theme) (Bubble, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return Bubble{}, err
	}

	sb := statusbar.New(yqtheme)
	sb.StatusMessageLifetime = time.Second * 10
	fs := fileselector.New(yqtheme)

	fs.SetInput(workingDirectory)

	inputData, err := inputdata.New(inputYAML, filename, yqtheme)
	if err != nil {
		return Bubble{}, err
	}
	queryInput := queryinput.New(yqtheme)
	if query != "" {
		queryInput.SetQuery(query)
	}

	b := Bubble{
		workingDirectory: workingDirectory,
		state:            state.Loading,
		queryinput:       queryInput,
		inputdata:        inputData,
		output:           output.New(yqtheme),
		help:             help.New(yqtheme),
		statusbar:        sb,
		fileselector:     fs,
		theme:            yqtheme,
		showInputPanel:   true,
	}
	return b, nil
}
