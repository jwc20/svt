package engine

import (
	"math"

	"github.com/jwc20/svt/internal/rand"
)

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

func NeedsAPIGateway(gs *GameState) bool {
	srv := ServerSpecs[gs.Server]
	db := DBSpecs[gs.Database]
	return srv.IsAWS || db.IsAWS
}

func APIGatewayCost(gs *GameState) int {
	if NeedsAPIGateway(gs) {
		return 129
	}
	return 0
}

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

func FixBugs(gs *GameState) (int, int) {
	if gs.ActionChoice != 2 {
		return 0, 0
	}
	bugsFixed := rand.GetRandomInt(4) + 1 // rand(2..5)
	debtFixed := rand.GetRandomInt(3)     // rand(1..3)

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

func SystemDeathRoll(ceiling int) int {
	return rand.GetRandomInt(ceiling)
}

// CalcCashBurn computes the monthly cost.
func CalcCashBurn(gs *GameState) int {
	srv := ServerSpecs[gs.Server]
	db := DBSpecs[gs.Database]
	perUserCost := (srv.PerUserCost + db.PerUserCost) * float64(gs.UserCount)
	return srv.MonthlyCost + db.MonthlyCost + int(math.Ceil(perUserCost)) + APIGatewayCost(gs)
}

func CalcRevenue(gs *GameState) int {
	// revenue = hype * 1.5 + rand(0..hype) + randomEvent (randomEvent handled separately)
	base := int(float64(gs.Hype) * 1.5)
	bonus := 0
	if gs.Hype > 0 {
		bonus = rand.GetRandomInt(gs.Hype+1) - 1 // rand(0..hype)
	}
	return base + bonus
}

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

func GenerateBugs(gs *GameState) int {
	srv := ServerSpecs[gs.Server]
	db := DBSpecs[gs.Database]
	// newBugs = floor(totalMiles/400) + server.bugMod + db.bugMod
	newBugs := gs.Miles / 400

	if srv.BugCeiling > 0 {
		newBugs += rand.GetRandomInt(srv.BugCeiling+1) - 1 // rand(0..BugCeiling)
	}
	if db.BugCeiling > 0 {
		newBugs += rand.GetRandomInt(db.BugCeiling+1) - 1 // rand(0..BugCeiling)
	}

	if newBugs < 0 {
		newBugs = 0
	}
	gs.BugCount += newBugs
	return newBugs
}

func UpdateUserCount(gs *GameState) {
	gs.UserCount = gs.Hype * 10
}

func TechHealth(gs *GameState) int {
	return 100 - gs.TechDebt - gs.BugCount*3
}

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
