package tui_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"

	"github.com/bluegreenhq/dogubako/tui"
)

func TestLineInput_SetValueAndValue(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("hello")
	assert.Equal(t, "hello", li.Value())
}

func TestLineInput_Reset(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("hello")
	li.Reset()
	assert.Equal(t, "", li.Value())
}

func TestLineInput_View_CursorAtEnd(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("abc")
	// カーソルは末尾、先頭スペース + テキスト + ブロックカーソル
	assert.Equal(t, " abc█", li.View(true))
}

func TestLineInput_View_CursorHidden(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("abc")
	assert.Equal(t, " abc", li.View(false))
}

func TestLineInput_View_CursorAtBeginning(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("abc")
	li.HandleKey(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl}) // C-a で先頭へ
	assert.Equal(t, " █bc", li.View(true))
}

func TestLineInput_View_CursorAtMiddle(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("abc")
	li.HandleKey(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
	li.HandleKey(tea.KeyPressMsg{Code: 'f', Mod: tea.ModCtrl}) // C-f で1つ右へ
	assert.Equal(t, " a█c", li.View(true))
}

func TestLineInput_View_WideCursorWidth(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("三井住友FG")
	li.HandleKey(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl}) // 先頭へ

	visible := li.View(true)
	hidden := li.View(false)
	// blink on/off で表示幅が変わらないこと
	assert.Equal(t, lipgloss.Width(hidden), lipgloss.Width(visible))
}

func TestLineInput_ViewWithWidth_WideCursorWidth(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("三井住友FG")
	li.HandleKey(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})

	visible := li.ViewWithWidth(0, true)
	hidden := li.ViewWithWidth(0, false)
	assert.Equal(t, lipgloss.Width(hidden), lipgloss.Width(visible))
}

func TestLineInput_HandleKey_Insert(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.HandleKey(tea.KeyPressMsg{Code: 'h', Text: "h"})
	li.HandleKey(tea.KeyPressMsg{Code: 'i', Text: "i"})
	assert.Equal(t, "hi", li.Value())
}

func TestLineInput_HandleKey_Backspace(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("abc")
	li.HandleKey(tea.KeyPressMsg{Code: tea.KeyBackspace})
	assert.Equal(t, "ab", li.Value())
}

func TestLineInput_HandleKey_Delete(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("abc")
	li.HandleKey(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl}) // 先頭へ
	li.HandleKey(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl}) // C-d
	assert.Equal(t, "bc", li.Value())
}

func TestLineInput_HandleKey_KillYank(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("abcdef")
	li.HandleKey(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl}) // 先頭へ
	li.HandleKey(tea.KeyPressMsg{Code: 'f', Mod: tea.ModCtrl}) // 1つ右
	li.HandleKey(tea.KeyPressMsg{Code: 'f', Mod: tea.ModCtrl}) // 2つ右
	li.HandleKey(tea.KeyPressMsg{Code: 'k', Mod: tea.ModCtrl}) // kill to end
	assert.Equal(t, "ab", li.Value())

	li.HandleKey(tea.KeyPressMsg{Code: 'y', Mod: tea.ModCtrl}) // yank
	assert.Equal(t, "abcdef", li.Value())
}

func TestLineInput_HandleKey_Enter(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	result := li.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	assert.Equal(t, tui.LineInputSubmit, result)
}

func TestLineInput_HandleKey_Escape(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	result := li.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	assert.Equal(t, tui.LineInputCancel, result)
}

func TestLineInput_HandleKey_CursorMovement(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("abc")

	// Home で先頭
	li.HandleKey(tea.KeyPressMsg{Code: tea.KeyHome})
	assert.Equal(t, " █bc", li.View(true))

	// End で末尾
	li.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnd})
	assert.Equal(t, " abc█", li.View(true))

	// Left で左
	li.HandleKey(tea.KeyPressMsg{Code: tea.KeyLeft})
	assert.Equal(t, " ab█", li.View(true))

	// Right で右
	li.HandleKey(tea.KeyPressMsg{Code: tea.KeyRight})
	assert.Equal(t, " abc█", li.View(true))
}

func TestLineInput_HandleKey_CursorBoundary(t *testing.T) {
	t.Parallel()

	li := tui.NewLineInput()
	li.SetValue("a")

	// 末尾でさらに右に移動しても変わらない
	li.HandleKey(tea.KeyPressMsg{Code: tea.KeyRight})
	assert.Equal(t, " a█", li.View(true))

	// 先頭でさらに左に移動しても変わらない
	li.HandleKey(tea.KeyPressMsg{Code: tea.KeyHome})
	li.HandleKey(tea.KeyPressMsg{Code: tea.KeyLeft})
	assert.Equal(t, " █", li.View(true))
}
