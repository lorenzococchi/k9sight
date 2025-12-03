package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/doganarif/k9sight/internal/k8s"
	"github.com/doganarif/k9sight/internal/ui/styles"
)

type LogsPanel struct {
	logs      []k8s.LogLine
	viewport  viewport.Model
	ready     bool
	width     int
	height    int
	following bool
	filter    string
	container string
}

func NewLogsPanel() LogsPanel {
	return LogsPanel{
		following: true,
	}
}

func (l LogsPanel) Init() tea.Cmd {
	return nil
}

func (l LogsPanel) Update(msg tea.Msg) (LogsPanel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "F":
			l.following = !l.following
			if l.following {
				l.viewport.GotoBottom()
			}
		case "e":
			l.jumpToNextError()
		case "g":
			l.viewport.GotoTop()
		case "G":
			l.viewport.GotoBottom()
		}
	}

	l.viewport, cmd = l.viewport.Update(msg)
	return l, cmd
}

func (l LogsPanel) View() string {
	if !l.ready {
		return styles.PanelStyle.Render("Loading logs...")
	}

	var header strings.Builder
	header.WriteString(styles.PanelTitleStyle.Render("Logs"))
	if l.container != "" {
		header.WriteString(styles.SubtitleStyle.Render(fmt.Sprintf(" [%s]", l.container)))
	}
	if l.following {
		header.WriteString(styles.StatusRunning.Render(" [Following]"))
	}
	header.WriteString("\n")

	return header.String() + l.viewport.View()
}

func (l *LogsPanel) SetLogs(logs []k8s.LogLine) {
	l.logs = logs
	l.updateContent()
}

func (l *LogsPanel) SetSize(width, height int) {
	l.width = width
	l.height = height - 2

	if !l.ready {
		l.viewport = viewport.New(width, l.height)
		l.ready = true
	} else {
		l.viewport.Width = width
		l.viewport.Height = l.height
	}

	l.updateContent()
}

func (l *LogsPanel) SetContainer(container string) {
	l.container = container
}

func (l *LogsPanel) SetFilter(filter string) {
	l.filter = filter
	l.updateContent()
}

func (l *LogsPanel) ToggleFollow() {
	l.following = !l.following
	if l.following {
		l.viewport.GotoBottom()
	}
}

func (l *LogsPanel) updateContent() {
	if !l.ready {
		return
	}

	var content strings.Builder
	filteredLogs := l.getFilteredLogs()

	for _, log := range filteredLogs {
		line := l.formatLogLine(log)
		content.WriteString(line)
		content.WriteString("\n")
	}

	l.viewport.SetContent(content.String())

	if l.following {
		l.viewport.GotoBottom()
	}
}

func (l LogsPanel) getFilteredLogs() []k8s.LogLine {
	if l.filter == "" {
		return l.logs
	}

	filter := strings.ToLower(l.filter)
	var filtered []k8s.LogLine
	for _, log := range l.logs {
		if strings.Contains(strings.ToLower(log.Content), filter) {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

func (l LogsPanel) formatLogLine(log k8s.LogLine) string {
	var b strings.Builder

	if !log.Timestamp.IsZero() {
		ts := log.Timestamp.Format("15:04:05")
		b.WriteString(styles.LogTimestamp.Render(ts))
		b.WriteString(" ")
	}

	if log.Container != "" && l.container == "" {
		b.WriteString(styles.LogContainer.Render(fmt.Sprintf("[%s]", log.Container)))
		b.WriteString(" ")
	}

	if log.IsError {
		b.WriteString(styles.LogError.Render(log.Content))
	} else {
		b.WriteString(styles.LogNormal.Render(log.Content))
	}

	return b.String()
}

func (l *LogsPanel) jumpToNextError() {
	content := l.viewport.View()
	lines := strings.Split(content, "\n")
	currentLine := l.viewport.YOffset

	for i := currentLine + 1; i < len(lines); i++ {
		if strings.Contains(strings.ToLower(lines[i]), "error") ||
			strings.Contains(strings.ToLower(lines[i]), "fatal") ||
			strings.Contains(strings.ToLower(lines[i]), "panic") {
			l.viewport.SetYOffset(i)
			return
		}
	}

	for i := 0; i < currentLine; i++ {
		if strings.Contains(strings.ToLower(lines[i]), "error") ||
			strings.Contains(strings.ToLower(lines[i]), "fatal") ||
			strings.Contains(strings.ToLower(lines[i]), "panic") {
			l.viewport.SetYOffset(i)
			return
		}
	}
}

func (l LogsPanel) IsFollowing() bool {
	return l.following
}

func (l LogsPanel) LogCount() int {
	return len(l.logs)
}

func (l LogsPanel) ErrorCount() int {
	count := 0
	for _, log := range l.logs {
		if log.IsError {
			count++
		}
	}
	return count
}
