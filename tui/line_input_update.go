package tui

import tea "charm.land/bubbletea/v2"

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
