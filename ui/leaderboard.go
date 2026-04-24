package ui

import (
	"fmt"
	"strconv"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/jwc20/svt/internal/engine"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type LeaderboardModel struct {
	table         table.Model
	width, height int
}

func NewLeaderboardModel(entries []engine.LeaderboardEntry, width, height int) LeaderboardModel {
	columns := []table.Column{
		{Title: "Rank", Width: 6},
		{Title: "Player", Width: 14},
		{Title: "Score", Width: 8},
		{Title: "Ended At", Width: 20},
	}

	rows := make([]table.Row, len(entries))
	for i, e := range entries {
		name := e.Username
		if name == "" {
			name = "anonymous"
		}
		rows[i] = table.Row{
			strconv.Itoa(e.Rank),
			name,
			strconv.Itoa(e.Score),
			e.EndedAt.Format("2006-01-02 15:04:05"),
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(56),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return LeaderboardModel{
		table:  t,
		width:  width,
		height: height,
	}
}

func (m LeaderboardModel) Init() tea.Cmd { return nil }

func (m LeaderboardModel) Update(msg tea.Msg) (LeaderboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return BackToLobbyMsg{} }
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m LeaderboardModel) View() tea.View {
	title := TitleStyle.Render("   LEADERBOARD")
	content := fmt.Sprintf("%s\n%s\n  %s\n%s",
		title,
		baseStyle.Render(m.table.View()),
		m.table.HelpView(),
		DimStyle.Render("  esc/q: back"),
	)
	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	v := tea.NewView(centered)
	v.AltScreen = true
	return v
}
