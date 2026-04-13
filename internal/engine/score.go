package engine

var serverScoreBonus = map[ServerOption]int{
	ServerThinkPad: 200,
	ServerLambda:   100,
	ServerEC2:      50,
	ServerFargate:  0,
}

var dbScoreBonus = map[DBOption]int{
	DBSQLite: 150,
	DBRDS:    75,
	DBAurora: 0,
}

// CalcScore computes the end-game score.
// Score = cash + (hype * 10) + (techHealth * 5) - (totalTurns * 20) + serverBonus + dbBonus
func CalcScore(gs *GameState) int {
	techHealth := TechHealth(gs)
	return gs.Cash + (gs.Hype * 10) + (techHealth * 5) - (gs.TurnNumber * 20) + serverScoreBonus[gs.Server] + dbScoreBonus[gs.Database]
}
