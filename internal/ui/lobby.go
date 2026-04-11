package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type LobbyModel struct {
	cursor        int
	choices       []string
	width, height int
	playerId      string
}

func NewLobbyModel(playerId string) LobbyModel {
	return LobbyModel{
		choices:  []string{"Play", "Leaderboard (TODO)"},
		playerId: playerId,
	}
}

func (m LobbyModel) Init() tea.Cmd { return nil }

func (m LobbyModel) Update(msg tea.Msg) (LobbyModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == 0 {
				return m, func() tea.Msg { return StartGameMsg{} }
			}
		}
	}
	return m, nil
}

func (m LobbyModel) View() tea.View {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("   THE SILICON TRAIL") + "\n\n")
	sb.WriteString(DimStyle.Render(fmt.Sprintf("Welcome, %s", m.playerId)) + "\n\n")
	for i, choice := range m.choices {
		cursor := "  "
		if i == m.cursor {
			cursor = SelectedStyle.Render("▸ ")
		}
		label := PlainLabel.Render(choice)
		if i == m.cursor {
			label = SelectedStyle.Render(choice)
		}
		sb.WriteString(fmt.Sprintf("%s%s\n", cursor, label))
	}
	sb.WriteString("\n" + DimStyle.Render("↑/↓: navigate  enter: select  ctrl+c: quit"))
	content := lipgloss.NewStyle().Align(lipgloss.Center).Render(sb.String())

	v := tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content))
	v.AltScreen = true
	return v
}
