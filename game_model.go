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

var (
	outerBorder = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Padding(0, 1).
			AlignVertical(lipgloss.Center).
			Padding(1, 2)

	promptPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	statusBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	logBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)
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
	ti.SetWidth(maxInt(w/3-6, 14))
	ti.Focus()

	vp := viewport.New()
	vp.SetWidth(maxInt(w/3-4, 14))
	vp.SetHeight(maxInt(h-22, 4))

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

func (m GameModel) rightW() int { return maxInt(m.width*30/100, 22) }
func (m GameModel) leftW() int  { return m.width - m.rightW() - 4 } // -4 for outer border

func (m GameModel) Update(msg tea.Msg) (GameModel, tea.Cmd) {
	switch msg := msg.(type) {
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

func (m *GameModel) Resize(w, h int) {
	m.width, m.height = w, h
	rw := m.rightW()
	m.input.SetWidth(maxInt(rw-6, 14))
	m.choiceVP.SetWidth(maxInt(rw-4, 14))
	m.choiceVP.SetHeight(maxInt(h-22, 4))
	m.refreshChoiceVP()
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

// ── phase handlers ────────────────────────────────────────────────

func (m GameModel) handleShooting(val string) (GameModel, tea.Cmd) {
	level, err := strconv.Atoi(val)
	if err != nil || !SetShootingLevel(&m.state, level) {
		m.addChoice("✗ Invalid — enter 1-5")
		return m, nil
	}
	labels := []string{"", "Ace Marksman", "Good Shot", "Fair to Middlin'", "Need More Practice", "Shaky Knees"}
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

// ── turn management ───────────────────────────────────────────────

func (m GameModel) startTurn() (GameModel, tea.Cmd) {
	m.state.Trip.TurnNumber++
	m.state.Trip.CurrentDate = m.state.Trip.TurnNumber
	if IsStarved(&m.state) {
		m.addChoice("✗ STARVED TO DEATH")
		m.setGameOver("starved", "Starved to death")
		return m, nil
	}
	if IsArrived(&m.state) {
		m.addChoice("★ ARRIVED IN OREGON!")
		m.setGameOver("won", "")
		return m, nil
	}
	m.phase = PhaseTurnAction
	m.setTurnPrompt()
	return m, nil
}

// ── prompt setters ────────────────────────────────────────────────

func (m *GameModel) setShootingPrompt() {
	m.promptTitle = "How good a shot are you?"
	m.promptLines = []string{
		"",
		"(1) Ace Marksman",
		"(2) Good Shot",
		"(3) Fair to Middlin'",
		"(4) Need More Practice",
		"(5) Shaky Knees",
	}
}

func (m *GameModel) setOxenPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("You have $%d to spend on your trip.", InitialCash),
		"",
		"How much do you want to spend on your oxen team?",
		DimStyle.Render("(Amount must be $200 – $300)"),
	}
}

func (m *GameModel) setFoodPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("Remaining budget: $%d", InitialCash-m.purchaseSpent),
		"", "How much do you want to spend on food?",
		DimStyle.Render("(Amount must be $100 – $200)"),
	}
}

func (m *GameModel) setAmmoPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("Remaining budget: $%d", InitialCash-m.purchaseSpent),
		"", "How much do you want to spend on ammo?",
		DimStyle.Render("(Amount must be $50 – $100)"),
	}
}

func (m *GameModel) setClothingPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("Remaining budget: $%d", InitialCash-m.purchaseSpent),
		"", "How much do you want to spend on clothing?",
		DimStyle.Render("(Amount must be $50 – $100)"),
	}
}

func (m *GameModel) setMiscPrompt() {
	m.promptTitle = "PURCHASE SUPPLIES"
	m.promptLines = []string{
		fmt.Sprintf("Remaining budget: $%d", InitialCash-m.purchaseSpent),
		"", "How much do you want to spend on miscellaneous supplies?",
		DimStyle.Render("(Amount must be $50 – $100)"),
	}
}

func (m *GameModel) setEatingPrompt() {
	m.promptTitle = "EATING"
	m.promptLines = []string{
		"Do you want to eat:",
		"",
		"(1) Poorly",
		"(2) Moderately",
		"(3) Well",
	}
}

func (m *GameModel) setTurnPrompt() {
	t := m.state.Trip.TurnNumber
	m.promptTitle = fmt.Sprintf("TURN %d — %s", t, DateName(t))
	m.promptLines = []string{
		"What do you want to do?",
		"",
		"(1) Continue on trail",
		"(2) Hunt for food",
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
		m.promptTitle = "✗ GAME OVER"
		m.promptLines = []string{
			"", BadStyle.Render("YOU RAN OUT OF FOOD AND STARVED TO DEATH."), "",
			fmt.Sprintf("Mileage reached: %d / %d", m.state.Trip.Mileage, TotalRequiredMileage),
			"", DimStyle.Render("Press enter or esc to return to lobby."),
		}
	default:
		m.promptTitle = "✗ GAME OVER"
		m.promptLines = []string{
			"", BadStyle.Render(deathMsg), "",
			fmt.Sprintf("Mileage reached: %d / %d", m.state.Trip.Mileage, TotalRequiredMileage),
			"", DimStyle.Render("Press enter or esc to return to lobby."),
		}
	}
}

// ── View ──────────────────────────────────────────────────────────

func (m GameModel) View() tea.View {
	lw := m.leftW()
	rw := m.rightW()
	innerH := m.height - 4 // outer border chrome

	// left: prompt panel, text centered in available space
	promptText := FocusLabel.Render(m.promptTitle) + "\n\n" +
		strings.Join(m.promptLines, "\n")
	leftContentH := maxInt(innerH-4, 6) // -4 for prompt border chrome
	centeredPrompt := lipgloss.Place(lw-6, leftContentH,
		lipgloss.Center, lipgloss.Center, promptText)
	leftCol := promptPanel.Width(lw - 2).Height(leftContentH).Render(centeredPrompt)

	// ── right column ──

	sBox := statusBox.Width(rw - 2).Render(m.renderStatus())
	statusRenderedH := lipgloss.Height(sBox)

	logPanelH := maxInt(innerH-statusRenderedH-2, 4) // -2 for log border chrome
	logInnerH := maxInt(logPanelH, 2)

	var logInner string
	if m.gameOver {
		// entire inner space is the viewport
		m.choiceVP.SetWidth(maxInt(rw-6, 10))
		m.choiceVP.SetHeight(logInnerH)
		logInner = m.choiceVP.View()
	} else {
		// viewport gets space minus 1 line for input
		inputLine := PromptStyle.Render("") + m.input.View()
		vpH := maxInt(logInnerH-2, 1) // -2: 1 for input, 1 for separator
		m.choiceVP.SetWidth(maxInt(rw-6, 10))
		m.choiceVP.SetHeight(vpH)
		separator := DimStyle.Render(strings.Repeat("─", maxInt(rw-6, 10)))
		logInner = m.choiceVP.View() + "\n" + separator + "\n" + inputLine
	}

	lBox := logBox.Width(rw - 2).Height(logInnerH).Render(logInner)

	rightCol := lipgloss.JoinVertical(lipgloss.Left, sBox, lBox)

	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)
	outer := outerBorder.Width(m.width - 2).Height(innerH).Render(columns)

	help := DimStyle.Render("ctrl+c: quit")
	content := lipgloss.JoinVertical(lipgloss.Left, outer, help)
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m GameModel) renderStatus() string {
	var sb strings.Builder

	turn := m.state.Trip.TurnNumber
	if turn < 1 {
		sb.WriteString(PlainLabel.Render("Date") + "\n")
		sb.WriteString("Preparing...\n")
	} else {
		sb.WriteString(PlainLabel.Render("Date") + "\n")
		sb.WriteString(DateName(turn) + "\n")
	}
	sb.WriteString(PlainLabel.Render("Mileage") + "\n")
	sb.WriteString(fmt.Sprintf("%d / %d\n", m.state.Trip.Mileage, TotalRequiredMileage))

	sb.WriteString("\n")
	sb.WriteString(PlainLabel.Render("Inventory") + "\n")

	inv := m.state.Inventory
	sb.WriteString(fmt.Sprintf("Oxen: %d\n", inv.Oxen))
	sb.WriteString(fmt.Sprintf("Food: %d\n", inv.Food))
	sb.WriteString(fmt.Sprintf("Ammo: %d\n", inv.Ammo))
	sb.WriteString(fmt.Sprintf("Clothing: %d\n", inv.Clothing))
	sb.WriteString(fmt.Sprintf("Misc: %d\n", inv.Miscellaneous))
	sb.WriteString(fmt.Sprintf("Cash: $%d", m.state.Player.Cash))

	if m.state.Flags.Injured {
		sb.WriteString("\n" + BadStyle.Render("⚠ INJURED"))
	}
	if m.state.Flags.Ill {
		sb.WriteString("\n" + BadStyle.Render("⚠ ILL"))
	}

	return sb.String()
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
