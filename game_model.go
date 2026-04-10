package svt

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"

	tea "charm.land/bubbletea/v2"
)

type GameModel struct {
	state         GameState
	phase         GamePhase
	store         GameStore
	purchaseSpent int

	promptTitle string
	promptLines []string

	choiceLog []string
	choiceVP  viewport.Model
	input     textinput.Model
	width     int
	height    int

	gameOver     bool
	gameResult   string
	deathMessage string
}

func NewGameModel(store GameStore, w, h int) GameModel {
	ti := textinput.New()
	ti.Placeholder = "Enter choice..."
	ti.CharLimit = 20
	ti.SetWidth(maxInt(w-8, 20))
	ti.Focus()

	vp := viewport.New()
	vp.SetWidth(maxInt(w-8, 20))
	vp.SetHeight(maxInt(h-20, 4))

	gs := InitState()

	m := GameModel{
		state:    gs,
		phase:    PhaseShooting,
		store:    store,
		input:    ti,
		choiceVP: vp,
		width:    w,
		height:   h,
	}

	//m.setShootingPrompt()
	return m
}

func (m GameModel) Init() tea.Cmd { return nil }

func (m GameModel) Update(msg tea.Msg) (GameModel, tea.Cmd) {
	return m, nil
}

func (m GameModel) View() tea.View {
	return tea.NewView("GameView")
}
