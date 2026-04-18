package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

const blinkInterval = 500 * time.Millisecond

// CursorBlinkMsg はカーソルの点滅状態を切り替えるメッセージ。
type CursorBlinkMsg struct {
	Owner int
	Tag   int
}

// CursorBlink はカーソルの点滅状態を管理する。
type CursorBlink struct {
	visible bool
	tag     int
	owner   int
}

// NewCursorBlink は新しい CursorBlink を生成する。
func NewCursorBlink(owner int) CursorBlink {
	return CursorBlink{visible: true, tag: 0, owner: owner}
}

// Visible はカーソルが表示状態かを返す。
func (cb *CursorBlink) Visible() bool { return cb.visible }

// Reset はカーソルを表示状態にリセットし、新しい blink タイマーを開始する。
func (cb *CursorBlink) Reset() tea.Cmd {
	cb.visible = true
	cb.tag++
	tag := cb.tag
	owner := cb.owner

	return tea.Tick(blinkInterval, func(_ time.Time) tea.Msg {
		return CursorBlinkMsg{Owner: owner, Tag: tag}
	})
}

// Stop はカーソルを表示状態にし、blink を停止する。
func (cb *CursorBlink) Stop() {
	cb.visible = true
	cb.tag++
}

// HandleMsg は CursorBlinkMsg を処理して blink 状態を切り替える。
// Owner または Tag が一致しない場合は無視して nil を返す。
func (cb *CursorBlink) HandleMsg(msg CursorBlinkMsg) tea.Cmd {
	if msg.Owner != cb.owner || msg.Tag != cb.tag {
		return nil
	}

	cb.visible = !cb.visible
	tag := cb.tag
	owner := cb.owner

	return tea.Tick(blinkInterval, func(_ time.Time) tea.Msg {
		return CursorBlinkMsg{Owner: owner, Tag: tag}
	})
}
