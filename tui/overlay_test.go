package tui_test

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"

	"github.com/bluegreenhq/dogubako/tui"
)

func TestComposeLine_ASCII(t *testing.T) {
	t.Parallel()

	base := "ABCDEFGHIJ" // 幅10
	overlay := "XY"      // 幅2
	result := tui.ComposeLine(base, overlay, 3, 2)
	// "ABC" + reset + "XY" + reset + "FGHIJ"
	assert.Equal(t, "ABC\x1b[0mXY\x1b[0mFGHIJ", result)
}

func TestComposeLine_OverlayPadding(t *testing.T) {
	t.Parallel()

	base := "ABCDEFGHIJ"
	overlay := "X" // 幅1, overlayW=3 → "X  " にパディング
	result := tui.ComposeLine(base, overlay, 2, 3)
	// "AB" + reset + "X  " + reset + "FGHIJ"
	assert.Equal(t, "AB\x1b[0mX  \x1b[0mFGHIJ", result)
}

func TestComposeLine_BaseShort(t *testing.T) {
	t.Parallel()

	base := "AB" // 幅2, startX=5 → prefix パディング
	overlay := "XY"
	result := tui.ComposeLine(base, overlay, 5, 2)
	// "AB   " + reset + "XY" + reset + ""
	assert.Equal(t, "AB   \x1b[0mXY\x1b[0m", result)
}

func TestComposeLine_CJK(t *testing.T) {
	t.Parallel()

	base := "あいうえお" // 幅10（全角5文字）
	overlay := "XY" // 幅2
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

func TestClampMenuOrigin_FitsInScreen(t *testing.T) {
	t.Parallel()

	x, y := tui.ClampMenuOrigin(10, 5, 3, 2, 80, 24)
	assert.Equal(t, 3, x)
	assert.Equal(t, 2, y)
}

func TestClampMenuOrigin_ClampRight(t *testing.T) {
	t.Parallel()

	// menuW=10, anchorX=75, screenW=80 → 75+10=85 > 80 → x=70
	x, y := tui.ClampMenuOrigin(10, 5, 75, 2, 80, 24)
	assert.Equal(t, 70, x)
	assert.Equal(t, 2, y)
}

func TestClampMenuOrigin_ClampBottom(t *testing.T) {
	t.Parallel()

	// menuH=5, anchorY=22, screenH=24 → 22+5=27 > 24 → y=19
	x, y := tui.ClampMenuOrigin(10, 5, 3, 22, 80, 24)
	assert.Equal(t, 3, x)
	assert.Equal(t, 19, y)
}

func TestClampMenuOrigin_MenuLargerThanScreen(t *testing.T) {
	t.Parallel()

	// メニューが画面より大きい → 0にクランプ
	x, y := tui.ClampMenuOrigin(100, 30, 5, 5, 80, 24)
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)
}

func TestOverlayMenuOnBase_EmptyMenu(t *testing.T) {
	t.Parallel()

	menu := &stubMenu{w: 4, h: 0, lines: nil}
	base := "AAAA\nBBBB"
	result := tui.OverlayMenuOnBase(menu, base, 0, 0, 10, 2)
	assert.Equal(t, base, result)
}

func TestOverlayMenuOnBase_Basic(t *testing.T) {
	t.Parallel()

	menu := &stubMenu{w: 2, h: 1, lines: []string{"XY"}}
	base := "AAAA\nBBBB\nCCCC"
	result := tui.OverlayMenuOnBase(menu, base, 1, 1, 10, 3)
	lines := strings.Split(result, "\n")
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "AAAA", lines[0])
	assert.Contains(t, lines[1], "XY")
}

func TestOverlayGeometry_Contains(t *testing.T) {
	t.Parallel()

	geo := tui.OverlayGeometry{StartX: 10, StartY: 5, OverlayW: 20, OverlayH: 10}

	assert.True(t, geo.Contains(10, 5))
	assert.True(t, geo.Contains(29, 14))
	assert.False(t, geo.Contains(9, 5))
	assert.False(t, geo.Contains(30, 5))
	assert.False(t, geo.Contains(10, 4))
	assert.False(t, geo.Contains(10, 15))
}

func TestCalcOverlayGeometry_Centered(t *testing.T) {
	t.Parallel()

	// 幅4, 高さ1のオーバーレイを80x24画面に中央配置
	geo := tui.CalcOverlayGeometry("ABCD", 80, 24, 1, 2, 1)
	assert.Equal(t, (80-4)/2, geo.StartX)
	assert.Equal(t, (24-1)/2, geo.StartY)
	assert.Equal(t, geo.StartX+1+2, geo.ContentX)
	assert.Equal(t, geo.StartY+1+1, geo.ContentY)
	assert.Equal(t, 4, geo.OverlayW)
	assert.Equal(t, 1, geo.OverlayH)
}

func TestCalcOverlayGeometry_LargerThanScreen(t *testing.T) {
	t.Parallel()

	// 画面より大きい → 0にクランプ
	geo := tui.CalcOverlayGeometry("ABCDEFGHIJ", 5, 1, 1, 0, 0)
	assert.Equal(t, 0, geo.StartX)
	assert.Equal(t, 0, geo.StartY)
}

func TestOverlayCentered_Basic(t *testing.T) {
	t.Parallel()

	base := "AAAA\nBBBB\nCCCC\nDDDD\nEEEE"
	overlay := "XY"
	result := tui.OverlayCentered(base, overlay, 4, 5)
	lines := strings.Split(result, "\n")
	assert.Equal(t, 5, len(lines))
	// 中央行(index 2)にオーバーレイされる
	assert.Contains(t, lines[2], "XY")
	// 他の行は変更なし
	assert.Equal(t, "AAAA", lines[0])
	assert.Equal(t, "EEEE", lines[4])
}

func TestOverlayCentered_PadsShortBase(t *testing.T) {
	t.Parallel()

	base := "AA"
	result := tui.OverlayCentered(base, "X", 4, 3)
	lines := strings.Split(result, "\n")
	assert.Equal(t, 3, len(lines))
}

// stubMenu は OverlayMenu インターフェースのテスト用スタブ。
type stubMenu struct {
	w, h  int
	lines []string
}

func (s *stubMenu) Width() int     { return s.w }
func (s *stubMenu) Height() int    { return s.h }
func (s *stubMenu) View() []string { return s.lines }
