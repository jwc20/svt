package engine

// TODO: these should be configurable
const TotalRequiredMileage = 2040
const InitialCash = 700

type Player struct {
	Cash          int
	ShootingLevel int
}

type Inventory struct {
	Oxen          int
	Food          int
	Ammo          int
	Clothing      int
	Miscellaneous int
}

type TripState struct {
	Mileage         int
	PreviousMileage int
	TurnNumber      int
	CurrentDate     int
	EatingChoice    int
	ActionChoice    int
	FortAvailable   bool
	FortSpending    int
}

type Flags struct {
	Injured              bool
	Ill                  bool
	ClearedSouthPass     bool
	ClearedBlueMountains bool
	SouthPassMileage     bool
}

type GameState struct {
	Player    Player
	Inventory Inventory
	Trip      TripState
	Flags     Flags
}

func InitState() GameState {
	return GameState{
		Player: Player{Cash: InitialCash},
		Trip:   TripState{FortAvailable: true},
	}
}
