package statusbar

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/arch-err/yqp/tui/theme"
)

type Bubble struct {
	styles                styles
	StatusMessageLifetime time.Duration
	statusMessage         string
	statusMessageTimer    *time.Timer
}

func (Bubble) Init() tea.Cmd {
	return nil
}

func (b Bubble) View() string {
	return b.styles.containerStyle.Render(b.statusMessage)
}

func (b *Bubble) SetSize(width int) {
	b.styles.containerStyle = b.styles.containerStyle.Width(width)
}

func (b *Bubble) hideStatusMessage() {
	b.statusMessage = ""
	if b.statusMessageTimer != nil {
		b.statusMessageTimer.Stop()
	}
}

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case statusMessageTimeoutMsg:
		b.hideStatusMessage()
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width)
	default:
	}
	return b, tea.Batch(cmd)
}

func New(yqtheme theme.Theme) Bubble {
	styles := defaultStyles()
	styles.successMessageStyle = styles.successMessageStyle.Foreground(yqtheme.Success)
	styles.errorMessageStyle = styles.errorMessageStyle.Foreground(yqtheme.Error)
	b := Bubble{
		styles: styles,
	}
	return b
}

type statusMessageTimeoutMsg struct{}

func (b *Bubble) NewStatusMessage(s string, success bool) tea.Cmd {
	if success {
		b.statusMessage = b.styles.successMessageStyle.Render(s)
	} else {
		b.statusMessage = b.styles.errorMessageStyle.Render(s)
	}

	if b.statusMessageTimer != nil {
		b.statusMessageTimer.Stop()
	}

	b.statusMessageTimer = time.NewTimer(b.StatusMessageLifetime)

	// Wait for timeout
	return func() tea.Msg {
		<-b.statusMessageTimer.C
		return statusMessageTimeoutMsg{}
	}
}
