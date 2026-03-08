package main

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type screen int

const (
	screenMain screen = iota
	screenProblem
	screenView
	screenToast
	screenCompiling
	screenCustom
)

type compileTarget int

const (
	compilePreset compileTarget = iota
	compileCustom
)

const (
	optionOpenCode          = "Open code"
	optionOpenDevlog        = "Open devlog"
	optionOpenJudge         = "Open judge"
	optionRunPresetTestCase = "Run preset test cases"
	optionRunCustomTestCase = "Run custom test cases"
)

var problemOptions = []string{
	optionOpenCode,
	optionOpenDevlog,
	optionOpenJudge,
	optionRunPresetTestCase,
	optionRunCustomTestCase,
}

type model struct {
	problems []Problem
	screen   screen

	cursor       int
	offset       int
	optionCursor int

	selected   Problem
	width      int
	height     int
	viewLine   int
	output     string
	toastText  string
	toastIsErr bool
	toastBack  screen
	toastID    int

	compileCancel context.CancelFunc
	compileBack   screen
	compileJobID  int

	// Custom interactive run
	customCmd       *exec.Cmd
	customStdin     io.WriteCloser
	customOutputCh  <-chan string
	customCleanup   func()
	customInputBuf  string
	customAllInput  string
	customAllOutput string
	customDisplay   string
	customCursorOn  bool
	customBlinkID   int
}

type toastTimeoutMsg struct{ id int }
type customOutputMsg string
type customDoneMsg struct{}
type cursorBlinkMsg struct{ id int }
type compileDoneMsg struct {
	id      int
	target  compileTarget
	exePath string
	cleanup func()
	err     error
}

func NewModel(problems []Problem) model {
	return model{problems: problems}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.adjustOffset()
		m.clampViewLine()
		return m, nil

	case toastTimeoutMsg:
		if m.screen == screenToast && msg.id == m.toastID {
			m.screen = m.toastBack
		}
		return m, nil

	case customOutputMsg:
		if m.screen != screenCustom {
			return m, nil
		}
		chunk := string(msg)
		m.customAllOutput += chunk
		m.customDisplay += chunk
		return m, waitForOutput(m.customOutputCh)

	case customDoneMsg:
		if m.screen != screenCustom {
			return m, nil
		}
		m.showCustomResult()
		return m, nil

	case cursorBlinkMsg:
		if m.screen != screenCustom || msg.id != m.customBlinkID {
			return m, nil
		}
		m.customCursorOn = !m.customCursorOn
		return m, blinkCursor(m.customBlinkID)

	case compileDoneMsg:
		if msg.id != m.compileJobID {
			if msg.cleanup != nil {
				msg.cleanup()
			}
			return m, nil
		}

		m.compileCancel = nil
		if m.screen != screenCompiling {
			if msg.cleanup != nil {
				msg.cleanup()
			}
			return m, nil
		}

		if msg.err != nil {
			m.screen = m.compileBack
			return m, m.showToast(msg.err.Error(), true, m.compileBack)
		}

		switch msg.target {
		case compilePreset:
			res := RunPresetTestCasesWithExe(m.selected, msg.exePath)
			if msg.cleanup != nil {
				msg.cleanup()
			}
			if res.Err != nil {
				m.screen = m.compileBack
				return m, m.showToast(res.Err.Error(), true, m.compileBack)
			}
			m.output = res.Summary
			m.viewLine = 0
			m.screen = screenView
			return m, nil
		case compileCustom:
			return m.startCustomWithExe(msg.exePath, msg.cleanup)
		default:
			if msg.cleanup != nil {
				msg.cleanup()
			}
			m.screen = m.compileBack
			return m, m.showToast("unknown compile target", true, m.compileBack)
		}

	case tea.KeyMsg:
		key := msg.String()

		if m.screen == screenCustom {
			return m.updateCustom(key)
		}

		if key == "ctrl+c" {
			if m.compileCancel != nil {
				m.compileCancel()
				m.compileCancel = nil
			}
			return m, tea.Quit
		}
		if key == "q" && m.screen == screenMain {
			return m, tea.Quit
		}

		switch key {
		case "up", "k":
			switch m.screen {
			case screenMain:
				if m.cursor > 0 {
					m.cursor--
				}
				m.adjustOffset()
			case screenProblem:
				if m.optionCursor > 0 {
					m.optionCursor--
				}
			case screenView:
				m.scrollView(-1)
			}
		case "down", "j":
			switch m.screen {
			case screenMain:
				if m.cursor < len(m.problems)-1 {
					m.cursor++
				}
				m.adjustOffset()
			case screenProblem:
				if m.optionCursor < len(problemOptions)-1 {
					m.optionCursor++
				}
			case screenView:
				m.scrollView(1)
			}
		case "enter":
			switch m.screen {
			case screenMain:
				if len(m.problems) == 0 {
					return m, nil
				}
				m.selected = m.problems[m.cursor]
				m.optionCursor = 0
				m.screen = screenProblem
			case screenProblem:
				option := problemOptions[m.optionCursor]
				switch option {
				case optionOpenCode:
					codeBytes, err := os.ReadFile(m.selected.CodePath)
					if err != nil {
						return m, m.showToast(err.Error(), true, screenProblem)
					}
					codeText := string(codeBytes)
					if strings.TrimSpace(codeText) == "" {
						return m, m.showToast("Code file is empty", true, screenProblem)
					}
					m.output = codeText
					m.viewLine = 0
					m.screen = screenView
				case optionOpenDevlog:
					devlogBytes, err := os.ReadFile(m.selected.DevlogPath)
					if err != nil {
						return m, m.showToast(err.Error(), true, screenProblem)
					}
					devlogText := string(devlogBytes)
					if strings.TrimSpace(devlogText) == "" {
						return m, m.showToast("Devlog file is empty", true, screenProblem)
					}
					m.output = devlogText
					m.viewLine = 0
					m.screen = screenView
				case optionOpenJudge:
					if m.selected.JudgeLink == "" {
						return m, m.showToast("judge link not found", true, screenProblem)
					}
					if err := OpenInSystem(m.selected.JudgeLink); err != nil {
						return m, m.showToast(err.Error(), true, screenProblem)
					}
				case optionRunPresetTestCase:
					return m, m.startCompile(compilePreset, screenProblem)
				case optionRunCustomTestCase:
					return m, m.startCompile(compileCustom, screenProblem)
				}
			}
		case "esc":
			switch m.screen {
			case screenProblem:
				m.screen = screenMain
			case screenView:
				m.screen = screenProblem
			case screenCompiling:
				if m.compileCancel != nil {
					m.compileCancel()
					m.compileCancel = nil
				}
				m.screen = m.compileBack
			case screenToast:
				m.screen = m.toastBack
			}
		}
	}
	return m, nil
}

func (m *model) startCompile(target compileTarget, back screen) tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background())
	m.compileCancel = cancel
	m.compileBack = back
	m.compileJobID++
	id := m.compileJobID
	sourcePath := m.selected.CodePath
	m.screen = screenCompiling

	return func() tea.Msg {
		exePath, cleanup, err := compileCPPContext(ctx, sourcePath)
		return compileDoneMsg{id: id, target: target, exePath: exePath, cleanup: cleanup, err: err}
	}
}

func (m model) startCustomWithExe(exePath string, cleanup func()) (tea.Model, tea.Cmd) {
	cmd := exec.Command(exePath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		cleanup()
		return m, m.showToast(err.Error(), true, screenProblem)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cleanup()
		return m, m.showToast(err.Error(), true, screenProblem)
	}
	if err := cmd.Start(); err != nil {
		cleanup()
		return m, m.showToast(err.Error(), true, screenProblem)
	}

	ch := make(chan string, 64)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, readErr := stdout.Read(buf)
			if n > 0 {
				ch <- string(buf[:n])
			}
			if readErr != nil {
				close(ch)
				return
			}
		}
	}()

	m.customCmd = cmd
	m.customStdin = stdin
	m.customOutputCh = ch
	m.customCleanup = cleanup
	m.customInputBuf = ""
	m.customAllInput = ""
	m.customAllOutput = ""
	m.customDisplay = ""
	m.customCursorOn = true
	m.customBlinkID++
	m.screen = screenCustom
	return m, tea.Batch(waitForOutput(ch), blinkCursor(m.customBlinkID))
}

func (m model) updateCustom(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c", "esc":
		m.cleanupCustom()
		m.screen = screenProblem
		return m, nil
	case "enter":
		line := m.customInputBuf + "\n"
		m.customAllInput += line
		m.customDisplay += line
		m.customInputBuf = ""
		if m.customStdin != nil {
			_, _ = io.WriteString(m.customStdin, line)
		}
		return m, nil
	case "backspace":
		if len(m.customInputBuf) > 0 {
			runes := []rune(m.customInputBuf)
			m.customInputBuf = string(runes[:len(runes)-1])
		}
		return m, nil
	case "space":
		m.customInputBuf += " "
		return m, nil
	case "tab":
		m.customInputBuf += "\t"
		return m, nil
	default:
		if len(key) == 1 && key[0] >= 32 && key[0] < 127 {
			m.customInputBuf += key
		}
		return m, nil
	}
}

func (m *model) cleanupCustom() {
	if m.customStdin != nil {
		_ = m.customStdin.Close()
	}
	if m.customCmd != nil && m.customCmd.Process != nil {
		_ = m.customCmd.Process.Kill()
		_ = m.customCmd.Wait()
	}
	if m.customCleanup != nil {
		m.customCleanup()
	}
	m.customCmd = nil
	m.customStdin = nil
	m.customOutputCh = nil
	m.customCleanup = nil
	m.customBlinkID++
}

func (m *model) showCustomResult() {
	var result strings.Builder
	result.WriteString("─── Input ───\n")
	if m.customAllInput == "" {
		result.WriteString("(no input)\n")
	} else {
		result.WriteString(m.customAllInput)
		if !strings.HasSuffix(m.customAllInput, "\n") {
			result.WriteString("\n")
		}
	}
	result.WriteString("\n─── Output ───\n")
	if m.customAllOutput == "" {
		result.WriteString("(no output)\n")
	} else {
		result.WriteString(m.customAllOutput)
	}
	m.output = result.String()
	m.viewLine = 0
	m.cleanupCustom()
	m.screen = screenView
}

func waitForOutput(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		data, ok := <-ch
		if !ok {
			return customDoneMsg{}
		}
		return customOutputMsg(data)
	}
}

func blinkCursor(id int) tea.Cmd {
	return tea.Tick(530*time.Millisecond, func(time.Time) tea.Msg {
		return cursorBlinkMsg{id: id}
	})
}

func (m *model) showToast(text string, isErr bool, back screen) tea.Cmd {
	m.toastText = text
	m.toastIsErr = isErr
	m.toastBack = back
	m.screen = screenToast
	m.toastID++
	id := m.toastID
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return toastTimeoutMsg{id: id}
	})
}

func (m *model) adjustOffset() {
	visible := m.visibleRows()
	if visible <= 0 {
		m.offset = 0
		return
	}

	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
}

func (m model) visibleRows() int {
	rows := m.height - 6
	if rows < 1 {
		rows = 1
	}
	return rows
}

func (m model) View() string {
	switch m.screen {
	case screenMain:
		title := titleStyle.Render("Competitive Programming Problems")
		var list string
		if len(m.problems) == 0 {
			list = errorStyle.Render("No .cpp problems found in ./problems")
		} else {
			visible := m.visibleRows()
			start := m.offset
			end := start + visible
			if end > len(m.problems) {
				end = len(m.problems)
			}

			lines := make([]string, 0, end-start)
			for i := start; i < end; i++ {
				line := m.problems[i].ProblemCode
				if i == m.cursor {
					line = cursorStyle.Render("> " + line)
				} else {
					line = "  " + line
				}
				lines = append(lines, line)
			}
			list = strings.Join(lines, "\n")
		}
		footer := footerStyle.Render("q: quit · enter: select")
		return strings.Join([]string{title, list, footer}, "\n")

	case screenToast:
		msg := m.toastText
		if m.toastIsErr {
			msg = errorStyle.Render(msg)
		} else {
			msg = statusStyle.Render(msg)
		}
		footer := footerStyle.Render("esc: go back")
		return strings.Join([]string{msg, footer}, "\n")

	case screenCompiling:
		body := statusStyle.Render("Compiling...")
		footer := footerStyle.Render("esc: cancel and go back")
		return strings.Join([]string{body, "", footer}, "\n")

	case screenView:
		var body string
		if m.output == "" {
			body = statusStyle.Render("no content loaded")
		} else {
			lines := m.viewLines()
			visible := m.viewVisibleRows()
			maxStart := len(lines) - visible
			if maxStart < 0 {
				maxStart = 0
			}

			start := m.viewLine
			if start < 0 {
				start = 0
			}
			if start > maxStart {
				start = maxStart
			}
			end := start + visible
			if end > len(lines) {
				end = len(lines)
			}
			body = strings.Join(lines[start:end], "\n")
		}
		footer := footerStyle.Render("esc: go back · j/k: scroll")
		return strings.Join([]string{body, "", footer}, "\n")

	case screenProblem:
		title := titleStyle.Render(m.selected.ProblemCode)
		lines := make([]string, 0, len(problemOptions))
		for i, option := range problemOptions {
			line := option
			if i == m.optionCursor {
				line = cursorStyle.Render("> " + line)
			} else {
				line = "  " + line
			}
			lines = append(lines, line)
		}
		list := strings.Join(lines, "\n")
		footer := footerStyle.Render("esc: go back · enter: select")
		return strings.Join([]string{title, list, footer}, "\n")

	case screenCustom:
		viewLines := m.customViewLines()
		visible := m.height - 2
		if visible < 1 {
			visible = 1
		}
		start := len(viewLines) - visible
		if start < 0 {
			start = 0
		}
		end := start + visible
		if end > len(viewLines) {
			end = len(viewLines)
		}
		body := strings.Join(viewLines[start:end], "\n")
		footer := footerStyle.Render("esc: quit program")
		return body + "\n" + footer

	default:
		return ""
	}
}

func (m model) customViewLines() []string {
	cursor := "█"
	if !m.customCursorOn {
		cursor = " "
	}
	full := m.customDisplay + m.customInputBuf + cursor

	clean := strings.ReplaceAll(full, "\r\n", "\n")
	clean = strings.ReplaceAll(clean, "\r", "\n")
	rawLines := strings.Split(clean, "\n")

	wrapWidth := m.width
	if wrapWidth < 1 {
		wrapWidth = 1
	}

	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		lines = append(lines, wrapLine(line, wrapWidth)...)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func (m *model) scrollView(delta int) {
	if m.output == "" {
		return
	}

	lines := m.viewLines()
	visible := m.viewVisibleRows()
	maxStart := len(lines) - visible
	if maxStart < 0 {
		maxStart = 0
	}

	m.viewLine += delta
	if m.viewLine < 0 {
		m.viewLine = 0
	}
	if m.viewLine > maxStart {
		m.viewLine = maxStart
	}
}

func (m *model) clampViewLine() {
	if m.output == "" {
		m.viewLine = 0
		return
	}

	lines := m.viewLines()
	visible := m.viewVisibleRows()
	maxStart := len(lines) - visible
	if maxStart < 0 {
		maxStart = 0
	}

	if m.viewLine < 0 {
		m.viewLine = 0
	}
	if m.viewLine > maxStart {
		m.viewLine = maxStart
	}
}

func (m model) viewLines() []string {
	cleanOutput := strings.ReplaceAll(m.output, "\r\n", "\n")
	cleanOutput = strings.ReplaceAll(cleanOutput, "\r", "\n")
	rawLines := strings.Split(cleanOutput, "\n")
	wrapWidth := m.width
	if wrapWidth < 1 {
		wrapWidth = 1
	}

	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		lines = append(lines, wrapLine(line, wrapWidth)...)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func (m model) viewVisibleRows() int {
	rows := m.height - 2
	if rows < 1 {
		rows = 1
	}
	return rows
}

func wrapLine(line string, width int) []string {
	if width < 1 {
		width = 1
	}
	if line == "" {
		return []string{""}
	}

	runes := []rune(line)
	out := make([]string, 0, (len(runes)/width)+1)
	start := 0
	for start < len(runes) {
		col := 0
		end := start
		lastSpace := -1
		for end < len(runes) {
			rw := runewidth.RuneWidth(runes[end])
			if rw <= 0 {
				rw = 1
			}
			if col+rw > width {
				break
			}
			col += rw
			if unicode.IsSpace(runes[end]) {
				lastSpace = end
			}
			end++
		}

		if end == len(runes) {
			out = append(out, string(runes[start:end]))
			break
		}

		if lastSpace >= start {
			segment := strings.TrimRightFunc(string(runes[start:lastSpace]), unicode.IsSpace)
			out = append(out, segment)
			start = lastSpace + 1
			for start < len(runes) && unicode.IsSpace(runes[start]) {
				start++
			}
			continue
		}

		out = append(out, string(runes[start:end]))
		start = end
	}

	if len(out) == 0 {
		return []string{""}
	}
	return out
}

var (
	titleStyle  = lipgloss.NewStyle().Bold(true)
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
