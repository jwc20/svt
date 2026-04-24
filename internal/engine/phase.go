package engine

type GamePhase int

const (
	PhaseDifficultyChoice GamePhase = iota
	PhaseServerChoice
	PhaseDBChoice
	PhaseTurnAction
	PhaseDeathRoll
	PhaseGameOver
)
