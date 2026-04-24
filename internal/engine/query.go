package engine

var Route = []string{
	"San Jose",
	"Santa Clara",
	"Sunnyvale",
	"Mountain View",
	"Palo Alto",
	"Menlo Park",
	"Redwood City",
	"San Mateo",
	"Hillsborough",
	"San Bruno",
	"Daly City",
	"San Francisco",
}

func CurrentLocation(turn int) string {
	if turn < 1 {
		return Route[0]
	}
	if turn > len(Route) {
		return Route[len(Route)-1]
	}
	return Route[turn-1]
}

func IsBankrupt(gs *GameState) bool {
	return gs.Cash < 0
}

func IsGhostTown(gs *GameState) bool {
	return gs.Hype < 5
}

func IsSystemFailure(gs *GameState) bool {
	return TechHealth(gs) < 0
}

func IsArrived(gs *GameState) bool {
	return gs.ProductReadiness >= TotalRequiredProduct
}

func CheckLoseCondition(gs *GameState) (string, bool) {
	if IsBankrupt(gs) {
		return "BANKRUPT! You ran out of cash.", true
	}
	if IsGhostTown(gs) {
		return "GHOST TOWN! Your hype dropped too low -- no one cares about your product anymore.", true
	}
	if IsSystemFailure(gs) {
		return "TOTAL SYSTEM FAILURE! Tech debt and bugs have destroyed your infrastructure.", true
	}
	return "", false
}
