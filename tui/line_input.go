package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var cursorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("7"))

// LineInputResult は LineInput のキー処理結果を表す。
type LineInputResult int

const (
	LineInputNone   LineInputResult = iota // 通常の編集操作
	LineInputSubmit                        // Enter が押された
	LineInputCancel                        // Escape が押された
)

// LineInput は単一行のインライン入力を管理する。
// カーソル移動・kill/yank など基本的な編集操作をサポートする。
type LineInput struct {
	value   []rune
	cursor  int
	killBuf []rune
}

// NewLineInput は新しい LineInput を生成する。
func NewLineInput() LineInput {
	return LineInput{value: nil, cursor: 0, killBuf: nil}
}

// Value は入力中のテキストを返す。
func (li *LineInput) Value() string { return string(li.value) }

// SetValue はテキストを設定し、カーソルを末尾に移動する。
func (li *LineInput) SetValue(s string) {
	li.value = []rune(s)
	li.cursor = len(li.value)
}

// Reset は入力をクリアする。
func (li *LineInput) Reset() {
	li.value = nil
	li.cursor = 0
}

// View はカーソル付きの表示文字列を返す。先頭にスペースを付与する。
// カーソル位置の文字をブロックカーソル（█）で置換して表示する。
// cursorVisible が false の場合はカーソルを非表示にする。
func (li *LineInput) View(cursorVisible bool) string {
	if !cursorVisible {
		return " " + string(li.value)
	}

	before := string(li.value[:li.cursor])
	if li.cursor >= len(li.value) {
		return " " + before + cursorStyle.Render(" ")
	}

	cursor := cursorStyle.Render(string(li.value[li.cursor]))
	after := string(li.value[li.cursor+1:])

	return " " + before + cursor + after
}

// ViewWithWidth は最大表示幅を考慮した表示文字列を返す。
// カーソル位置が常に見えるよう、左側を切り詰める。
func (li *LineInput) ViewWithWidth(maxWidth int, cursorVisible bool) string {
	runes := li.value
	cursor := li.cursor

	if maxWidth <= 0 || len(runes) <= maxWidth {
		return li.viewContent(runes, cursor, cursorVisible)
	}

	// カーソルが見えるようにウィンドウをスライド
	start := 0
	if cursor > maxWidth-1 {
		start = cursor - maxWidth + 1
	}

	end := min(start+maxWidth, len(runes))

	return li.viewContent(runes[start:end], cursor-start, cursorVisible)
}

// HandleKey はキー入力を処理し、結果を返す。
func (li *LineInput) HandleKey(msg tea.KeyPressMsg) LineInputResult { //nolint:cyclop // キーバインド分岐
	switch {
	case msg.Code == tea.KeyEnter:
		return LineInputSubmit
	case msg.Code == tea.KeyEscape:
		return LineInputCancel
	case msg.Code == 'a' && msg.Mod == tea.ModCtrl:
		li.cursor = 0
	case msg.Code == 'e' && msg.Mod == tea.ModCtrl:
		li.cursor = len(li.value)
	case msg.Code == 'f' && msg.Mod == tea.ModCtrl:
		li.cursorRight()
	case msg.Code == 'b' && msg.Mod == tea.ModCtrl:
		li.cursorLeft()
	case msg.Code == 'd' && msg.Mod == tea.ModCtrl:
		li.delete()
	case msg.Code == 'k' && msg.Mod == tea.ModCtrl:
		li.killToEnd()
	case msg.Code == 'y' && msg.Mod == tea.ModCtrl:
		li.yank()
	case msg.Code == tea.KeyBackspace, msg.Code == 'h' && msg.Mod == tea.ModCtrl:
		li.backspace()
	case msg.Code == tea.KeyDelete:
		li.delete()
	case msg.Code == tea.KeyLeft:
		li.cursorLeft()
	case msg.Code == tea.KeyRight:
		li.cursorRight()
	case msg.Code == tea.KeyHome:
		li.cursor = 0
	case msg.Code == tea.KeyEnd:
		li.cursor = len(li.value)
	default:
		if msg.Text != "" && (msg.Mod == 0 || msg.Mod == tea.ModShift) {
			li.insertText(msg.Text)
		}
	}

	return LineInputNone
}

func (li *LineInput) viewContent(runes []rune, cursorPos int, cursorVisible bool) string {
	if !cursorVisible {
		return string(runes)
	}

	before := string(runes[:cursorPos])
	if cursorPos >= len(runes) {
		return before + cursorStyle.Render(" ")
	}

	cursor := cursorStyle.Render(string(runes[cursorPos]))
	after := string(runes[cursorPos+1:])

	return before + cursor + after
}

func (li *LineInput) insertText(s string) {
	runes := []rune(s)
	newValue := make([]rune, 0, len(li.value)+len(runes))
	newValue = append(newValue, li.value[:li.cursor]...)
	newValue = append(newValue, runes...)
	newValue = append(newValue, li.value[li.cursor:]...)
	li.value = newValue
	li.cursor += len(runes)
}

func (li *LineInput) backspace() {
	if li.cursor > 0 {
		li.value = append(li.value[:li.cursor-1], li.value[li.cursor:]...)
		li.cursor--
	}
}

func (li *LineInput) delete() {
	if li.cursor < len(li.value) {
		li.value = append(li.value[:li.cursor], li.value[li.cursor+1:]...)
	}
}

func (li *LineInput) killToEnd() {
	if li.cursor < len(li.value) {
		killed := make([]rune, len(li.value)-li.cursor)
		copy(killed, li.value[li.cursor:])
		li.killBuf = killed
		li.value = li.value[:li.cursor]
	}
}

func (li *LineInput) yank() {
	if len(li.killBuf) > 0 {
		li.insertText(string(li.killBuf))
	}
}

func (li *LineInput) cursorLeft() {
	if li.cursor > 0 {
		li.cursor--
	}
}

func (li *LineInput) cursorRight() {
	if li.cursor < len(li.value) {
		li.cursor++
	}
}
