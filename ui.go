package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	screenList = iota
	screenAction
	screenView
	screenToast
	screenCustomInput
	screenRunning
)

var actions = []string{"Show devlog", "Show code", "Open judge", "Run sample inputs", "Run custom input"}

type model struct {
	problems []Problem
	cursor   int
	screen   int

	selected *Problem
	action   int
	output   string
	customIn string

	width, height int
	listOffset    int
	viewLine      int

	toastText  string
	toastIsErr bool
	toastID    int

	runningJobID int
	prevScreen   int

	viewIsSample      bool
	showSampleDetails bool
	lastSampleResults []TestResult
}

type toastTimeoutMsg struct{ id int }
type runDoneMsg struct {
	id            int
	output        string
	sampleResults []TestResult
}

func NewModel(problems []Problem) model { return model{problems: problems, screen: screenList} }
func (m model) Init() tea.Cmd           { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.adjustListOffset()
		m.clampViewLine()
		return m, nil
	case toastTimeoutMsg:
		if m.screen == screenToast && msg.id == m.toastID {
			m.screen = m.prevScreen
		}
		return m, nil
	case runDoneMsg:
		if msg.id == m.runningJobID && m.screen == screenRunning {
			if msg.sampleResults != nil {
				m.lastSampleResults = msg.sampleResults
				m.viewIsSample = true
				m.output = formatTestResults(msg.sampleResults, m.showSampleDetails)
			} else {
				m.viewIsSample = false
				m.output = msg.output
			}
			m.viewLine, m.screen = 0, screenView
		}
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	switch key {
	case "q", "esc":
		return m.handleBack(key == "q")
	case "up", "k":
		if m.screen == screenList && m.cursor > 0 {
			m.cursor--
			m.adjustListOffset()
		} else if m.screen == screenAction && m.action > 0 {
			m.action--
		} else if m.screen == screenView {
			m.scrollView(-1)
		}
	case "down", "j":
		if m.screen == screenList && m.cursor < len(m.problems)-1 {
			m.cursor++
			m.adjustListOffset()
		} else if m.screen == screenAction && m.action < len(actions)-1 {
			m.action++
		} else if m.screen == screenView {
			m.scrollView(1)
		}
	case "pgup", "ctrl+u":
		if m.screen == screenView {
			m.scrollView(-m.viewRows())
		}
	case "pgdown", "ctrl+d":
		if m.screen == screenView {
			m.scrollView(m.viewRows())
		}
	case "ctrl+o":
		if m.screen == screenView && m.viewIsSample {
			m.showSampleDetails = !m.showSampleDetails
			m.output = formatTestResults(m.lastSampleResults, m.showSampleDetails)
			m.clampViewLine()
		}
	case "enter":
		if m.screen == screenCustomInput {
			m.customIn += "\n"
			return m, nil
		}
		if m.screen == screenList {
			if len(m.problems) > 0 {
				m.selected, m.action, m.screen = &m.problems[m.cursor], 0, screenAction
			}
			return m, nil
		}
		if m.screen == screenAction {
			return m.selectAction()
		}
	case "ctrl+r":
		if m.screen == screenCustomInput {
			return m.startRunCustom(m.customIn)
		}
	case "backspace":
		if m.screen == screenCustomInput && len(m.customIn) > 0 {
			r := []rune(m.customIn)
			m.customIn = string(r[:len(r)-1])
		}
	default:
		if m.screen == screenCustomInput && len(msg.Runes) > 0 {
			m.customIn += string(msg.Runes)
		}
	}
	return m, nil
}

func (m model) handleBack(qPressed bool) (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenList:
		if qPressed {
			return m, tea.Quit
		}
	case screenAction:
		m.screen = screenList
	case screenView, screenCustomInput:
		m.screen = screenAction
	case screenToast, screenRunning:
		m.screen = m.prevScreen
	}
	return m, nil
}

func (m model) selectAction() (tea.Model, tea.Cmd) {
	switch m.action {
	case 0:
		text, err := readText(m.selected.DevlogPath)
		if err != nil {
			return m, m.showToast("Could not read devlog: "+err.Error(), true, screenAction)
		}
		m.output, m.viewLine, m.screen = text, 0, screenView
		m.viewIsSample = false
	case 1:
		text, err := readText(m.selected.CodePath)
		if err != nil {
			return m, m.showToast("Could not read code: "+err.Error(), true, screenAction)
		}
		m.output, m.viewLine, m.screen = text, 0, screenView
		m.viewIsSample = false
	case 2:
		if m.selected.JudgeURL == "" {
			return m, m.showToast("Judge URL not found", true, screenAction)
		}
		if err := OpenInSystem(m.selected.JudgeURL); err != nil {
			return m, m.showToast("Could not open judge: "+err.Error(), true, screenAction)
		}
		return m, m.showToast("Opened judge link", false, screenAction)
	case 3:
		return m.startRunSamples()
	case 4:
		m.customIn, m.screen = "", screenCustomInput
	}
	return m, nil
}

func (m model) View() string {
	switch m.screen {
	case screenList:
		return m.renderList()
	case screenAction:
		return m.renderAction()
	case screenView:
		return m.renderView()
	case screenToast:
		msg := m.toastText
		if m.toastIsErr {
			msg = errorStyle.Render(msg)
		} else {
			msg = infoStyle.Render(msg)
		}
		return msg + "\n\n" + footerStyle.Render("esc/q: dismiss")
	case screenCustomInput:
		return m.renderCustomInput()
	case screenRunning:
		return titleStyle.Render(m.selected.Name) + "\n\n" + infoStyle.Render("Running...") + "\n\n" + footerStyle.Render("esc/q: back")
	default:
		return ""
	}
}

func (m model) renderList() string {
	if len(m.problems) == 0 {
		return titleStyle.Render("Competitive Programming Problems") + "\n\n" + errorStyle.Render("No problems found") + "\n\n" + footerStyle.Render("q: quit")
	}
	visible, start := m.listRows(), m.listOffset
	end := min(start+visible, len(m.problems))
	lines := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		line := "  " + m.problems[i].Name
		if i == m.cursor {
			line = cursorStyle.Render("> " + m.problems[i].Name)
		}
		lines = append(lines, line)
	}
	return titleStyle.Render("Competitive Programming Problems") + "\n\n" + strings.Join(lines, "\n") + "\n\n" + footerStyle.Render("enter: select | q: quit")
}

func (m model) renderAction() string {
	lines := make([]string, 0, len(actions))
	for i, action := range actions {
		line := "  " + action
		if i == m.action {
			line = cursorStyle.Render("> " + action)
		}
		lines = append(lines, line)
	}
	return titleStyle.Render(m.selected.Name) + "\n\n" + strings.Join(lines, "\n") + "\n\n" + footerStyle.Render("enter: select | esc/q: back")
}

func (m model) renderView() string {
	lines, visible := m.outputLines(), m.viewRows()
	start := clamp(m.viewLine, 0, max(0, len(lines)-visible))
	end := min(start+visible, len(lines))
	footer := "j/k: scroll | esc/q: back"
	if m.viewIsSample {
		footer = "j/k: scroll | ctrl+o: toggle io | esc/q: back"
	}
	return titleStyle.Render(m.selected.Name) + "\n\n" + strings.Join(lines[start:end], "\n") + "\n\n" + footerStyle.Render(footer)
}

func (m model) renderCustomInput() string {
	width := max(1, m.width-2)
	lines := []string{}
	for _, raw := range strings.Split(strings.ReplaceAll(m.customIn+"█", "\r\n", "\n"), "\n") {
		lines = append(lines, wrapLine(raw, width)...)
	}
	if len(lines) == 0 {
		lines = []string{"█"}
	}
	start := max(0, len(lines)-m.viewRows())
	return titleStyle.Render("Custom Input") + "\n\n" + strings.Join(lines[start:], "\n") + "\n\n" + footerStyle.Render("enter: newline | ctrl+r: run | esc/q: back")
}

func (m *model) startRunSamples() (tea.Model, tea.Cmd) {
	m.runningJobID++
	id := m.runningJobID
	m.prevScreen, m.screen = screenAction, screenRunning
	selected := *m.selected
	return m, func() tea.Msg { return runDoneMsg{id: id, sampleResults: RunTestCases(selected)} }
}

func (m *model) startRunCustom(input string) (tea.Model, tea.Cmd) {
	m.runningJobID++
	id := m.runningJobID
	m.prevScreen, m.screen = screenCustomInput, screenRunning
	selected := *m.selected
	return m, func() tea.Msg {
		out, err := RunCustomInput(selected, input)
		return runDoneMsg{id: id, output: formatCustomRunOutput(input, out, err)}
	}
}

func (m *model) showToast(text string, isErr bool, back int) tea.Cmd {
	m.toastText, m.toastIsErr, m.prevScreen = text, isErr, back
	m.toastID++
	id := m.toastID
	m.screen = screenToast
	return tea.Tick(2500*time.Millisecond, func(time.Time) tea.Msg { return toastTimeoutMsg{id: id} })
}

func (m *model) adjustListOffset() {
	visible := m.listRows()
	if m.cursor < m.listOffset {
		m.listOffset = m.cursor
	}
	if m.cursor >= m.listOffset+visible {
		m.listOffset = m.cursor - visible + 1
	}
	if m.listOffset < 0 {
		m.listOffset = 0
	}
}

func (m *model) scrollView(delta int) { m.viewLine += delta; m.clampViewLine() }

func (m *model) clampViewLine() {
	maxStart := max(0, len(m.outputLines())-m.viewRows())
	m.viewLine = clamp(m.viewLine, 0, maxStart)
}

func (m model) listRows() int { return m.visibleRows() }
func (m model) viewRows() int { return m.visibleRows() }

func (m model) visibleRows() int {
	if m.height <= 0 {
		return 18
	}
	return max(1, m.height-6)
}

func (m model) outputLines() []string {
	if strings.TrimSpace(m.output) == "" {
		return []string{"(empty output)"}
	}
	width := max(1, m.width-2)
	text := strings.ReplaceAll(strings.ReplaceAll(m.output, "\r\n", "\n"), "\r", "\n")
	raw, lines := strings.Split(text, "\n"), make([]string, 0)
	for _, line := range raw {
		lines = append(lines, wrapLine(line, width)...)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func wrapLine(s string, width int) []string {
	if width < 1 {
		width = 1
	}
	r := []rune(s)
	if len(r) == 0 {
		return []string{""}
	}
	out := make([]string, 0, (len(r)/width)+1)
	for len(r) > width {
		out, r = append(out, string(r[:width])), r[width:]
	}
	return append(out, string(r))
}

func readText(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("file path is empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	text := string(data)
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("file is empty")
	}
	return text, nil
}

func formatTestResults(results []TestResult, showDetails bool) string {
	var b strings.Builder
	passed := 0
	for _, r := range results {
		if r.Passed {
			passed++
			b.WriteString(passStyle.Render(fmt.Sprintf("✓ %s", r.InputFile)) + "\n")
			if showDetails {
				appendLabeledBlock(&b, "input", r.Input)
				appendLabeledBlock(&b, "got", r.Got)
			}
			continue
		}
		b.WriteString(errorStyle.Render(fmt.Sprintf("✗ %s", r.InputFile)) + "\n")
		if showDetails {
			appendLabeledBlock(&b, "input", r.Input)
			if r.Expected != "" {
				appendLabeledBlock(&b, "expected", r.Expected)
			}
			gotText := r.Got
			if r.InputFile == "compile" && strings.Contains(gotText, bitsCompilerHint) {
				gotText = strings.TrimSpace(strings.TrimSuffix(gotText, bitsCompilerHint))
			}
			appendLabeledBlock(&b, "got", gotText)
			if r.InputFile == "compile" && strings.Contains(r.Got, bitsCompilerHint) {
				b.WriteString("    " + infoStyle.Render(bitsCompilerHint) + "\n")
			}
		} else {
			if r.Expected != "" {
				b.WriteString(fmt.Sprintf("  expected: %s\n", oneLine(r.Expected)))
			}
			b.WriteString(fmt.Sprintf("  got:      %s\n", oneLine(r.Got)))
		}
	}
	b.WriteString(fmt.Sprintf("\nPassed %d/%d\n", passed, len(results)))
	return b.String()
}

func appendLabeledBlock(b *strings.Builder, label, text string) {
	b.WriteString("  " + label + ":\n")
	if strings.TrimSpace(text) == "" {
		b.WriteString("    (empty)\n")
		return
	}
	for _, line := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			b.WriteString("    \n")
		} else {
			b.WriteString("    " + line + "\n")
		}
	}
}

func formatCustomRunOutput(input, out string, err error) string {
	var b strings.Builder
	b.WriteString("Input:\n")
	if strings.TrimSpace(input) == "" {
		b.WriteString("(empty)\n")
	} else {
		b.WriteString(input)
		if !strings.HasSuffix(input, "\n") {
			b.WriteString("\n")
		}
	}
	b.WriteString("\nOutput:\n")
	if err != nil {
		errText := err.Error()
		if strings.Contains(errText, bitsCompilerHint) {
			errText = strings.ReplaceAll(errText, bitsCompilerHint, infoStyle.Render(bitsCompilerHint))
		}
		b.WriteString(errText)
	} else if out == "" {
		b.WriteString("(empty)")
	} else {
		b.WriteString(out)
	}
	return b.String()
}

func oneLine(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	first := strings.TrimSpace(strings.Split(s, "\n")[0])
	if strings.Contains(s, "\n") {
		return first + " ..."
	}
	return first
}

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	cursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	passStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	infoStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("45"))
	errorStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
