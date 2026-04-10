package svt

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var gamePanel = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62")).
	Padding(1, 2)

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

	m.setShootingPrompt()
	return m
}

func (m GameModel) Init() tea.Cmd { return nil }

func (m GameModel) Update(msg tea.Msg) (GameModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.input.SetWidth(maxInt(m.width-8, 20))
		m.choiceVP.SetWidth(maxInt(m.width-8, 20))
		m.choiceVP.SetHeight(maxInt(m.height-20, 4))
		m.refreshChoiceVP()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.gameOver {
				return m, func() tea.Msg { return BackToLobbyMsg{} }
			}
			return m, nil
		case "enter":
			return m.handleInput()
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}
func (m GameModel) handleInput() (GameModel, tea.Cmd) {
	val := strings.TrimSpace(m.input.Value())
	m.input.Reset()
	if val == "" {
		return m, nil
	}
	switch m.phase {
	case PhaseShooting:
		return m.handleShooting(val)
	case PhasePurchaseOxen, PhasePurchaseFood, PhasePurchaseAmmo,
		PhasePurchaseClothing, PhasePurchaseMisc:
		return m.handlePurchase(val)
	case PhaseTurnAction:
		return m.handleTurnAction(val)
	case PhaseEating:
		return m.handleEating(val)
	case PhaseGameOver:
		return m, func() tea.Msg { return BackToLobbyMsg{} }
	}
	return m, nil
}

func (m GameModel) handleShooting(val string) (GameModel, tea.Cmd) {
	level, err := strconv.Atoi(val)
	if err != nil || !SetShootingLevel(&m.state, level) {
		m.addChoice("✗ Invalid — enter 1-5")
		return m, nil
	}
	labels := []string{"", "Ace Marksman", "Good Shot", "Fair to Middlin'",
		"Need More Practice", "Shaky Knees"}
	m.addChoice(fmt.Sprintf("Shooting: (%d) %s", level, labels[level]))
	m.phase = PhasePurchaseOxen
	m.purchaseSpent = 0
	m.setOxenPrompt()
	return m, nil
}

func (m GameModel) handlePurchase(val string) (GameModel, tea.Cmd) {
	amount, err := strconv.Atoi(val)
	if err != nil {
		m.addChoice("✗ Enter a number")
		return m, nil
	}
	ok, errMsg := PurchaseItem(&m.state, m.phase, amount)
	if !ok {
		m.addChoice("✗ " + errMsg)
		return m, nil
	}
	m.purchaseSpent += amount
	switch m.phase {
	case PhasePurchaseOxen:
		m.addChoice(fmt.Sprintf("Bought $%d Oxen", amount))
		m.phase = PhasePurchaseFood
		m.setFoodPrompt()
	case PhasePurchaseFood:
		m.addChoice(fmt.Sprintf("Bought $%d Food", amount))
		m.phase = PhasePurchaseAmmo
		m.setAmmoPrompt()
	case PhasePurchaseAmmo:
		m.addChoice(fmt.Sprintf("Bought $%d Ammo", amount))
		m.phase = PhasePurchaseClothing
		m.setClothingPrompt()
	case PhasePurchaseClothing:
		m.addChoice(fmt.Sprintf("Bought $%d Clothing", amount))
		m.phase = PhasePurchaseMisc
		m.setMiscPrompt()
	case PhasePurchaseMisc:
		m.addChoice(fmt.Sprintf("Bought $%d Misc", amount))
		ok, remaining := FinalizePurchases(&m.state)
		if !ok {
			m.addChoice("✗ OVERSPENT — GAME OVER")
			m.setGameOver("died", "Overspent on supplies")
			return m, nil
		}
		m.addChoice(fmt.Sprintf("Cash left: $%d", remaining))
		return m.startTurn()
	}
	return m, nil
}

func (m GameModel) handleTurnAction(val string) (GameModel, tea.Cmd) {
	choice, err := strconv.Atoi(val)
	if err != nil {
		choice = 1
	}
	m.state.Trip.ActionChoice = choice
	if choice == 2 {
		m.addChoice("Chose: Hunt")
	} else {
		m.addChoice("Chose: Continue on trail")
	}
	m.phase = PhaseEating
	m.setEatingPrompt()
	return m, nil
}

func (m GameModel) handleEating(val string) (GameModel, tea.Cmd) {
	choice, err := strconv.Atoi(val)
	if err != nil || choice < 1 || choice > 3 {
		choice = 2
	}
	labels := []string{"", "Poorly", "Moderately", "Well"}
	m.addChoice(fmt.Sprintf("Eating: (%d) %s", choice, labels[choice]))
	ApplyEating(&m.state, choice)
	AdvanceMileage(&m.state)
	event := GenerateEvent(&m.state)
	if event != "" {
		m.addChoice("⚡ " + event)
	}
	if NeedsAilmentCheck(&m.state) {
		survived, msg := HandleAilment(&m.state)
		m.addChoice(msg)
		if !survived {
			m.setGameOver("died", msg)
			return m, nil
		}
	}
	return m.startTurn()
}

func (m GameModel) startTurn() (GameModel, tea.Cmd) {
	m.state.Trip.TurnNumber++
	m.state.Trip.CurrentDate = m.state.Trip.TurnNumber
	if IsStarved(&m.state) {
		m.addChoice("STARVED TO DEATH")
		m.setGameOver("starved", "Starved to death")
		return m, nil
	}
	if IsArrived(&m.state) {
		m.addChoice("ARRIVED IN OREGON!")
		m.setGameOver("won", "")
		return m, nil
	}
	m.phase = PhaseTurnAction
	m.setTurnPrompt()
	return m, nil
}

func (m GameModel) View() tea.View {
	var status strings.Builder
	turn := m.state.Trip.TurnNumber
	if turn < 1 {
		status.WriteString("Date: Preparing...   Mileage: 0 / " + fmt.Sprint(TotalRequiredMileage))
	} else {
		status.WriteString(fmt.Sprintf("Date: %s   Mileage: %d / %d",
			DateName(turn), m.state.Trip.Mileage, TotalRequiredMileage))
	}

	inv := m.state.Inventory
	inventory := fmt.Sprintf("Oxen: %d  Food: %d  Ammo: %d  Clothing: %d  Misc: %d  Cash: $%d",
		inv.Oxen, inv.Food, inv.Ammo, inv.Clothing, inv.Miscellaneous, m.state.Player.Cash)
	if m.state.Flags.Injured {
		inventory += "  " + BadStyle.Render("INJURED")
	}
	if m.state.Flags.Ill {
		inventory += "  " + BadStyle.Render("ILL")
	}

	prompt := FocusLabel.Render(m.promptTitle) + "\n\n" +
		strings.Join(m.promptLines, "\n")

	var inputLine string
	if !m.gameOver {
		inputLine = "\n" + PromptStyle.Render("enter: ") + m.input.View()
	}

	inner := status.String() + "\n" +
		inventory + "\n\n" +
		prompt +
		inputLine

	panel := gamePanel.Width(maxInt(m.width-2, 20)).Render(inner)
	help := DimStyle.Render("ctrl+c: quit")
	content := lipgloss.JoinVertical(lipgloss.Left, panel, help)

	return tea.NewView(content)
}

func (m *GameModel) addChoice(text string) {
	m.choiceLog = append(m.choiceLog, text)
	m.refreshChoiceVP()
}

func (m *GameModel) refreshChoiceVP() {
	m.choiceVP.SetContent(strings.Join(m.choiceLog, "\n"))
	m.choiceVP.GotoBottom()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m *GameModel) setShootingPrompt() {
	m.promptTitle = "HOW GOOD A SHOT ARE YOU WITH YOUR RIFLE?"
	m.promptLines = []string{
		"",
		"(1) ACE MARKSMAN",
		"(2) GOOD SHOT",
		"(3) FAIR TO MIDDLIN'",
		"(4) NEED MORE PRACTICE",
		"(5) SHAKY KNEES",
		"",
		PromptStyle.Render("Enter 1-5:"),
	}
}

func (m *GameModel) setOxenPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("You have $%d to spend on your trip.", InitialCash),
		"",
		"HOW MUCH DO YOU WANT TO SPEND ON",
		"YOUR OXEN TEAM?",
		"",
		DimStyle.Render("(Amount must be $200 – $300)"),
		"",
		PromptStyle.Render("Enter amount:"),
	}
}

func (m *GameModel) setFoodPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("Remaining budget: $%d", InitialCash-m.purchaseSpent),
		"", "HOW MUCH DO YOU WANT TO SPEND ON FOOD?",
		"", DimStyle.Render("(Amount must be $100 – $200)"),
		"", PromptStyle.Render("Enter amount:"),
	}
}

func (m *GameModel) setAmmoPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("Remaining budget: $%d", InitialCash-m.purchaseSpent),
		"", "HOW MUCH DO YOU WANT TO SPEND ON AMMO?",
		"", DimStyle.Render("(Amount must be $50 – $100)"),
		"", PromptStyle.Render("Enter amount:"),
	}
}

func (m *GameModel) setClothingPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("Remaining budget: $%d", InitialCash-m.purchaseSpent),
		"", "HOW MUCH DO YOU WANT TO SPEND ON CLOTHING?",
		"", DimStyle.Render("(Amount must be $50 – $100)"),
		"", PromptStyle.Render("Enter amount:"),
	}
}

func (m *GameModel) setMiscPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("Remaining budget: $%d", InitialCash-m.purchaseSpent),
		"", "HOW MUCH DO YOU WANT TO SPEND ON",
		"MISCELLANEOUS SUPPLIES?",
		"", DimStyle.Render("(Amount must be $50 – $100)"),
		"", PromptStyle.Render("Enter amount:"),
	}
}

func (m *GameModel) setEatingPrompt() {
	m.promptTitle = "EATING"
	m.promptLines = []string{
		"DO YOU WANT TO EAT:", "",
		"(1) POORLY", "(2) MODERATELY", "(3) WELL",
		"", PromptStyle.Render("Enter 1-3:"),
	}
}

func (m *GameModel) setTurnPrompt() {
	t := m.state.Trip.TurnNumber
	m.promptTitle = fmt.Sprintf("TURN %d — %s", t, DateName(t))
	m.promptLines = []string{
		"WHAT DO YOU WANT TO DO?", "",
		"(1) CONTINUE ON TRAIL", "(2) HUNT FOR FOOD",
		"", PromptStyle.Render("Enter 1-2:"),
	}
}

func (m *GameModel) setGameOver(result, deathMsg string) {
	m.phase = PhaseGameOver
	m.gameOver = true
	m.gameResult = result
	m.deathMessage = deathMsg

	switch result {
	case "won":
		m.promptTitle = "★ CONGRATULATIONS! ★"
		m.promptLines = []string{
			"", GoodStyle.Render("YOU MADE IT TO OREGON CITY!"), "",
			fmt.Sprintf("Turns taken: %d", m.state.Trip.TurnNumber),
			fmt.Sprintf("Cash remaining: $%d", m.state.Player.Cash),
			"", DimStyle.Render("Press enter or esc to return to lobby."),
		}
	case "starved":
		m.promptTitle = "GAME OVER"
		m.promptLines = []string{
			"", BadStyle.Render("YOU RAN OUT OF FOOD AND STARVED TO DEATH."), "",
			fmt.Sprintf("Mileage reached: %d / %d", m.state.Trip.Mileage, TotalRequiredMileage),
			"", DimStyle.Render("Press enter or esc to return to lobby."),
		}
	default:
		m.promptTitle = "GAME OVER"
		m.promptLines = []string{
			"", BadStyle.Render(deathMsg), "",
			fmt.Sprintf("Mileage reached: %d / %d", m.state.Trip.Mileage, TotalRequiredMileage),
			"", DimStyle.Render("Press enter or esc to return to lobby."),
		}
	}
}
