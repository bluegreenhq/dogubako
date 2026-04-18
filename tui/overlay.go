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
