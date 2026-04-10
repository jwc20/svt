package svt

import tea "charm.land/bubbletea/v2"

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
