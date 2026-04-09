package svt

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

// ********************************************************************************************************************

type GameStore interface {
	SaveState(state GameState) error
	LoadState() (GameState, error)
}

type GamePhase int

const (
	PhaseShooting GamePhase = iota
	PhasePurchaseOxen
	PhasePurchaseFood
	PhasePurchaseAmmo
	PhasePurchaseClothing
	PhasePurchaseMisc
	PhaseTurnAction
	PhaseEating
	PhaseGameOver
)

func InitState() GameState {
	return GameState{
		Player: Player{Cash: InitialCash},
		Trip:   TripState{FortAvailable: true},
	}
}

func SetShootingLevel(gs *GameState, level int) bool {
	if level < 1 || level > 5 {
		return false
	}
	gs.Player.ShootingLevel = level
	return true
}
