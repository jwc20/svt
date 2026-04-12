package engine

// Route is the 12-turn path from San Jose to San Francisco.
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

// CurrentLocation returns the name of the current city based on turn number.
func CurrentLocation(turn int) string {
	if turn < 1 {
		return Route[0]
	}
	if turn > len(Route) {
		return Route[len(Route)-1]
	}
	return Route[turn-1]
}

// IsBankrupt returns true if cash has gone below 0.
func IsBankrupt(gs *GameState) bool {
	return gs.Cash < 0
}

// IsGhostTown returns true if hype dropped below 5.
func IsGhostTown(gs *GameState) bool {
	return gs.Hype < 5
}

// IsSystemFailure returns true if tech health dropped below 0.
func IsSystemFailure(gs *GameState) bool {
	return TechHealth(gs) < 0
}

// IsArrived returns true if mileage reached the win threshold.
func IsArrived(gs *GameState) bool {
	return gs.Miles >= TotalRequiredMileage
}

// CheckLoseCondition returns ("", false) if no lose condition met,
// or (reason, true) if the game is lost.
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
