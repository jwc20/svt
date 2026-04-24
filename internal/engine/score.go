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

	cashScore := gs.Cash
	hypeScore := gs.Hype
	if gs.Difficulty == Diff1 {
		cashScore = cashScore - DifficultySpecs[Diff1].BonusCash
		hypeScore = hypeScore - DifficultySpecs[Diff1].BonusHype
	} else if gs.Difficulty == Diff2 {
		cashScore = cashScore - DifficultySpecs[Diff2].BonusCash
		hypeScore = hypeScore - DifficultySpecs[Diff2].BonusHype
	}

	return cashScore + (hypeScore * 10) + (techHealth * 5) - (gs.TurnNumber * 20) + serverScoreBonus[gs.Server] + dbScoreBonus[gs.Database]
}
