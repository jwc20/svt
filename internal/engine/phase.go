package engine

type GamePhase int

const (
	PhaseServerChoice GamePhase = iota
	PhaseDBChoice
	PhaseTurnAction
	PhaseGameOver
)
