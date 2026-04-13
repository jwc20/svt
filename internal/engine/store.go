package engine

type GameStore interface {
	CreatePlayer(publicKey string) (playerID int64, err error)
	GetPlayerByKey(publicKey string) (playerID int64, err error)
	SaveGame(playerID int64, gameID int64, state *GameState) error
	NewGame(playerID int64, state *GameState) (gameID int64, err error)
	LoadActiveGame(playerID int64) (gameID int64, state *GameState, err error)
	FinishGame(gameID int64, score *int) error
}
