package engine

import (
	"math"

	"github.com/jwc20/svt/internal/rand"
)

// SetServer sets the server infrastructure choice.
func SetServer(gs *GameState, choice int) bool {
	switch choice {
	case 1:
		gs.Server = ServerFargate
	case 2:
		gs.Server = ServerEC2
	case 3:
		gs.Server = ServerLambda
	case 4:
		gs.Server = ServerThinkPad
	default:
		return false
	}
	return true
}

// SetDatabase sets the database choice.
func SetDatabase(gs *GameState, choice int) bool {
	switch choice {
	case 1:
		gs.Database = DBAurora
	case 2:
		gs.Database = DBRDS
	case 3:
		gs.Database = DBSQLite
	default:
		return false
	}
	return true
}

// NeedsAPIGateway returns true if either server or database is an AWS service.
func NeedsAPIGateway(gs *GameState) bool {
	srv := ServerSpecs[gs.Server]
	db := DBSpecs[gs.Database]
	return srv.IsAWS || db.IsAWS
}

// APIGatewayCost returns 129 if API gateway is needed, 0 otherwise.
func APIGatewayCost(gs *GameState) int {
	if NeedsAPIGateway(gs) {
		return 129
	}
	return 0
}

// AdvanceMileage calculates and adds miles for this turn.
// If action is "fix bugs" (2), miles are halved.
func AdvanceMileage(gs *GameState) int {
	// miles += 140 + hype/5 + rand(-20..20)
	r := rand.GetRandomInt(41) - 21 // gives -20 to 20
	miles := 140 + gs.Hype/5 + r
	if gs.ActionChoice == 2 {
		miles /= 2
	}
	if miles < 0 {
		miles = 0
	}
	gs.Miles += miles
	return miles
}

// FixBugs reduces bugs and tech debt when action choice is "fix bugs".
// Returns (bugsFixed, debtFixed).
func FixBugs(gs *GameState) (int, int) {
	if gs.ActionChoice != 2 {
		return 0, 0
	}
	bugsFixed := rand.GetRandomInt(4) + 1 // rand(2..5)
	debtFixed := rand.GetRandomInt(3)      // rand(1..3)

	gs.BugCount -= bugsFixed
	if gs.BugCount < 0 {
		gs.BugCount = 0
	}
	gs.TechDebt -= debtFixed
	if gs.TechDebt < 0 {
		gs.TechDebt = 0
	}
	return bugsFixed, debtFixed
}

// DeathRoll simulates the WoW death roll mechanic to determine if bugs are completely fixed.
// Player and "system" alternate rolling from the current ceiling down.
// If either rolls 1, they lose. Returns true if player wins (bugs fully fixed).
func DeathRoll(gs *GameState) (bool, []int) {
	ceiling := 100
	rolls := []int{}
	playerTurn := true

	for ceiling > 1 {
		roll := rand.GetRandomInt(ceiling)
		rolls = append(rolls, roll)
		if roll == 1 {
			// whoever rolled 1 loses
			return playerTurn, rolls // if player rolled 1, player loses (false would mean bugs NOT fixed)
		}
		ceiling = roll
		playerTurn = !playerTurn
	}
	// ceiling reached 1, last roller loses
	return !playerTurn, rolls
}

// CalcCashBurn computes the monthly cost.
func CalcCashBurn(gs *GameState) int {
	srv := ServerSpecs[gs.Server]
	db := DBSpecs[gs.Database]
	perUserCost := (srv.PerUserCost + db.PerUserCost) * float64(gs.UserCount)
	return srv.MonthlyCost + db.MonthlyCost + int(math.Ceil(perUserCost)) + APIGatewayCost(gs)
}

// CalcRevenue computes the monthly revenue.
func CalcRevenue(gs *GameState) int {
	// revenue = hype * 1.5 + rand(0..hype) + randomEvent (randomEvent handled separately)
	base := int(float64(gs.Hype) * 1.5)
	bonus := 0
	if gs.Hype > 0 {
		bonus = rand.GetRandomInt(gs.Hype+1) - 1 // rand(0..hype)
	}
	return base + bonus
}

// ApplyHypeDecay reduces hype each turn.
func ApplyHypeDecay(gs *GameState) int {
	// hype = hype - 3 - bugCount/2 + rand(-5..5)
	r := rand.GetRandomInt(11) - 6 // gives -5 to 5
	decay := 3 + gs.BugCount/2 - r
	gs.Hype -= decay
	if gs.Hype < 0 {
		gs.Hype = 0
	}
	return decay
}

// AccumulateTechDebt adds tech debt for this turn.
func AccumulateTechDebt(gs *GameState) int {
	srv := ServerSpecs[gs.Server]
	db := DBSpecs[gs.Database]
	// techDebt += totalMiles/200 + server.debtMod + db.debtMod + rand(0..3)
	added := gs.Miles/200 + srv.DebtPerTurn + db.DebtPerTurn + rand.GetRandomInt(4) - 1
	if added < 0 {
		added = 0
	}
	gs.TechDebt += added
	return added
}

// GenerateBugs adds new bugs for this turn.
func GenerateBugs(gs *GameState) int {
	srv := ServerSpecs[gs.Server]
	db := DBSpecs[gs.Database]
	// newBugs = floor(totalMiles/400) + server.bugMod + db.bugMod
	newBugs := gs.Miles / 400

	if srv.BugCeiling > 0 {
		newBugs += rand.GetRandomInt(srv.BugCeiling + 1) - 1 // rand(0..BugCeiling)
	}
	if db.BugCeiling > 0 {
		newBugs += rand.GetRandomInt(db.BugCeiling + 1) - 1 // rand(0..BugCeiling)
	}

	if newBugs < 0 {
		newBugs = 0
	}
	gs.BugCount += newBugs
	return newBugs
}

// UpdateUserCount recalculates user count from hype.
func UpdateUserCount(gs *GameState) {
	gs.UserCount = gs.Hype * 10
}

// TechHealth computes the derived tech health value.
func TechHealth(gs *GameState) int {
	return 100 - gs.TechDebt - gs.BugCount*3
}

// ApplyEndOfTurn runs all end-of-turn effects in order.
// Returns (cashBurn, revenue, hypeDecay, techDebtAdded, bugsAdded, eventMsg).
func ApplyEndOfTurn(gs *GameState) (int, int, int, int, int, string) {
	// 1. Cash burn deducted
	cashBurn := CalcCashBurn(gs)
	gs.Cash -= cashBurn

	// 2. Revenue collected
	revenue := CalcRevenue(gs)
	gs.Cash += revenue

	// 3. Hype decay applied
	hypeDecay := ApplyHypeDecay(gs)

	// 4. Tech debt accumulates
	techDebtAdded := AccumulateTechDebt(gs)

	// 5. Bugs generated
	bugsAdded := GenerateBugs(gs)

	// 6. Update user count
	UpdateUserCount(gs)

	// 7. Random event rolled
	eventMsg := GenerateEvent(gs)

	return cashBurn, revenue, hypeDecay, techDebtAdded, bugsAdded, eventMsg
}
