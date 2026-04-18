package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// boxButtonPadding は罫線ボタンの追加幅（│ + space + space + │）。
const boxButtonPadding = 4

var (
	boxBorderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	boxLabelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	boxLabelHover  = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("4"))
)

// BoxButton は罫線で囲まれたボタンコンポーネント。
type BoxButton struct {
	label   string
	hovered bool
}

// NewBoxButton は新しい BoxButton を生成する。
func NewBoxButton(label string) BoxButton {
	return BoxButton{label: label, hovered: false}
}

// Label はボタンのラベルを返す。
func (b *BoxButton) Label() string { return b.label }

// Hovered はホバー状態を返す。
func (b *BoxButton) Hovered() bool { return b.hovered }

// SetHovered はホバー状態を設定する。
func (b *BoxButton) SetHovered(h bool) { b.hovered = h }

// DisplayWidth はボタンの表示幅（│ + space + label + space + │）を返す。
func (b *BoxButton) DisplayWidth() int {
	return lipgloss.Width(b.label) + boxButtonPadding
}

// HitTest は startX からの相対位置でボタン範囲内かを判定する。
func (b *BoxButton) HitTest(x, startX int) bool {
	return x >= startX && x < startX+b.DisplayWidth()
}

// ViewTop は上罫線行（┌────┐）を返す。
func (b *BoxButton) ViewTop() string {
	return boxBorderStyle.Render("┌" + strings.Repeat("─", lipgloss.Width(b.innerLabel())) + "┐")
}

// ViewMiddle はラベル行（│ label │）を返す。
func (b *BoxButton) ViewMiddle() string {
	inner := b.innerLabel()

	var labelStyled string
	if b.hovered {
		labelStyled = boxLabelHover.Render(inner)
	} else {
		labelStyled = boxLabelStyle.Render(inner)
	}

	return boxBorderStyle.Render("│") + labelStyled + boxBorderStyle.Render("│")
}

// ViewBottom は下罫線行（└────┘）を返す。
func (b *BoxButton) ViewBottom() string {
	return boxBorderStyle.Render("└" + strings.Repeat("─", lipgloss.Width(b.innerLabel())) + "┘")
}

// innerLabel は罫線内の文字列（スペース含む）を返す。
func (b *BoxButton) innerLabel() string {
	return " " + b.label + " "
}
