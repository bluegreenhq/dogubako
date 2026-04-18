package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// MenuItem はポップアップメニューの1項目を表す。
type MenuItem struct {
	Label    string
	Disabled bool
}

// NewMenuItem は有効な MenuItem を生成する。
func NewMenuItem(label string) MenuItem {
	return MenuItem{Label: label, Disabled: false}
}

// NewDisabledMenuItem は無効な MenuItem を生成する。
func NewDisabledMenuItem(label string) MenuItem {
	return MenuItem{Label: label, Disabled: true}
}

var (
	menuItemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	menuItemHoverStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("4"))
	menuItemDisabledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	menuBorderStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// PopupMenu は汎用ポップアップメニューコンポーネント。
type PopupMenu struct {
	items []MenuItem
	hover int // ホバー中の項目インデックス (-1 = なし)
}

// NewPopupMenu は新しい PopupMenu を生成する。
func NewPopupMenu(items []MenuItem) *PopupMenu {
	return &PopupMenu{items: items, hover: -1}
}

// Items はメニュー項目のリストを返す。
func (m *PopupMenu) Items() []MenuItem { return m.items }

const (
	menuBorderLines = 2 // 上枠 + 下枠
	menuPadding     = 4 // │ + space + space + │
)

// Height はメニュー全体の高さ（上枠 + 項目数*2-1(項目間空行) + 下枠）を返す。
func (m *PopupMenu) Height() int {
	n := len(m.items)
	if n == 0 {
		return 0
	}

	// 上枠1 + (項目数 + 項目間空行(n-1)) + 下枠1
	return menuBorderLines + n + n - 1
}

// Width はメニュー全体の幅（左枠 + スペース + 最長ラベル + スペース + 右枠）を返す。
func (m *PopupMenu) Width() int {
	maxLen := 0
	for _, item := range m.items {
		if w := lipgloss.Width(item.Label); w > maxLen {
			maxLen = w
		}
	}

	return maxLen + menuPadding
}

// Hover はホバー中の項目インデックスを返す。-1 はホバーなし。
func (m *PopupMenu) Hover() int { return m.hover }

// SetHover はホバー中の項目インデックスを設定する。
func (m *PopupMenu) SetHover(h int) { m.hover = h }

// View はポップアップメニューを枠付きで描画する。
// 各行は独立した文字列のスライスとして返す。
func (m *PopupMenu) View() []string {
	if len(m.items) == 0 {
		return nil
	}

	innerWidth := m.Width() - menuBorderLines // 左右の枠を除く

	lines := make([]string, 0, m.Height())

	// 上枠
	lines = append(lines, menuBorderStyle.Render("┌"+strings.Repeat("─", innerWidth)+"┐"))

	// 項目
	for i, item := range m.items {
		if i > 0 {
			// 項目間の空行（枠付き）
			emptyLine := menuBorderStyle.Render("│") +
				menuBorderStyle.Render(strings.Repeat(" ", innerWidth)) +
				menuBorderStyle.Render("│")
			lines = append(lines, emptyLine)
		}

		padRight := innerWidth - 1 - lipgloss.Width(item.Label)
		label := " " + item.Label + strings.Repeat(" ", padRight)

		var styled string

		switch {
		case item.Disabled:
			styled = menuItemDisabledStyle.Render(label)
		case i == m.hover:
			styled = menuItemHoverStyle.Render(label)
		default:
			styled = menuItemStyle.Render(label)
		}

		line := menuBorderStyle.Render("│") + styled + menuBorderStyle.Render("│")
		lines = append(lines, line)
	}

	// 下枠
	lines = append(lines, menuBorderStyle.Render("└"+strings.Repeat("─", innerWidth)+"┘"))

	return lines
}

// HitTest は座標からメニュー項目のインデックスを返す。
// メニュー左上を (0,0) とする相対座標。
// 項目間には空行があり、偶数行(1,3,5...)が項目、奇数行(2,4,...)が空行。
// 項目外、空行、または disabled の場合は -1 を返す。
func (m *PopupMenu) HitTest(x, y int) int {
	if x < 0 || x >= m.Width() {
		return -1
	}

	// y=0: 上枠, y=Height()-1: 下枠
	if y <= 0 || y >= m.Height()-1 {
		return -1
	}

	// 枠の内側の行番号 (1始まり)
	innerY := y - 1

	// 偶数番目(0,2,4...) が項目行、奇数番目(1,3,5...) が空行
	if innerY%menuBorderLines != 0 {
		return -1 // 空行
	}

	idx := innerY / menuBorderLines

	if idx < 0 || idx >= len(m.items) {
		return -1
	}

	if m.items[idx].Disabled {
		return -1
	}

	return idx
}

// MoveHoverDown はホバーを次の有効な項目に移動する。
func (m *PopupMenu) MoveHoverDown() {
	for i := m.hover + 1; i < len(m.items); i++ {
		if !m.items[i].Disabled {
			m.hover = i

			return
		}
	}
}

// MoveHoverUp はホバーを前の有効な項目に移動する。
func (m *PopupMenu) MoveHoverUp() {
	start := m.hover - 1
	if m.hover < 0 {
		start = len(m.items) - 1
	}

	for i := start; i >= 0; i-- {
		if !m.items[i].Disabled {
			m.hover = i

			return
		}
	}
}

// SelectHover はホバー中の項目インデックスを返す。ホバーなしなら -1。
func (m *PopupMenu) SelectHover() int {
	return m.hover
}

// HandleKeyNav はキー入力に応じてホバーを上下に移動する。
func (m *PopupMenu) HandleKeyNav(msg tea.KeyPressMsg) {
	switch {
	case msg.Code == tea.KeyDown || msg.Code == 'j' ||
		(msg.Code == 'n' && msg.Mod&tea.ModCtrl != 0):
		m.MoveHoverDown()
	case msg.Code == tea.KeyUp || msg.Code == 'k' ||
		(msg.Code == 'p' && msg.Mod&tea.ModCtrl != 0):
		m.MoveHoverUp()
	}
}

// SetHoverByPos はマウス座標からホバー状態を更新する。
// メニュー左上を (0,0) とする相対座標。
func (m *PopupMenu) SetHoverByPos(x, y int) {
	m.hover = m.HitTest(x, y)
}

// HandleClick はクリック座標から選択された項目インデックスを返す。
// メニュー左上を (0,0) とする相対座標。
// 戻り値は (インデックス, ヒットしたか)。
func (m *PopupMenu) HandleClick(x, y int) (int, bool) {
	idx := m.HitTest(x, y)
	if idx < 0 {
		return -1, false
	}

	return idx, true
}
