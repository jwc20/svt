package svt

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"

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

func (m GameModel) Update(msg tea.Msg) (GameModel, tea.Cmd) {
	return m, nil
}

func (m GameModel) View() tea.View {
	return tea.NewView("GameView")
}
