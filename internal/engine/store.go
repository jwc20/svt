package engine

type GameStore interface {
	SaveState(state GameState) error
	LoadState() (GameState, error)
}
