package engine

import "math/rand"

const (
	TotalRequiredProduct = 2040
	InitialCash          = 1500
	InitialHypeBase      = 50
	TotalTurns           = 12
)

// ServerOption represents infrastructure server choices.
type ServerOption int

const (
	ServerFargate  ServerOption = iota // AWS Fargate
	ServerEC2                          // AWS EC2
	ServerLambda                       // AWS Lambda
	ServerThinkPad                     // Lenovo ThinkPad
)

// DBOption represents database choices.
type DBOption int

const (
	DBAurora DBOption = iota // AWS Aurora
	DBRDS                    // AWS RDS
	DBSQLite                 // SQLite
)

type ServerSpec struct {
	Name        string
	MonthlyCost int
	PerUserCost float64
	DebtPerTurn int
	BugCeiling  int // rand(0..BugCeiling) bugs per turn; 0 means no random bugs
	IsAWS       bool
}

type DBSpec struct {
	Name        string
	MonthlyCost int
	PerUserCost float64
	DebtPerTurn int
	BugCeiling  int
	IsAWS       bool
}

var ServerSpecs = map[ServerOption]ServerSpec{
	ServerFargate:  {Name: "AWS Fargate", MonthlyCost: 0, PerUserCost: 0.05, DebtPerTurn: 0, BugCeiling: 0, IsAWS: true},
	ServerEC2:      {Name: "AWS EC2", MonthlyCost: 40, PerUserCost: 0, DebtPerTurn: 1, BugCeiling: 1, IsAWS: true},
	ServerLambda:   {Name: "AWS Lambda", MonthlyCost: 0, PerUserCost: 0.03, DebtPerTurn: 2, BugCeiling: 2, IsAWS: true},
	ServerThinkPad: {Name: "Lenovo ThinkPad", MonthlyCost: 0, PerUserCost: 0, DebtPerTurn: 4, BugCeiling: 3, IsAWS: false},
}

var DBSpecs = map[DBOption]DBSpec{
	DBAurora: {Name: "AWS Aurora", MonthlyCost: 0, PerUserCost: 0.04, DebtPerTurn: 0, BugCeiling: 0, IsAWS: true},
	DBRDS:    {Name: "AWS RDS", MonthlyCost: 30, PerUserCost: 0, DebtPerTurn: 1, BugCeiling: 1, IsAWS: true},
	DBSQLite: {Name: "SQLite", MonthlyCost: 0, PerUserCost: 0, DebtPerTurn: 3, BugCeiling: 2, IsAWS: false},
}

type DeathRollOutcome int

const (
	DeathRollNone DeathRollOutcome = iota // action was push forward
	DeathRollWin
	DeathRollLoss
)

type TurnEntry struct {
	Action    int              // 1 = push forward, 2 = fix bugs
	DeathRoll DeathRollOutcome // only meaningful when Action == 2
	EventID   int              // 0 = no event, 1-19 = specific event
}

type GameState struct {
	Cash             int
	Hype             int
	TechDebt         int
	BugCount         int
	ProductReadiness int
	UserCount        int

	Server   ServerOption
	Database DBOption

	TurnNumber   int
	ActionChoice int // 1 = push forward, 2 = fix bugs, 3 = marketing push

	TurnHistory []TurnEntry
}

func InitState(bonusHype int) GameState {
	hype := InitialHypeBase + rand.Intn(50) + bonusHype
	return GameState{
		Cash:             InitialCash,
		Hype:             hype,
		TechDebt:         0,
		BugCount:         0,
		ProductReadiness: 0,
		Server:           ServerFargate, // default, will be chosen by player
		Database:         DBAurora,      // default, will be chosen by player
	}
}
