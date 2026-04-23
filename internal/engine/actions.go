package engine

import (
	"math"
	"math/rand"

	"charm.land/log/v2"
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
	r := rand.Intn(41) - 20 // -20..20
	log.Info("AdvanceMileage called", "randomInt", r)
	miles := 140 + gs.Hype/5 + r
	if gs.ActionChoice == 2 || gs.ActionChoice == 3 {
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
	bugsFixedRand := rand.Intn(4) // 0..3
	debtFixedRand := rand.Intn(3) // 0..2
	log.Info("FixBugs called", "bugsFixedRandomInt", bugsFixedRand, "debtFixedRandomInt", debtFixedRand)

	bugsFixed := bugsFixedRand + 2 // 2..5
	debtFixed := debtFixedRand + 1 // 1..3

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

func MarketingPushCost(gs *GameState) int {
	return int(math.Round(50 + float64(gs.Hype)*0.5))
}

func MarketingPush(gs *GameState) (int, int) {
	cost := MarketingPushCost(gs)
	gs.Cash -= cost
	hypeGained := 16 + rand.Intn(11)
	log.Info("MarketingPush called", "cost", cost, "hypeGained", hypeGained)
	gs.Hype += hypeGained
	return hypeGained, cost
}

// returns a roll in [1, ceiling]
func SystemDeathRoll(ceiling int) int {
	r := rand.Intn(ceiling) + 1
	log.Info("SystemDeathRoll called", "ceiling", ceiling, "randomInt", r)
	return r
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
		bonus = rand.Intn(gs.Hype + 1) // 0..hype
		log.Info("CalcRevenue called", "randomInt", bonus)
	}
	return base + bonus
}

func ApplyHypeDecay(gs *GameState) int {
	// hype = hype - 3 - bugCount/2 + rand(-5..5)
	r := rand.Intn(11) - 5 // -5..5
	log.Info("ApplyHypeDecay called", "randomInt", r)
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
	r := rand.Intn(4) // 0..3
	log.Info("AccumulateTechDebt called", "randomInt", r)
	added := gs.Miles/200 + srv.DebtPerTurn + db.DebtPerTurn + r
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
		r := rand.Intn(srv.BugCeiling + 1) // 0..BugCeiling
		log.Info("GenerateBugs called", "source", "server", "randomInt", r)
		newBugs += r
	}
	if db.BugCeiling > 0 {
		r := rand.Intn(db.BugCeiling + 1) // 0..BugCeiling
		log.Info("GenerateBugs called", "source", "database", "randomInt", r)
		newBugs += r
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

func ApplyEndOfTurn(gs *GameState) (int, int, int, int, int, string, int) {
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
	eventMsg, eventID := GenerateEvent(gs)

	return cashBurn, revenue, hypeDecay, techDebtAdded, bugsAdded, eventMsg, eventID
}
