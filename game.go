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

// ********************************************************************************************************************

func SetShootingLevel(gs *GameState, level int) bool {
	if level < 1 || level > 5 {
		return false
	}
	gs.Player.ShootingLevel = level
	return true
}

func PurchaseItem(gs *GameState, phase GamePhase, amount int) (bool, string) {
	switch phase {
	case PhasePurchaseOxen:
		if amount < 200 || amount > 300 {
			return false, "AMOUNT MUST BE BETWEEN $200 AND $300"
		}
		gs.Inventory.Oxen = amount
	case PhasePurchaseFood:
		if amount < 100 || amount > 200 {
			return false, "AMOUNT MUST BE BETWEEN $100 AND $200"
		}
		gs.Inventory.Food = amount
	case PhasePurchaseAmmo:
		if amount < 50 || amount > 100 {
			return false, "AMOUNT MUST BE BETWEEN $50 AND $100"
		}
		gs.Inventory.Ammo = amount
	case PhasePurchaseClothing:
		if amount < 50 || amount > 100 {
			return false, "AMOUNT MUST BE BETWEEN $50 AND $100"
		}
		gs.Inventory.Clothing = amount
	case PhasePurchaseMisc:
		if amount < 50 || amount > 100 {
			return false, "AMOUNT MUST BE BETWEEN $50 AND $100"
		}
		gs.Inventory.Miscellaneous = amount
	default:
		panic("unhandled default case")
	}
	return true, ""
}
