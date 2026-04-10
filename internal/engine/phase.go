package engine

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
