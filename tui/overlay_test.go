package tui_test

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"

	"github.com/bluegreenhq/dogubako/tui"
)

func TestComposeLine_ASCII(t *testing.T) {
	t.Parallel()

	base := "ABCDEFGHIJ" // 幅10
	overlay := "XY"       // 幅2
	result := tui.ComposeLine(base, overlay, 3, 2)
	// "ABC" + reset + "XY" + reset + "FGHIJ"
	assert.Equal(t, "ABC\x1b[0mXY\x1b[0mFGHIJ", result)
}

func TestComposeLine_OverlayPadding(t *testing.T) {
	t.Parallel()

	base := "ABCDEFGHIJ"
	overlay := "X"        // 幅1, overlayW=3 → "X  " にパディング
	result := tui.ComposeLine(base, overlay, 2, 3)
	// "AB" + reset + "X  " + reset + "FGHIJ"
	assert.Equal(t, "AB\x1b[0mX  \x1b[0mFGHIJ", result)
}

func TestComposeLine_BaseShort(t *testing.T) {
	t.Parallel()

	base := "AB"          // 幅2, startX=5 → prefix パディング
	overlay := "XY"
	result := tui.ComposeLine(base, overlay, 5, 2)
	// "AB   " + reset + "XY" + reset + ""
	assert.Equal(t, "AB   \x1b[0mXY\x1b[0m", result)
}

func TestComposeLine_CJK(t *testing.T) {
	t.Parallel()

	base := "あいうえお" // 幅10（全角5文字）
	overlay := "XY"       // 幅2
	result := tui.ComposeLine(base, overlay, 4, 2)
	// prefix: "あい"(幅4) + "XY" + suffix: "えお"(幅4)
	assert.Equal(t, "あい\x1b[0mXY\x1b[0mえお", result)
}

func TestComposeLine_CJKBoundary(t *testing.T) {
	t.Parallel()

	base := "あいうえお"
	overlay := "X"
	// startX=3: "��"(幅2) の次は "い"(幅4始まり) → prefix は "あ" + " " = 幅3
	result := tui.ComposeLine(base, overlay, 3, 1)
	// 結果の表示幅はベースと同じ10
	assert.Equal(t, 10, lipgloss.Width(result))
}

func TestComposeLine_Empty(t *testing.T) {
	t.Parallel()

	result := tui.ComposeLine("", "XY", 0, 2)
	assert.Equal(t, "\x1b[0mXY\x1b[0m", result)
}

func TestOverlayLines_Basic(t *testing.T) {
	t.Parallel()

	body := []string{
		"AAAAAAAAAA",
		"BBBBBBBBBB",
		"CCCCCCCCCC",
		"DDDDDDDDDD",
	}
	overlay := []string{
		"┌──┐",
		"│XY│",
		"└──┘",
	}
	tui.OverlayLines(body, overlay, 3, 1)

	assert.Equal(t, "AAAAAAAAAA", body[0]) // 変更なし
	assert.Contains(t, body[1], "┌──┐")
	assert.Contains(t, body[2], "│XY│")
	assert.Contains(t, body[3], "└──┘")
}

func TestOverlayLines_OutOfBounds(t *testing.T) {
	t.Parallel()

	body := []string{
		"AAAAAAAAAA",
		"BBBBBBBBBB",
	}
	overlay := []string{
		"XX",
		"YY",
		"ZZ", // body の範囲外 → 無視される
	}
	tui.OverlayLines(body, overlay, 0, 1)

	assert.Equal(t, "AAAAAAAAAA", body[0]) // 変更なし
	assert.Contains(t, body[1], "XX")
}

func TestOverlayLines_Empty(t *testing.T) {
	t.Parallel()

	body := []string{"AAAA"}
	tui.OverlayLines(body, nil, 0, 0)
	assert.Equal(t, "AAAA", body[0]) // 変更なし
}
