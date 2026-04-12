package ui

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/jwc20/svt/internal/engine"
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
	state engine.GameState
	phase engine.GamePhase
	store engine.GameStore

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

	deathRollCeiling int // current upper bound for player's next roll
}

func NewGameModel(store engine.GameStore, w, h int) GameModel {
	ti := textinput.New()
	ti.Placeholder = "Enter choice..."
	ti.CharLimit = 20
	ti.SetWidth(maxInt(w/3-6, 14))
	ti.Focus()

	vp := viewport.New()
	vp.SetWidth(maxInt(w/3-4, 14))
	vp.SetHeight(maxInt(h-22, 4))

	gs := engine.InitState()

	m := GameModel{
		state:    gs,
		phase:    engine.PhaseServerChoice,
		store:    store,
		input:    ti,
		choiceVP: vp,
		width:    w,
		height:   h,
	}

	m.setServerPrompt()
	return m
}

func (m GameModel) Init() tea.Cmd { return nil }

func (m GameModel) rightW() int { return maxInt(m.width*45/100, 22) }
func (m GameModel) leftW() int  { return m.width - m.rightW() - 4 }

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
	case engine.PhaseServerChoice:
		return m.handleServerChoice(val)
	case engine.PhaseDBChoice:
		return m.handleDBChoice(val)
	case engine.PhaseTurnAction:
		return m.handleTurnAction(val)
	case engine.PhaseDeathRoll:
		return m.handleDeathRoll(val)
	case engine.PhaseGameOver:
		return m, func() tea.Msg { return BackToLobbyMsg{} }
	}
	return m, nil
}

// ── phase handlers ────────────────────────────────────────────────

func (m GameModel) handleServerChoice(val string) (GameModel, tea.Cmd) {
	choice, err := strconv.Atoi(val)
	if err != nil || !engine.SetServer(&m.state, choice) {
		m.addChoice("Invalid -- enter 1-4")
		return m, nil
	}
	srv := engine.ServerSpecs[m.state.Server]
	m.addChoice(fmt.Sprintf("Server: %s", srv.Name))
	m.phase = engine.PhaseDBChoice
	m.setDBPrompt()
	return m, nil
}

func (m GameModel) handleDBChoice(val string) (GameModel, tea.Cmd) {
	choice, err := strconv.Atoi(val)
	if err != nil || !engine.SetDatabase(&m.state, choice) {
		m.addChoice("Invalid -- enter 1-3")
		return m, nil
	}
	db := engine.DBSpecs[m.state.Database]
	m.addChoice(fmt.Sprintf("Database: %s", db.Name))

	engine.UpdateUserCount(&m.state)

	gw := ""
	if engine.NeedsAPIGateway(&m.state) {
		gw = " (+ $129/mo API Gateway)"
	}
	m.addChoice(fmt.Sprintf("Config set%s", gw))
	m.addChoice(fmt.Sprintf("Starting: $%d cash, %d hype, %d users", m.state.Cash, m.state.Hype, m.state.UserCount))
	m.addChoice("---")

	return m.startTurn()
}

func (m GameModel) handleTurnAction(val string) (GameModel, tea.Cmd) {
	choice, err := strconv.Atoi(val)
	if err != nil || (choice != 1 && choice != 2) {
		m.addChoice("Invalid -- enter 1 or 2")
		return m, nil
	}
	m.state.ActionChoice = choice

	if choice == 1 {
		m.addChoice(">> Push forward")
		miles := engine.AdvanceMileage(&m.state)
		m.addChoice(fmt.Sprintf("  +%d miles (total: %d/%d)", miles, m.state.Miles, engine.TotalRequiredMileage))
		return m.finishTurn()
	}

	m.addChoice(">> Fix bugs")
	m.addChoice("  -- Death Roll --")

	sysRoll := engine.SystemDeathRoll(100)
	m.addChoice(fmt.Sprintf("  System: rolls %d (1-100)", sysRoll))

	if sysRoll == 1 {
		return m.deathRollWin()
	}

	m.deathRollCeiling = sysRoll
	m.phase = engine.PhaseDeathRoll
	m.setDeathRollPrompt()
	return m, nil
}

func (m GameModel) handleDeathRoll(val string) (GameModel, tea.Cmd) {
	if strings.ToLower(val) != "roll" {
		m.addChoice("  Type \"roll\" to roll!")
		return m, nil
	}

	playerRoll := engine.SystemDeathRoll(m.deathRollCeiling)
	m.addChoice(fmt.Sprintf("  You: rolls %d (1-%d)", playerRoll, m.deathRollCeiling))

	if playerRoll == 1 {
		return m.deathRollLose()
	}

	// system rolls automatically
	m.deathRollCeiling = playerRoll
	sysRoll := engine.SystemDeathRoll(m.deathRollCeiling)
	m.addChoice(fmt.Sprintf("  System: rolls %d (1-%d)", sysRoll, m.deathRollCeiling))

	if sysRoll == 1 {
		return m.deathRollWin()
	}

	// back to player
	m.deathRollCeiling = sysRoll
	m.setDeathRollPrompt()
	return m, nil
}

func (m GameModel) deathRollWin() (GameModel, tea.Cmd) {
	m.addChoice(GoodStyle.Render("  You won the death roll!"))

	// apply bug fix + mileage only on win
	bugsFixed, debtFixed := engine.FixBugs(&m.state)
	miles := engine.AdvanceMileage(&m.state)
	m.addChoice(fmt.Sprintf("  Fixed %d bugs, reduced %d tech debt", bugsFixed, debtFixed))
	m.addChoice(fmt.Sprintf("  +%d miles (total: %d/%d)", miles, m.state.Miles, engine.TotalRequiredMileage))
	m.addChoice("")
	return m.finishTurn()
}

func (m GameModel) deathRollLose() (GameModel, tea.Cmd) {
	m.addChoice(WarnStyle.Render("  You rolled 1! You lost the death roll."))

	miles := engine.AdvanceMileage(&m.state)
	m.addChoice(fmt.Sprintf("  +%d miles (total: %d/%d)", miles, m.state.Miles, engine.TotalRequiredMileage))
	m.addChoice("")
	return m.finishTurn()
}

func (m GameModel) finishTurn() (GameModel, tea.Cmd) {
	cashBurn, revenue, _, techDebtAdded, bugsAdded, eventMsg := engine.ApplyEndOfTurn(&m.state)

	m.addChoice(fmt.Sprintf("  Cash burn: -$%d | Revenue: +$%d | Net: $%d", cashBurn, revenue, m.state.Cash))
	m.addChoice(fmt.Sprintf("  Tech debt: +%d (total: %d) | New bugs: +%d (total: %d)", techDebtAdded, m.state.TechDebt, bugsAdded, m.state.BugCount))
	m.addChoice(fmt.Sprintf("  Hype: %d | Users: %d | Tech Health: %d", m.state.Hype, m.state.UserCount, engine.TechHealth(&m.state)))

	if eventMsg != "" {
		m.addChoice(EventStyle.Render("  EVENT: " + eventMsg))
	}

	survived, incidentMsg := engine.CheckIncident(&m.state)
	if !survived {
		m.addChoice(BadStyle.Render("  " + incidentMsg))
	}

	m.addChoice("---")

	if reason, lost := engine.CheckLoseCondition(&m.state); lost {
		m.addChoice(BadStyle.Render(reason))
		m.setGameOver("died", reason)
		return m, nil
	}

	// check win
	if engine.IsArrived(&m.state) {
		m.addChoice(GoodStyle.Render("You made it to San Francisco!"))
		m.setGameOver("won", "")
		return m, nil
	}

	if m.state.TurnNumber >= engine.TotalTurns {
		m.addChoice(BadStyle.Render("You ran out of turns before reaching San Francisco!"))
		m.setGameOver("died", "Ran out of turns")
		return m, nil
	}

	return m.startTurn()
}

// ── turn management ───────────────────────────────────────────────

func (m GameModel) startTurn() (GameModel, tea.Cmd) {
	m.state.TurnNumber++
	m.phase = engine.PhaseTurnAction
	m.setTurnPrompt()
	return m, nil
}

// ── prompt setters ────────────────────────────────────────────────

func (m *GameModel) setServerPrompt() {
	m.promptTitle = "CHOOSE YOUR SERVER INFRASTRUCTURE"
	m.promptLines = []string{
		"Your startup needs a server. Choose wisely!",
		"",
		"(1) AWS Fargate    $0/mo + $0.05/user  | 0 debt/turn | 0 bugs/turn",
		"(2) AWS EC2        $40/mo + $0/user     | +1 debt/turn | 0-1 bugs/turn",
		"(3) AWS Lambda     $0/mo + $0.03/user   | +2 debt/turn | 0-2 bugs/turn",
		"(4) Lenovo ThinkPad $0/mo + $0/user     | +4 debt/turn | 0-3 bugs/turn",
	}
}

func (m *GameModel) setDBPrompt() {
	m.promptTitle = "CHOOSE YOUR DATABASE"
	m.promptLines = []string{
		fmt.Sprintf("Server: %s", engine.ServerSpecs[m.state.Server].Name),
		"",
		"(1) AWS Aurora  $0/mo + $0.04/user  | 0 debt/turn | 0 bugs/turn",
		"(2) AWS RDS     $30/mo + $0/user     | +1 debt/turn | 0-1 bugs/turn",
		"(3) SQLite      $0/mo + $0/user      | +3 debt/turn | 0-2 bugs/turn",
		"",
		DimStyle.Render("Note: AWS API Gateway ($129/mo) applies if any AWS service is selected."),
	}
}

func (m *GameModel) setTurnPrompt() {
	t := m.state.TurnNumber
	location := engine.CurrentLocation(t)
	m.promptTitle = fmt.Sprintf("TURN %d / %d -- %s", t, engine.TotalTurns, location)
	m.promptLines = []string{
		"What do you want to do?",
		"",
		"(1) Push forward -- advance miles normally",
		"(2) Fix bugs -- miles halved, fix bugs + death roll",
	}
}

func (m *GameModel) setDeathRollPrompt() {
	m.promptTitle = "DEATH ROLL"
	m.promptLines = []string{
		"Type \"roll\" to roll!",
		"",
		DimStyle.Render("If you roll 1, you lose."),
		DimStyle.Render("If the system rolls 1, you win!"),
	}
}

func (m *GameModel) setGameOver(result, deathMsg string) {
	m.phase = engine.PhaseGameOver
	m.gameOver = true
	m.gameResult = result
	m.deathMessage = deathMsg

	switch result {
	case "won":
		m.promptTitle = "CONGRATULATIONS!"
		m.promptLines = []string{
			"", GoodStyle.Render("YOUR STARTUP MADE IT TO SAN FRANCISCO!"), "",
			fmt.Sprintf("Turns taken: %d", m.state.TurnNumber),
			fmt.Sprintf("Cash remaining: $%d", m.state.Cash),
			fmt.Sprintf("Final hype: %d", m.state.Hype),
			fmt.Sprintf("Final users: %d", m.state.UserCount),
			fmt.Sprintf("Tech health: %d", engine.TechHealth(&m.state)),
			"", DimStyle.Render("Press esc to return to lobby."),
		}
	default:
		m.promptTitle = "GAME OVER"
		m.promptLines = []string{
			"", BadStyle.Render(deathMsg), "",
			fmt.Sprintf("Mileage reached: %d / %d", m.state.Miles, engine.TotalRequiredMileage),
			fmt.Sprintf("Turns played: %d", m.state.TurnNumber),
			fmt.Sprintf("Cash: $%d", m.state.Cash),
			fmt.Sprintf("Hype: %d", m.state.Hype),
			fmt.Sprintf("Tech health: %d", engine.TechHealth(&m.state)),
			"", DimStyle.Render("Press esc to return to lobby."),
		}
	}
}

// ********************************************************************************************************************
// ── View ──────────────────────────────────────────────────────────
// ********************************************************************************************************************

func (m GameModel) View() tea.View {
	lw := m.leftW()
	rw := m.rightW()
	innerH := m.height - 4

	promptText := FocusLabel.Render(m.promptTitle) + "\n\n" +
		strings.Join(m.promptLines, "\n")
	leftContentH := maxInt(innerH-4, 6)
	centeredPrompt := lipgloss.Place(lw-6, leftContentH,
		lipgloss.Center, lipgloss.Center, promptText)
	leftCol := promptPanel.Width(lw - 2).Height(leftContentH).Render(centeredPrompt)

	sBox := statusBox.Width(rw - 2).Render(m.renderStatus())
	statusRenderedH := lipgloss.Height(sBox)

	logPanelH := maxInt(innerH-statusRenderedH-2, 4)
	logInnerH := maxInt(logPanelH, 2)

	vpW := maxInt(rw-6, 10)
	wrapped := lipgloss.Wrap(strings.Join(m.choiceLog, "\n"), vpW, "")
	m.choiceVP.SetContent(wrapped)

	var logInner string
	if m.gameOver {
		m.choiceVP.SetWidth(vpW)
		m.choiceVP.SetHeight(logInnerH)
		m.choiceVP.GotoBottom()
		logInner = m.choiceVP.View()
	} else {
		inputLine := PromptStyle.Render("") + m.input.View()
		vpH := maxInt(logInnerH-2, 1)
		m.choiceVP.SetWidth(vpW)
		m.choiceVP.SetHeight(vpH)
		m.choiceVP.GotoBottom()
		separator := DimStyle.Render(strings.Repeat("~", vpW))
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

	turn := m.state.TurnNumber
	if turn < 1 {
		sb.WriteString(PlainLabel.Render("Location") + "\n")
		sb.WriteString("Setting up...\n")
	} else {
		sb.WriteString(PlainLabel.Render("Location") + "\n")
		sb.WriteString(fmt.Sprintf("Turn %d: %s\n", turn, engine.CurrentLocation(turn)))
	}

	sb.WriteString(PlainLabel.Render("Mileage") + "\n")
	sb.WriteString(fmt.Sprintf("%d / %d\n", m.state.Miles, engine.TotalRequiredMileage))

	sb.WriteString("\n")
	sb.WriteString(PlainLabel.Render("Startup Stats") + "\n")

	cashStr := fmt.Sprintf("$%d", m.state.Cash)
	if m.state.Cash < 200 {
		cashStr = BadStyle.Render(cashStr)
	} else if m.state.Cash < 500 {
		cashStr = WarnStyle.Render(cashStr)
	}
	sb.WriteString(fmt.Sprintf("Cash: %s\n", cashStr))

	hypeStr := fmt.Sprintf("%d", m.state.Hype)
	if m.state.Hype < 15 {
		hypeStr = BadStyle.Render(hypeStr)
	} else if m.state.Hype < 30 {
		hypeStr = WarnStyle.Render(hypeStr)
	}
	sb.WriteString(fmt.Sprintf("Hype: %s\n", hypeStr))

	sb.WriteString(fmt.Sprintf("Users: %d\n", m.state.UserCount))

	th := engine.TechHealth(&m.state)
	thStr := fmt.Sprintf("%d", th)
	if th < 20 {
		thStr = BadStyle.Render(thStr)
	} else if th < 40 {
		thStr = WarnStyle.Render(thStr)
	}
	sb.WriteString(fmt.Sprintf("Tech Health: %s\n", thStr))

	sb.WriteString(fmt.Sprintf("Tech Debt: %d\n", m.state.TechDebt))
	sb.WriteString(fmt.Sprintf("Bugs: %d\n", m.state.BugCount))

	sb.WriteString("\n")
	sb.WriteString(PlainLabel.Render("Infrastructure") + "\n")
	sb.WriteString(fmt.Sprintf("Server: %s\n", engine.ServerSpecs[m.state.Server].Name))
	sb.WriteString(fmt.Sprintf("DB: %s", engine.DBSpecs[m.state.Database].Name))

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
