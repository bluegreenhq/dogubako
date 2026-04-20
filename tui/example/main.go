package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/bluegreenhq/dogubako/tui"
)

// demo はデモの種類を表す。
type demo int

const (
	demoLineInput demo = iota
	demoBoxButton
	demoPopupMenu
	demoConfirmDialog
	demoCount
)

var demoNames = [demoCount]string{
	"LineInput",
	"BoxButton",
	"PopupMenu",
	"ConfirmDialog",
}

type model struct {
	current demo
	width   int
	height  int

	// LineInput (blink有効/無効)
	inputs      [2]tui.LineInput
	inputFocus  int
	cursorBlink tui.CursorBlink

	// BoxButton
	buttons [3]tui.BoxButton

	// PopupMenu
	popup *tui.PopupMenu

	// ConfirmDialog
	confirm       tui.ConfirmDialog
	confirmResult string
}

const cursorOwner = 1

func newModel() model {
	return model{
		current: demoLineInput,
		width:   0,
		height:  0,
		inputs: [2]tui.LineInput{
			tui.NewLineInput(),
			tui.NewLineInputNoBlink(),
		},
		inputFocus:    0,
		cursorBlink:   tui.NewCursorBlink(cursorOwner),
		confirmResult: "",
		buttons: [3]tui.BoxButton{
			tui.NewBoxButton("OK"),
			tui.NewBoxButton("キャンセル"),
			tui.NewBoxButton("Apply"),
		},
		popup: tui.NewPopupMenu([]tui.MenuItem{
			tui.NewMenuItem("追加"),
			tui.NewMenuItem("編集"),
			tui.NewDisabledMenuItem("削除 (disabled)"),
			tui.NewMenuItem("閉じる"),
		}),
		confirm: tui.NewConfirmDialog("確認", "この操作を実行しますか？"),
	}
}

func (m model) Init() tea.Cmd {
	return m.cursorBlink.Reset()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.confirm.SetScreenSize(msg.Width, msg.Height)

		return m, nil

	case tui.CursorBlinkMsg:
		cmd := m.cursorBlink.HandleMsg(msg)

		return m, cmd

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Code == 'q' || (msg.Code == 'c' && msg.Mod == tea.ModCtrl):
		return m, tea.Quit

	case msg.Code == tea.KeyTab && msg.Mod == 0:
		m.current = (m.current + 1) % demoCount

		return m, m.cursorBlink.Reset()

	case msg.Code == tea.KeyTab && msg.Mod == tea.ModShift:
		m.current = (m.current - 1 + demoCount) % demoCount

		return m, m.cursorBlink.Reset()
	}

	switch m.current {
	case demoLineInput:
		switch {
		case msg.Code == tea.KeyDown,
			msg.Code == 'n' && msg.Mod&tea.ModCtrl != 0:
			if m.inputFocus < len(m.inputs)-1 {
				m.inputFocus++
			}
		case msg.Code == tea.KeyUp,
			msg.Code == 'p' && msg.Mod&tea.ModCtrl != 0:
			if m.inputFocus > 0 {
				m.inputFocus--
			}
		default:
			m.inputs[m.inputFocus].HandleKey(msg)
		}

		return m, m.cursorBlink.Reset()

	case demoBoxButton:

	case demoPopupMenu:
		m.popup.HandleKeyNav(msg)

	case demoConfirmDialog:
		result := m.confirm.Update(msg)
		switch result {
		case tui.ConfirmYes:
			m.confirmResult = "Yes"
		case tui.ConfirmNo:
			m.confirmResult = "No"
		case tui.ConfirmContinue:
		}

	case demoCount:
	}

	return m, nil
}

var (
	tabStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Padding(0, 1)
	tabActiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("4")).Padding(0, 1).Bold(true)
	titleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	helpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func (m model) View() tea.View {
	var b strings.Builder

	// タブバー
	for i := range demoCount {
		style := tabStyle
		if i == m.current {
			style = tabActiveStyle
		}

		b.WriteString(style.Render(demoNames[i]))
		b.WriteString(" ")
	}

	b.WriteString("\n\n")

	// デモコンテンツ
	switch m.current {
	case demoLineInput:
		b.WriteString(m.viewLineInput())
	case demoBoxButton:
		b.WriteString(m.viewBoxButton())
	case demoPopupMenu:
		b.WriteString(m.viewPopupMenu())
	case demoConfirmDialog:
		b.WriteString(m.viewConfirmDialog())
	case demoCount:
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Tab/Shift+Tab: switch demo  q: quit"))

	v := tea.NewView(b.String())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion

	return v
}

var inputLabels = [2]string{"blink:    ", "no blink: "}

func (m model) viewLineInput() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("LineInput (複数フィールド)"))
	b.WriteString("\n\n")

	for i := range m.inputs {
		label := inputLabels[i]
		if i == m.inputFocus {
			b.WriteString("  " + label + m.inputs[i].ViewWithWidth(0, m.cursorBlink.Visible()))
		} else {
			v := m.inputs[i].Value()
			if v == "" {
				v = " "
			}

			b.WriteString("  " + label + v)
		}

		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  ↑/↓: フィールド移動  C-a: 先頭  C-e: 末尾  C-f/C-b: 左右"))

	return b.String()
}

func (m model) viewBoxButton() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("BoxButton"))
	b.WriteString("\n\n")

	tops := make([]string, len(m.buttons))
	mids := make([]string, len(m.buttons))
	bots := make([]string, len(m.buttons))

	for i := range m.buttons {
		tops[i] = m.buttons[i].ViewTop()
		mids[i] = m.buttons[i].ViewMiddle()
		bots[i] = m.buttons[i].ViewBottom()
	}

	gap := "  "
	b.WriteString("  " + strings.Join(tops, gap) + "\n")
	b.WriteString("  " + strings.Join(mids, gap) + "\n")
	b.WriteString("  " + strings.Join(bots, gap))

	return b.String()
}

func (m model) viewPopupMenu() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("PopupMenu"))
	b.WriteString("\n\n")

	lines := m.popup.View()
	for _, line := range lines {
		b.WriteString("  " + line + "\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  j/↓: 下  k/↑: 上"))

	return b.String()
}

func (m model) viewConfirmDialog() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("ConfirmDialog"))
	b.WriteString("\n\n")
	b.WriteString(m.confirm.View())

	if m.confirmResult != "" {
		b.WriteString("\n\n")
		fmt.Fprintf(&b, "  Result: %s", m.confirmResult)
	}

	return b.String()
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
