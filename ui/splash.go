package ui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SplashModel struct {
	width, height int
}

type tickMsg time.Time

func tick() tea.Msg {
	return tickMsg(time.Now())
}

func (m SplashModel) Init() tea.Cmd { return nil }

func (m SplashModel) Update(msg tea.Msg) (SplashModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m, func() tea.Msg { return BackToLobbyMsg{} }
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		return m, tick
	}
	return m, nil
}

var style = lipgloss.NewStyle()

func (m SplashModel) View() tea.View {
	var v tea.View
	v.AltScreen = true

	content := style.Render(
		style.Render(
			style.Render("Are you sure you want to eat that "),
		),
	)

	v.SetContent(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content))

	return v
}
