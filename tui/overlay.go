package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

const ansiReset = "\033[0m"

// ComposeLine はベース行の [startX, startX+overlayW) をオーバーレイ行で置き換える。
// CJK 文字の境界や ANSI シーケンスの破損を防ぐ。
func ComposeLine(baseLine, overlayLine string, startX, overlayW int) string {
	baseW := lipgloss.Width(baseLine)

	// 1. prefix（0〜startX）
	prefix := ansi.Truncate(baseLine, startX, "")

	prefixW := lipgloss.Width(prefix)
	if prefixW < startX {
		prefix += strings.Repeat(" ", startX-prefixW)
	}

	// 2. オーバーレイ行を overlayW 幅にパディング
	lineW := lipgloss.Width(overlayLine)
	if lineW < overlayW {
		overlayLine += strings.Repeat(" ", overlayW-lineW)
	}

	// 3. suffix（startX+overlayW〜末尾）
	rightStart := startX + overlayW

	suffix := ""
	if baseW > rightStart {
		suffix = ansi.Cut(baseLine, rightStart, baseW)
		suffixW := lipgloss.Width(suffix)

		expectedW := baseW - rightStart
		if suffixW > expectedW {
			extra := suffixW - expectedW
			suffix = strings.Repeat(" ", extra) + ansi.Cut(baseLine, rightStart+extra, baseW)
		} else if suffixW < expectedW {
			suffix = strings.Repeat(" ", expectedW-suffixW) + suffix
		}
	}

	return prefix + ansiReset + overlayLine + ansiReset + suffix
}

// OverlayMenu はアンカー位置に表示するポップアップメニューの共通インターフェース。
type OverlayMenu interface {
	Width() int
	Height() int
	View() []string
}

// ClampMenuOrigin はメニューが画面内に収まるようにアンカー座標をクランプする。
func ClampMenuOrigin(menuW, menuH, anchorX, anchorY, screenW, screenH int) (int, int) {
	x := anchorX
	y := anchorY

	if x+menuW > screenW {
		x = screenW - menuW
	}

	if x < 0 {
		x = 0
	}

	if y+menuH > screenH {
		y = screenH - menuH
	}

	if y < 0 {
		y = 0
	}

	return x, y
}

// OverlayMenuOnBase はベース画面上にメニューをアンカー位置に合成する。
func OverlayMenuOnBase(menu OverlayMenu, base string, anchorX, anchorY, screenW, screenH int) string {
	menuLines := menu.View()
	if len(menuLines) == 0 {
		return base
	}

	ox, oy := ClampMenuOrigin(menu.Width(), menu.Height(), anchorX, anchorY, screenW, screenH)
	baseLines := strings.Split(base, "\n")

	for len(baseLines) < screenH {
		baseLines = append(baseLines, "")
	}

	OverlayLines(baseLines, menuLines, ox, oy)

	return strings.Join(baseLines[:screenH], "\n")
}

// OverlayGeometry はオーバーレイの画面上の配置情報。
type OverlayGeometry struct {
	StartX, StartY     int
	ContentX, ContentY int
	OverlayW, OverlayH int
}

// Contains は座標がオーバーレイ内にあるかを判定する。
func (g OverlayGeometry) Contains(x, y int) bool {
	return x >= g.StartX && x < g.StartX+g.OverlayW &&
		y >= g.StartY && y < g.StartY+g.OverlayH
}

// CalcOverlayGeometry は実際のレンダリング結果から中央配置の座標を計算する。
func CalcOverlayGeometry(rendered string, screenW, screenH, borderW, padLeft, padTop int) OverlayGeometry {
	lines := strings.Split(rendered, "\n")
	overlayW := lipgloss.Width(rendered)
	overlayH := len(lines)

	startX := (screenW - overlayW) / 2
	startY := (screenH - overlayH) / 2

	if startX < 0 {
		startX = 0
	}

	if startY < 0 {
		startY = 0
	}

	return OverlayGeometry{
		StartX:   startX,
		StartY:   startY,
		ContentX: startX + borderW + padLeft,
		ContentY: startY + borderW + padTop,
		OverlayW: overlayW,
		OverlayH: overlayH,
	}
}

// OverlayCentered はベース画面上にオーバーレイを中央配置で合成する。
func OverlayCentered(base, overlay string, width, height int) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for len(baseLines) < height {
		baseLines = append(baseLines, "")
	}

	overlayW := lipgloss.Width(overlay)

	startY := (height - len(overlayLines)) / 2
	startX := (width - overlayW) / 2

	if startY < 0 {
		startY = 0
	}

	if startX < 0 {
		startX = 0
	}

	OverlayLines(baseLines, overlayLines, startX, startY)

	return strings.Join(baseLines[:height], "\n")
}

// OverlayLines は bodyLines の [startY, startY+len(overlayLines)) 行に対して
// startX の位置からオーバーレイを合成する。bodyLines を直接書き換える。
func OverlayLines(bodyLines []string, overlayLines []string, startX, startY int) {
	if len(overlayLines) == 0 {
		return
	}

	overlayW := lipgloss.Width(overlayLines[0])

	for i, line := range overlayLines {
		y := startY + i
		if y < 0 || y >= len(bodyLines) {
			continue
		}

		bodyLines[y] = ComposeLine(bodyLines[y], line, startX, overlayW)
	}
}
