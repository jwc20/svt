package ui

import (
	"image/color"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SplashModel struct {
	width, height int
	rate          float64
	startTime     time.Time
	cursor        tea.Cursor
	blink         bool
}

func (m SplashModel) Init() tea.Cmd {
	return tick
}

func (m SplashModel) Update(msg tea.Msg) (SplashModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m, func() tea.Msg { return BackToLobbyMsg{} }
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tickMsg:
		if m.startTime.IsZero() {
			m.startTime = time.Now()
		}
		return m, tick
	}
	return m, nil
}

func (m SplashModel) View() tea.View {
	var v tea.View
	v.AltScreen = true

	content := "svt"

	v.SetContent(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content))
	c := tea.NewCursor(m.width/2+1, m.height/2-1)
	c.Shape = m.cursor.Shape
	c.Blink = m.blink
	c.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	v.Cursor = c

	return v
}

type tickMsg time.Time

func tick() tea.Msg {
	return tickMsg(time.Now())
}
