package engine

func IsStarved(gs *GameState) bool {
	return gs.Inventory.Food < 0
}

func IsArrived(gs *GameState) bool {
	return gs.Trip.Mileage >= TotalRequiredMileage
}

func NeedsAilmentCheck(gs *GameState) bool {
	return gs.Flags.Injured || gs.Flags.Ill
}

func DateName(turn int) string {
	dates := []string{
		"MARCH 29", "APRIL 12", "APRIL 26", "MAY 10", "MAY 24",
		"JUNE 7", "JUNE 21", "JULY 5", "JULY 19", "AUGUST 2",
		"AUGUST 16", "AUGUST 31", "SEPTEMBER 13", "SEPTEMBER 27",
		"OCTOBER 11", "OCTOBER 25", "NOVEMBER 8", "NOVEMBER 22",
		"DECEMBER 6", "DECEMBER 20",
	}
	weekdays := []string{
		"SATURDAY", "SUNDAY", "MONDAY", "TUESDAY",
		"WEDNESDAY", "THURSDAY", "FRIDAY",
	}
	if turn < 1 || turn > len(dates) {
		return "WINTER"
	}
	idx := turn - 1
	return weekdays[idx%len(weekdays)] + ", " + dates[idx] + ", 1847"
}
