package svt

import tea "charm.land/bubbletea/v2"

type ViewState int

const (
	LobbyView ViewState = iota
	GameView
)

type RootModel struct {
	state         ViewState
	lobby         LobbyModel
	game          *GameModel
	store         GameStore
	playerId      string
	width, height int
}

func NewRootModel(store GameStore, userID string) RootModel {
	return RootModel{
		state:    LobbyView,
		lobby:    NewLobbyModel(userID),
		store:    store,
		playerId: userID,
	}
}

func (m RootModel) Init() tea.Cmd { return nil }

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.lobby.width, m.lobby.height = msg.Width, msg.Height
		if m.game != nil {
			m.game.width, m.game.height = msg.Width, msg.Height
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.state {
	case LobbyView:
		var cmd tea.Cmd
		m.lobby, cmd = m.lobby.Update(msg)
		return m, cmd
	case GameView:
		if m.game != nil {
			var cmd tea.Cmd
			*m.game, cmd = m.game.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m RootModel) View() tea.View {
	var v tea.View
	if m.state == GameView && m.game != nil {
		v = m.game.View()
	} else {
		v = m.lobby.View()
	}
	v.AltScreen = true
	return v
}
