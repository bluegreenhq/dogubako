package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ConfirmResult はダイアログの操作結果を表す。
type ConfirmResult int

const (
	// ConfirmContinue はダイアログ継続中。
	ConfirmContinue ConfirmResult = iota
	// ConfirmYes は確定。
	ConfirmYes
	// ConfirmNo はキャンセル。
	ConfirmNo
)

const (
	confirmDialogPaddingH = 2
	confirmDialogWidth    = 40
)

// ConfirmDialog は汎用確認ダイアログコンポーネント。
type ConfirmDialog struct {
	Title       string
	Detail      string
	yesBtn      BoxButton
	noBtn       BoxButton
	screenWidth int // 画面幅
	bodyHeight  int // ボディ領域の高さ
}

// NewConfirmDialog は新しい ConfirmDialog を生成する。
func NewConfirmDialog(title, detail string) ConfirmDialog {
	return ConfirmDialog{
		Title:       title,
		Detail:      detail,
		yesBtn:      NewBoxButton("[Y]es"),
		noBtn:       NewBoxButton("[N]o"),
		screenWidth: 0,
		bodyHeight:  0,
	}
}

const (
	confirmContentLines  = 7 // title + empty + detail + empty + button top + button middle + button bottom
	confirmBorderPad     = 3 // border上1 + padding上1 + border下1
	confirmButtonGapCols = 2 // ボタン間のスペース数
	confirmButtonRows    = 3 // ボタン描画行数（上枠・ラベル・下枠）
	confirmPaddingSides  = 2 // 左右パディングの個数
	confirmCenterDiv     = 2 // センタリング用除数
)

// SetScreenSize は画面サイズを設定する。
func (d *ConfirmDialog) SetScreenSize(screenWidth, bodyHeight int) {
	d.screenWidth = screenWidth
	d.bodyHeight = bodyHeight
}

// Origin はダイアログの画面上のコンテンツ左上座標を返す。
// border + padding 分を加算した座標。
func (d *ConfirmDialog) Origin() (int, int) {
	rendered := d.View()
	dialogLines := strings.Split(rendered, "\n")

	const (
		borderPaddingLines = 2 // border上 + padding上
		borderPaddingCols  = 3 // border左1 + padding左2
	)

	dialogWidth := lipgloss.Width(dialogLines[0])

	sy := max((d.bodyHeight-len(dialogLines))/confirmCenterDiv, 0) + borderPaddingLines
	sx := max((d.screenWidth-dialogWidth)/confirmCenterDiv, 0) + borderPaddingCols

	return sx, sy
}

// View はダイアログの描画内容を返す。
func (d *ConfirmDialog) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true)
	b.WriteString(titleStyle.Render(d.Title))
	b.WriteString("\n\n")
	b.WriteString(d.Detail)
	b.WriteString("\n\n")

	gap := strings.Repeat(" ", confirmButtonGapCols)
	pad := strings.Repeat(" ", d.buttonPadLeft())
	b.WriteString(pad + d.yesBtn.ViewTop() + gap + d.noBtn.ViewTop())
	b.WriteString("\n")
	b.WriteString(pad + d.yesBtn.ViewMiddle() + gap + d.noBtn.ViewMiddle())
	b.WriteString("\n")
	b.WriteString(pad + d.yesBtn.ViewBottom() + gap + d.noBtn.ViewBottom())

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("9")).
		PaddingTop(1).
		PaddingBottom(0).
		PaddingLeft(confirmDialogPaddingH).
		PaddingRight(confirmDialogPaddingH).
		Width(confirmDialogWidth)

	return style.Render(b.String())
}

// DialogLines はレンダリング後の行数を返す。
func (d *ConfirmDialog) DialogLines() int {
	// padding上1 + border上1 + content + padding下1 + border下1
	return d.contentLines() + confirmBorderPad
}

// Update はキー入力に応じてダイアログの状態を更新する。
func (d *ConfirmDialog) Update(msg tea.Msg) ConfirmResult {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return ConfirmContinue
	}

	switch keyMsg.Code {
	case 'y', 'Y', tea.KeyEnter:
		return ConfirmYes
	case 'n', 'N', tea.KeyEscape:
		return ConfirmNo
	}

	return ConfirmContinue
}

// HandleClick は相対座標でのクリックを処理する。
// relY はダイアログコンテンツ内の行番号。
func (d *ConfirmDialog) HandleClick(relX, relY int) ConfirmResult {
	buttonTop := d.contentLines() - confirmButtonRows // ボタン3行の先頭行

	if relY < buttonTop || relY > buttonTop+2 {
		return ConfirmContinue
	}

	padLeft := d.buttonPadLeft()
	noStartX := padLeft + d.yesBtn.DisplayWidth() + confirmButtonGapCols

	if d.yesBtn.HitTest(relX, padLeft) {
		return ConfirmYes
	}

	if d.noBtn.HitTest(relX, noStartX) {
		return ConfirmNo
	}

	return ConfirmContinue
}

// HandleMotion は相対座標でのマウスホバーを処理する。
func (d *ConfirmDialog) HandleMotion(relX, relY int) {
	d.yesBtn.SetHovered(false)
	d.noBtn.SetHovered(false)

	buttonTop := d.contentLines() - confirmButtonRows

	if relY < buttonTop || relY > buttonTop+2 {
		return
	}

	padLeft := d.buttonPadLeft()
	noStartX := padLeft + d.yesBtn.DisplayWidth() + confirmButtonGapCols

	if d.yesBtn.HitTest(relX, padLeft) {
		d.yesBtn.SetHovered(true)
	} else if d.noBtn.HitTest(relX, noStartX) {
		d.noBtn.SetHovered(true)
	}
}

// ClearHover はホバー状態をリセットする。
func (d *ConfirmDialog) ClearHover() {
	d.yesBtn.SetHovered(false)
	d.noBtn.SetHovered(false)
}

// HandleClickAbs は画面絶対座標でのクリックを処理する。
func (d *ConfirmDialog) HandleClickAbs(absX, absY int) ConfirmResult {
	originX, originY := d.Origin()

	return d.HandleClick(absX-originX, absY-originY)
}

// HandleMotionAbs は画面絶対座標でのマウスホバーを処理する。
func (d *ConfirmDialog) HandleMotionAbs(absX, absY int) {
	originX, originY := d.Origin()
	d.HandleMotion(absX-originX, absY-originY)
}

// buttonPadLeft はボタン領域のセンタリング用左パディングを返す。
func (d *ConfirmDialog) buttonPadLeft() int {
	buttonsWidth := d.yesBtn.DisplayWidth() + confirmButtonGapCols + d.noBtn.DisplayWidth()
	contentWidth := confirmDialogWidth - confirmDialogPaddingH*confirmPaddingSides

	return (contentWidth - buttonsWidth) / confirmCenterDiv
}

// contentLines はダイアログコンテンツ（padding内側）の行数を返す。
func (d *ConfirmDialog) contentLines() int {
	return confirmContentLines
}
