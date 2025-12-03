package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/doganarif/k9sight/internal/k8s"
	"github.com/doganarif/k9sight/internal/ui/components"
	"github.com/doganarif/k9sight/internal/ui/keys"
	"github.com/doganarif/k9sight/internal/ui/styles"
)

type PanelFocus int

const (
	FocusLogs PanelFocus = iota
	FocusEvents
	FocusMetrics
	FocusManifest
)

type Dashboard struct {
	pod        *k8s.PodInfo
	logs       components.LogsPanel
	events     components.EventsPanel
	metrics    components.MetricsPanel
	manifest   components.ManifestPanel
	statusBar  components.StatusBar
	breadcrumb components.Breadcrumb
	help       components.HelpPanel
	focus      PanelFocus
	fullscreen bool
	width      int
	height     int
	keys       keys.KeyMap
}

func NewDashboard() Dashboard {
	return Dashboard{
		logs:       components.NewLogsPanel(),
		events:     components.NewEventsPanel(),
		metrics:    components.NewMetricsPanel(),
		manifest:   components.NewManifestPanel(),
		statusBar:  components.NewStatusBar(),
		breadcrumb: components.NewBreadcrumb(),
		help:       components.NewHelpPanel(),
		focus:      FocusLogs,
		keys:       keys.DefaultKeyMap(),
	}
}

func (d Dashboard) Init() tea.Cmd {
	return nil
}

func (d Dashboard) Update(msg tea.Msg) (Dashboard, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if d.help.IsVisible() {
			if msg.String() == "?" || msg.String() == "esc" {
				d.help.Hide()
				return d, nil
			}
			return d, nil
		}

		switch {
		case key.Matches(msg, d.keys.Help):
			d.help.Toggle()
			return d, nil

		case key.Matches(msg, d.keys.NextPanel):
			d.nextPanel()
			return d, nil

		case key.Matches(msg, d.keys.PrevPanel):
			d.prevPanel()
			return d, nil

		case key.Matches(msg, d.keys.Panel1):
			d.focus = FocusLogs
			return d, nil

		case key.Matches(msg, d.keys.Panel2):
			d.focus = FocusEvents
			return d, nil

		case key.Matches(msg, d.keys.Panel3):
			d.focus = FocusMetrics
			return d, nil

		case key.Matches(msg, d.keys.Panel4):
			d.focus = FocusManifest
			return d, nil

		case key.Matches(msg, d.keys.ToggleFullView):
			d.fullscreen = !d.fullscreen
			return d, nil
		}
	}

	switch d.focus {
	case FocusLogs:
		d.logs, cmd = d.logs.Update(msg)
		cmds = append(cmds, cmd)
	case FocusEvents:
		d.events, cmd = d.events.Update(msg)
		cmds = append(cmds, cmd)
	case FocusMetrics:
		d.metrics, cmd = d.metrics.Update(msg)
		cmds = append(cmds, cmd)
	case FocusManifest:
		d.manifest, cmd = d.manifest.Update(msg)
		cmds = append(cmds, cmd)
	}

	return d, tea.Batch(cmds...)
}

func (d *Dashboard) nextPanel() {
	d.focus = (d.focus + 1) % 4
}

func (d *Dashboard) prevPanel() {
	d.focus = (d.focus + 3) % 4
}

func (d Dashboard) View() string {
	if d.pod == nil {
		return styles.PanelStyle.Render("No pod selected")
	}

	var b strings.Builder

	b.WriteString(d.breadcrumb.View())
	b.WriteString("\n")

	if d.fullscreen {
		// Render only the focused panel in fullscreen
		b.WriteString(d.renderFullscreenPanel())
	} else {
		// Normal 4-panel layout
		topRow := d.renderTopRow()
		bottomRow := d.renderBottomRow()

		b.WriteString(topRow)
		b.WriteString("\n")
		b.WriteString(bottomRow)
	}

	content := b.String()

	if d.help.IsVisible() {
		return d.renderFloatingHelp(content)
	}

	return content
}

func (d Dashboard) renderFullscreenPanel() string {
	panelWidth := d.width - 4
	panelHeight := d.height - 8

	var content string
	switch d.focus {
	case FocusLogs:
		d.logs.SetSize(panelWidth, panelHeight)
		content = d.logs.View()
	case FocusEvents:
		d.events.SetSize(panelWidth, panelHeight)
		content = d.events.View()
	case FocusMetrics:
		d.metrics.SetSize(panelWidth, panelHeight)
		content = d.metrics.View()
	case FocusManifest:
		d.manifest.SetSize(panelWidth, panelHeight)
		content = d.manifest.View()
	}

	return d.wrapPanel(content, panelWidth, panelHeight, true)
}

func (d Dashboard) renderTopRow() string {
	halfWidth := d.width / 2
	panelHeight := (d.height - 6) / 2

	d.logs.SetSize(halfWidth-2, panelHeight)
	d.events.SetSize(halfWidth-2, panelHeight)

	logsView := d.wrapPanel(d.logs.View(), halfWidth-2, panelHeight, d.focus == FocusLogs)
	eventsView := d.wrapPanel(d.events.View(), halfWidth-2, panelHeight, d.focus == FocusEvents)

	return lipgloss.JoinHorizontal(lipgloss.Top, logsView, eventsView)
}

func (d Dashboard) renderBottomRow() string {
	halfWidth := d.width / 2
	panelHeight := (d.height - 6) / 2

	d.metrics.SetSize(halfWidth-2, panelHeight)
	d.manifest.SetSize(halfWidth-2, panelHeight)

	metricsView := d.wrapPanel(d.metrics.View(), halfWidth-2, panelHeight, d.focus == FocusMetrics)
	manifestView := d.wrapPanel(d.manifest.View(), halfWidth-2, panelHeight, d.focus == FocusManifest)

	return lipgloss.JoinHorizontal(lipgloss.Top, metricsView, manifestView)
}

func (d Dashboard) wrapPanel(content string, width, height int, active bool) string {
	style := styles.PanelStyle
	if active {
		style = styles.ActivePanelStyle
	}

	return style.
		Width(width).
		Height(height).
		Render(content)
}

func (d Dashboard) renderFloatingHelp(base string) string {
	// Get the help content
	helpContent := d.help.View()

	// Use lipgloss.Place to center the help modal over the content
	return lipgloss.Place(
		d.width,
		d.height-4,
		lipgloss.Center,
		lipgloss.Center,
		helpContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(styles.Background),
	)
}

func (d *Dashboard) SetPod(pod *k8s.PodInfo) {
	d.pod = pod
	d.manifest.SetPod(pod)
	d.metrics.SetPod(pod)
}

func (d *Dashboard) SetLogs(logs []k8s.LogLine) {
	d.logs.SetLogs(logs)
}

func (d *Dashboard) SetEvents(events []k8s.EventInfo) {
	d.events.SetEvents(events)
}

func (d *Dashboard) SetMetrics(metrics *k8s.PodMetrics) {
	d.metrics.SetMetrics(metrics)
}

func (d *Dashboard) SetRelated(related *k8s.RelatedResources) {
	d.manifest.SetRelated(related)
}

func (d *Dashboard) SetHelpers(helpers []k8s.DebugHelper) {
	d.manifest.SetHelpers(helpers)
}

func (d *Dashboard) SetSize(width, height int) {
	d.width = width
	d.height = height
	d.statusBar.SetWidth(width)
	d.breadcrumb.SetWidth(width)
	d.help.SetSize(width, height)
}

func (d *Dashboard) SetBreadcrumb(items ...string) {
	d.breadcrumb.SetItems(items...)
}

func (d *Dashboard) SetContext(ctx string) {
	d.statusBar.SetContext(ctx)
}

func (d *Dashboard) SetNamespace(ns string) {
	d.statusBar.SetNamespace(ns)
}

func (d *Dashboard) SetResource(res string) {
	d.statusBar.SetResource(res)
}

func (d Dashboard) Focus() PanelFocus {
	return d.focus
}

func (d Dashboard) HelpVisible() bool {
	return d.help.IsVisible()
}

func (d Dashboard) ShortHelp() string {
	return d.help.ShortHelp()
}
