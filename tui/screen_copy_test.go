package tui_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"

	"github.com/bluegreenhq/dogubako/tui"
)

// mockModel はテスト用の tea.Model。
type mockModel struct {
	content      string
	lastMsg      tea.Msg
	updateCalled bool
}

func (m *mockModel) Init() tea.Cmd { return nil }

func (m *mockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.lastMsg = msg
	m.updateCalled = true

	return m, nil
}

func (m *mockModel) View() tea.View {
	return tea.NewView(m.content)
}

func TestWrapScreenCopy_DelegatesView(t *testing.T) {
	t.Parallel()

	inner := &mockModel{content: "hello world"}
	wrapped := tui.WrapScreenCopy(inner)

	v := wrapped.View()
	assert.Contains(t, v.Content, "hello world")
}

func TestWrapScreenCopy_DelegatesUpdate(t *testing.T) {
	t.Parallel()

	inner := &mockModel{content: "hello"}
	wrapped := tui.WrapScreenCopy(inner)

	_, _ = wrapped.Update(tea.KeyPressMsg{Code: 'a'})

	assert.True(t, inner.updateCalled)
	assert.Equal(t, tea.KeyPressMsg{Code: 'a'}, inner.lastMsg)
}

func TestWrapScreenCopy_DelegatesInit(t *testing.T) {
	t.Parallel()

	inner := &mockModel{content: "hello"}
	wrapped := tui.WrapScreenCopy(inner)

	cmd := wrapped.Init()
	assert.Nil(t, cmd)
}

func TestWrapScreenCopy_CtrlShift4_SetsFlash(t *testing.T) {
	t.Parallel()

	inner := &mockModel{content: "screen content"}
	wrapped := tui.WrapScreenCopy(inner)

	// View を呼んで lastRendered をキャッシュ
	wrapped.View()

	// Ctrl+Shift+4 を送信
	_, cmd := wrapped.Update(tea.KeyPressMsg{Code: '4', Mod: tea.ModCtrl | tea.ModShift})

	// フラッシュメッセージがオーバーレイされる
	v := wrapped.View()
	// コピー成功時は "Screen copied" が含まれる
	// CI環境でクリップボードが使えない場合は "Failed" が含まれる
	assert.True(t,
		strings.Contains(v.Content, "Screen copied") || strings.Contains(v.Content, "Failed"),
		"flash message should be shown, got: %s", v.Content,
	)

	// Tick コマンドが返される（フラッシュ消去用）
	assert.NotNil(t, cmd)
}

func TestWrapScreenCopy_CtrlShift4_DoesNotDelegateToInner(t *testing.T) {
	t.Parallel()

	inner := &mockModel{content: "screen content"}
	wrapped := tui.WrapScreenCopy(inner)

	wrapped.View()
	_, _ = wrapped.Update(tea.KeyPressMsg{Code: '4', Mod: tea.ModCtrl | tea.ModShift})

	// inner には委譲されない
	assert.False(t, inner.updateCalled)
}

func TestWrapScreenCopy_ClearFlashMsg(t *testing.T) {
	t.Parallel()

	inner := &mockModel{content: "screen content"}
	wrapped := tui.WrapScreenCopy(inner)

	wrapped.View()
	_, cmd := wrapped.Update(tea.KeyPressMsg{Code: '4', Mod: tea.ModCtrl | tea.ModShift})

	// フラッシュがある状態を確認
	v := wrapped.View()
	hasFlash := strings.Contains(v.Content, "Screen copied") || strings.Contains(v.Content, "Failed")
	assert.True(t, hasFlash)

	// Tick コマンドを実行して clearFlashMsg を取得
	assert.NotNil(t, cmd)
	clearMsg := cmd()

	// clearFlashMsg を送信
	_, _ = wrapped.Update(clearMsg)

	// フラッシュが消えている
	v2 := wrapped.View()
	assert.NotContains(t, v2.Content, "Screen copied")
	assert.NotContains(t, v2.Content, "Failed")
}
