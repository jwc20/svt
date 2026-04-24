package ui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/jwc20/svt/internal/engine"
)

type ViewState int

const (
	SplashView ViewState = iota
	LobbyView
	GameView
	LeaderboardView
)

type RootModel struct {
	state         ViewState
	splash        SplashModel
	lobby         LobbyModel
	game          *GameModel
	leaderboard   *LeaderboardModel
	store         engine.GameStore
	playerId      int64
	userName      string
	bonusHype     int
	width, height int
}

func NewRootModel(store engine.GameStore, playerID int64, userName string, bonusHype int) RootModel {
	return RootModel{
		state: SplashView,
		//lobby:     NewLobbyModel(userName),
		store:     store,
		playerId:  playerID,
		userName:  userName,
		bonusHype: bonusHype,
	}
}

func (m RootModel) Init() tea.Cmd { return nil }

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	//	m.lobby.width, m.lobby.height = msg.Width, msg.Height
	//	m.splash.width, m.splash.height = msg.Width, msg.Height
	//	if m.game != nil {
	//		m.game.Resize(msg.Width, msg.Height)
	//	}
	//	if m.leaderboard != nil {
	//		m.leaderboard.width, m.leaderboard.height = msg.Width, msg.Height
	//	}
	//	return m, nil
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case StartGameMsg:
		gm := NewGameModel(m.store, m.playerId, m.bonusHype, m.width, m.height)
		m.game = &gm
		m.state = GameView
		return m, m.game.Init()

	case ShowLeaderboardMsg:
		entries, err := m.store.Leaderboard(10)
		if err != nil {
			return m, nil
		}
		lb := NewLeaderboardModel(entries, m.width, m.height)
		m.leaderboard = &lb
		m.state = LeaderboardView
		return m, nil

	case BackToLobbyMsg:
		m.game = nil
		m.leaderboard = nil
		l := NewLobbyModel(m.userName, m.width, m.height)
		m.lobby = l
		m.state = LobbyView
		return m, nil
	}

	switch m.state {
	case SplashView:
		var cmd tea.Cmd
		m.splash, cmd = m.splash.Update(msg)
		return m, cmd
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
	case LeaderboardView:
		if m.leaderboard != nil {
			var cmd tea.Cmd
			*m.leaderboard, cmd = m.leaderboard.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m RootModel) View() tea.View {
	switch m.state {
	case SplashView:
		return m.splash.View()
	case GameView:
		if m.game != nil {
			return m.game.View()
		}
	case LeaderboardView:
		if m.leaderboard != nil {
			return m.leaderboard.View()
		}
	}
	return m.lobby.View()
}
