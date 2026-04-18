package tui_test

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"

	"github.com/bluegreenhq/dogubako/tui"
)

func TestBoxButtonDisplayWidth_ASCII(t *testing.T) {
	t.Parallel()

	b := tui.NewBoxButton("OK")
	// "OK" は表示幅2 + padding 4 = 6
	assert.Equal(t, 6, b.DisplayWidth())
}

func TestBoxButtonDisplayWidth_Multibyte(t *testing.T) {
	t.Parallel()

	b := tui.NewBoxButton("キャンセル")
	// "キャンセル" は全角5文字 = 表示幅10 + padding 4 = 14
	// ※ len("キャンセル") = 15 (バイト長) なので、バイト長だと 19 になってしまう
	assert.Equal(t, 14, b.DisplayWidth())
}

func TestBoxButtonViewTop_Multibyte(t *testing.T) {
	t.Parallel()

	b := tui.NewBoxButton("キャンセル")
	top := b.ViewTop()
	// ViewTop の表示幅は DisplayWidth と一致すべき
	assert.Equal(t, b.DisplayWidth(), lipgloss.Width(top))
}

func TestBoxButtonViewBottom_Multibyte(t *testing.T) {
	t.Parallel()

	b := tui.NewBoxButton("キャンセル")
	bottom := b.ViewBottom()
	// ViewBottom の表示幅は DisplayWidth と一致すべき
	assert.Equal(t, b.DisplayWidth(), lipgloss.Width(bottom))
}
