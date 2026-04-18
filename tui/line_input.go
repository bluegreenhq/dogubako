package tui

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
		return " " + before + "█"
	}

	after := string(li.value[li.cursor+1:])

	return " " + before + "█" + after
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

func (li *LineInput) viewContent(runes []rune, cursorPos int, cursorVisible bool) string {
	if !cursorVisible {
		return string(runes)
	}

	before := string(runes[:cursorPos])
	if cursorPos >= len(runes) {
		return before + "█"
	}

	after := string(runes[cursorPos+1:])

	return before + "█" + after
}
