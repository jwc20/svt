package svt

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

func FinalizePurchases(gs *GameState) (bool, int) {
	spent := gs.Inventory.Oxen + gs.Inventory.Food + gs.Inventory.Ammo + gs.Inventory.Clothing + gs.Inventory.Miscellaneous
	remaining := InitialCash - spent
	if remaining < 0 {
		return false, remaining
	}
	gs.Player.Cash = remaining
	return true, remaining
}

func ApplyEating(gs *GameState, choice int) {
	if choice < 1 || choice > 3 {
		choice = 2
	}
	gs.Trip.EatingChoice = choice
	gs.Inventory.Food -= 8 + 5*choice
}

func AdvanceMileage(gs *GameState) {
	gs.Trip.PreviousMileage = gs.Trip.Mileage
	r := GetRandomInt(10)
	miles := 200 + (gs.Inventory.Oxen-220)/5 + r*10
	if gs.Trip.ActionChoice != 1 {
		miles /= 2
	}
	gs.Trip.Mileage += miles
}

func GenerateEvent(gs *GameState) string {
	r := GetRandomInt(100)

	switch {
	case r <= 6:
		gs.Trip.Mileage -= 15 + GetRandomInt(10)
		gs.Inventory.Miscellaneous -= 8 + GetRandomInt(5)
		return "WAGON BREAKS DOWN — LOSS OF TIME AND SUPPLIES."
	case r <= 11:
		gs.Trip.Mileage -= 25
		gs.Inventory.Oxen -= 20
		return "OX INJURED — LOSS OF TIME."
	case r <= 15:
		gs.Flags.Injured = true
		gs.Inventory.Miscellaneous -= 5 + GetRandomInt(4)
		return "BAD LUCK — YOUR DAUGHTER BROKE HER ARM."
	case r <= 20:
		gs.Inventory.Ammo -= 10 + GetRandomInt(5)
		if gs.Inventory.Ammo < 0 {
			gs.Inventory.Food -= 30 + GetRandomInt(20)
			return "WILD ANIMALS ATTACK! YOU RAN OUT OF BULLETS — THEY GOT SOME FOOD."
		}
		return "WILD ANIMALS ATTACK!"
	case r <= 25:
		if gs.Inventory.Clothing < 20 {
			gs.Flags.Ill = true
			return "COLD WEATHER — BRRRR! YOU DON'T HAVE ENOUGH CLOTHING."
		}
		return "COLD WEATHER — BRRRR! BUT YOU'RE DRESSED WARM."
	case r <= 30:
		gs.Trip.Mileage -= 10 + GetRandomInt(5)
		gs.Inventory.Food -= 10
		gs.Inventory.Ammo -= 5 + GetRandomInt(5)
		gs.Inventory.Miscellaneous -= 5 + GetRandomInt(5)
		return "HEAVY RAINS — TIME LOST AND SUPPLIES DAMAGED."
	case r <= 33:
		gs.Inventory.Food -= 10 + GetRandomInt(10)
		gs.Player.Cash -= 10 + GetRandomInt(15)
		if gs.Player.Cash < 0 {
			gs.Player.Cash = 0
		}
		return "BANDITS ATTACK!"
	case r <= 36:
		gs.Inventory.Food -= 40 + GetRandomInt(30)
		gs.Inventory.Ammo -= 20 + GetRandomInt(20)
		gs.Inventory.Miscellaneous -= 10 + GetRandomInt(10)
		return "FIRE IN YOUR WAGON — LOSS OF SUPPLIES."
	case r <= 40:
		gs.Inventory.Food += 14 + GetRandomInt(5)
		return "HELPFUL INDIANS SHOW YOU WHERE TO FIND FOOD."
	default:
		return ""
	}
}
