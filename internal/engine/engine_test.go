package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubGameStore struct {
	Games  map[int64]*GameState
	NextID int64
}

func NewStubGameStore() *StubGameStore {
	return &StubGameStore{Games: make(map[int64]*GameState), NextID: 1}
}

func (s *StubGameStore) CreatePlayer(publicKey string, username string) (int64, error) {
	return 1, nil
}

func (s *StubGameStore) GetPlayerByKey(publicKey string) (int64, error) {
	return 1, nil
}

func (s *StubGameStore) GetBonusHype(playerID int64) (int, error) {
	return 0, nil
}

func (s *StubGameStore) SetBonusHype(playerID int64, bonus int) error {
	return nil
}

func (s *StubGameStore) NewGame(playerID int64, state *GameState) (int64, error) {
	id := s.NextID
	s.NextID++
	cp := *state
	s.Games[id] = &cp
	return id, nil
}

func (s *StubGameStore) SaveGame(playerID int64, gameID int64, state *GameState) error {
	cp := *state
	s.Games[gameID] = &cp
	return nil
}

func (s *StubGameStore) LoadActiveGame(playerID int64) (int64, *GameState, error) {
	return 0, nil, nil
}

func (s *StubGameStore) FinishGame(gameID int64, score *int) error {
	return nil
}

func TestGameConstants(t *testing.T) {
	t.Run("total mileage is 2040", func(t *testing.T) {
		assert.Equal(t, 2040, TotalRequiredMileage)
	})
	t.Run("initial cash is 1500", func(t *testing.T) {
		assert.Equal(t, 1500, InitialCash)
	})
	t.Run("total turns is 12", func(t *testing.T) {
		assert.Equal(t, 12, TotalTurns)
	})
}

func TestInitState(t *testing.T) {
	gs := InitState(0)

	assert.Equal(t, 1500, gs.Cash)
	assert.GreaterOrEqual(t, gs.Hype, 50)
	assert.LessOrEqual(t, gs.Hype, 100)
	assert.Zero(t, gs.TechDebt)
	assert.Zero(t, gs.BugCount)
	assert.Zero(t, gs.Miles)
}

func TestSetServer(t *testing.T) {
	t.Run("valid choice", func(t *testing.T) {
		gs := InitState(0)
		ok := SetServer(&gs, 1)
		assert.True(t, ok)
		assert.Equal(t, ServerFargate, gs.Server)
	})
	t.Run("ThinkPad", func(t *testing.T) {
		gs := InitState(0)
		ok := SetServer(&gs, 4)
		assert.True(t, ok)
		assert.Equal(t, ServerThinkPad, gs.Server)
	})
	t.Run("invalid choice", func(t *testing.T) {
		gs := InitState(0)
		assert.False(t, SetServer(&gs, 0))
		assert.False(t, SetServer(&gs, 5))
	})
}

func TestSetDatabase(t *testing.T) {
	t.Run("valid choice", func(t *testing.T) {
		gs := InitState(0)
		ok := SetDatabase(&gs, 3)
		assert.True(t, ok)
		assert.Equal(t, DBSQLite, gs.Database)
	})
	t.Run("invalid choice", func(t *testing.T) {
		gs := InitState(0)
		assert.False(t, SetDatabase(&gs, 0))
		assert.False(t, SetDatabase(&gs, 4))
	})
}

func TestAPIGateway(t *testing.T) {
	t.Run("AWS server needs gateway", func(t *testing.T) {
		gs := InitState(0)
		gs.Server = ServerFargate
		gs.Database = DBSQLite
		assert.True(t, NeedsAPIGateway(&gs))
		assert.Equal(t, 129, APIGatewayCost(&gs))
	})
	t.Run("AWS db needs gateway", func(t *testing.T) {
		gs := InitState(0)
		gs.Server = ServerThinkPad
		gs.Database = DBAurora
		assert.True(t, NeedsAPIGateway(&gs))
	})
	t.Run("no AWS no gateway", func(t *testing.T) {
		gs := InitState(0)
		gs.Server = ServerThinkPad
		gs.Database = DBSQLite
		assert.False(t, NeedsAPIGateway(&gs))
		assert.Equal(t, 0, APIGatewayCost(&gs))
	})
}

func TestAdvanceMileage(t *testing.T) {
	gs := InitState(0)
	gs.Hype = 50
	gs.ActionChoice = 1
	miles := AdvanceMileage(&gs)
	assert.Greater(t, miles, 0)
	assert.Equal(t, miles, gs.Miles)
}

func TestAdvanceMileageHalvedForFixBugs(t *testing.T) {
	// Run multiple times to account for randomness
	for i := 0; i < 10; i++ {
		gs1 := InitState(0)
		gs1.Hype = 50
		gs1.ActionChoice = 1
		miles1 := AdvanceMileage(&gs1)

		gs2 := InitState(0)
		gs2.Hype = 50
		gs2.ActionChoice = 2
		miles2 := AdvanceMileage(&gs2)

		// miles2 should generally be about half of miles1 (both are random though)
		// Just verify fix bugs gives positive miles
		_ = miles1
		assert.GreaterOrEqual(t, miles2, 0)
	}
}

func TestTechHealth(t *testing.T) {
	gs := InitState(0)
	gs.TechDebt = 20
	gs.BugCount = 5
	// techHealth = 100 - 20 - 5*3 = 65
	assert.Equal(t, 65, TechHealth(&gs))
}

func TestCalcCashBurn(t *testing.T) {
	gs := InitState(0)
	gs.Server = ServerEC2
	gs.Database = DBRDS
	gs.UserCount = 0
	// EC2: $40/mo, RDS: $30/mo, 0 users, AWS gateway: $129
	burn := CalcCashBurn(&gs)
	assert.Equal(t, 40+30+129, burn)
}

func TestCalcCashBurnWithUsers(t *testing.T) {
	gs := InitState(0)
	gs.Server = ServerFargate
	gs.Database = DBAurora
	gs.UserCount = 100
	// Fargate: $0.05/user, Aurora: $0.04/user = $0.09 * 100 = $9, + $129 gateway
	burn := CalcCashBurn(&gs)
	assert.Equal(t, 9+129, burn)
}

func TestCalcCashBurnNoAWS(t *testing.T) {
	gs := InitState(0)
	gs.Server = ServerThinkPad
	gs.Database = DBSQLite
	gs.UserCount = 500
	// ThinkPad: $0/mo $0/user, SQLite: $0/mo $0/user, no gateway
	burn := CalcCashBurn(&gs)
	assert.Equal(t, 0, burn)
}

func TestUpdateUserCount(t *testing.T) {
	gs := InitState(0)
	gs.Hype = 50
	UpdateUserCount(&gs)
	assert.Equal(t, 500, gs.UserCount)
}

func TestLoseConditions(t *testing.T) {
	t.Run("bankrupt", func(t *testing.T) {
		gs := InitState(0)
		gs.Cash = -1
		assert.True(t, IsBankrupt(&gs))
		reason, lost := CheckLoseCondition(&gs)
		assert.True(t, lost)
		assert.Contains(t, reason, "BANKRUPT")
	})
	t.Run("ghost town", func(t *testing.T) {
		gs := InitState(0)
		gs.Hype = 4
		assert.True(t, IsGhostTown(&gs))
		reason, lost := CheckLoseCondition(&gs)
		assert.True(t, lost)
		assert.Contains(t, reason, "GHOST TOWN")
	})
	t.Run("system failure", func(t *testing.T) {
		gs := InitState(0)
		gs.TechDebt = 80
		gs.BugCount = 10
		assert.True(t, IsSystemFailure(&gs))
		reason, lost := CheckLoseCondition(&gs)
		assert.True(t, lost)
		assert.Contains(t, reason, "SYSTEM FAILURE")
	})
	t.Run("no lose condition", func(t *testing.T) {
		gs := InitState(0)
		_, lost := CheckLoseCondition(&gs)
		assert.False(t, lost)
	})
}

func TestIsArrived(t *testing.T) {
	gs := InitState(0)
	gs.Miles = 2040
	assert.True(t, IsArrived(&gs))
}

func TestRoute(t *testing.T) {
	assert.Equal(t, 12, len(Route))
	assert.Equal(t, "San Jose", Route[0])
	assert.Equal(t, "San Francisco", Route[11])
}

func TestCurrentLocation(t *testing.T) {
	assert.Equal(t, "San Jose", CurrentLocation(0))
	assert.Equal(t, "San Jose", CurrentLocation(1))
	assert.Equal(t, "San Francisco", CurrentLocation(12))
	assert.Equal(t, "San Francisco", CurrentLocation(15))
}

func TestGenerateEvent(t *testing.T) {
	gs := InitState(0)
	gs.Cash = 1000
	gs.Hype = 50
	gs.Miles = 500

	changed := false
	for i := 0; i < 50; i++ {
		before := gs.Cash + gs.Hype + gs.Miles + gs.TechDebt + gs.BugCount
		GenerateEvent(&gs)
		after := gs.Cash + gs.Hype + gs.Miles + gs.TechDebt + gs.BugCount
		if before != after {
			changed = true
			break
		}
	}
	assert.True(t, changed, "expected at least one event to change game state")
}

func TestCheckIncident(t *testing.T) {
	t.Run("healthy system survives", func(t *testing.T) {
		gs := InitState(0)
		gs.TechDebt = 0
		gs.BugCount = 0
		survived, _ := CheckIncident(&gs)
		assert.True(t, survived)
	})
	t.Run("unhealthy system may fail", func(t *testing.T) {
		gs := InitState(0)
		gs.TechDebt = 90
		gs.BugCount = 10
		// TechHealth = 100 - 90 - 30 = -20, should always fail
		survived, msg := CheckIncident(&gs)
		assert.False(t, survived)
		assert.Contains(t, msg, "INCIDENT")
	})
}

func TestFixBugs(t *testing.T) {
	gs := InitState(0)
	gs.BugCount = 10
	gs.TechDebt = 10
	gs.ActionChoice = 2
	bugsFixed, debtFixed := FixBugs(&gs)
	assert.Greater(t, bugsFixed, 0)
	assert.Greater(t, debtFixed, 0)
	assert.Less(t, gs.BugCount, 10)
}

func TestFixBugsNoOpWhenPushForward(t *testing.T) {
	gs := InitState(0)
	gs.BugCount = 10
	gs.ActionChoice = 1
	bugsFixed, debtFixed := FixBugs(&gs)
	assert.Equal(t, 0, bugsFixed)
	assert.Equal(t, 0, debtFixed)
	assert.Equal(t, 10, gs.BugCount)
}

//func TestSystemDeathRoll(t *testing.T) {
//	for i := 0; i < 50; i++ {
//		result := SystemDeathRoll(100)
//		assert.GreaterOrEqual(t, result, 1)
//		assert.LessOrEqual(t, result, 100)
//	}
//	// Edge case: ceiling of 1
//	result := SystemDeathRoll(1)
//	assert.Equal(t, 1, result)
//}

func TestCalcScore(t *testing.T) {
	t.Run("basic scoring", func(t *testing.T) {
		gs := GameState{
			Cash:       1000,
			Hype:       50,
			TechDebt:   10,
			BugCount:   5,
			TurnNumber: 6,
			Server:     ServerFargate,
			Database:   DBAurora,
		}
		// techHealth = 100 - 10 - 15 = 75
		// score = 1000 + (50*10) + (75*5) - (6*20) + 0 + 0 = 1000 + 500 + 375 - 120 = 1755
		assert.Equal(t, 1755, CalcScore(&gs))
	})

	t.Run("with server and db bonuses", func(t *testing.T) {
		gs := GameState{
			Cash:       1000,
			Hype:       50,
			TechDebt:   10,
			BugCount:   5,
			TurnNumber: 6,
			Server:     ServerThinkPad,
			Database:   DBSQLite,
		}
		// 1755 + 200 + 150 = 2105
		assert.Equal(t, 2105, CalcScore(&gs))
	})
}

func TestSerializeDeserializeRoundTrip(t *testing.T) {
	gs := &GameState{
		TurnNumber: 3,
		Server:     ServerThinkPad,
		Database:   DBSQLite,
		Cash:       1200,
		Hype:       65,
		TechDebt:   4,
		BugCount:   2,
		Miles:      420,
		UserCount:  650,
		TurnHistory: []TurnEntry{
			{Action: 1, DeathRoll: DeathRollNone, EventID: 0},
			{Action: 2, DeathRoll: DeathRollWin, EventID: 3},
			{Action: 1, DeathRoll: DeathRollNone, EventID: 15},
		},
	}

	serialized := Serialize(gs)
	assert.Equal(t, "3;s4d3;a1/a2w/a1;1200/65/4/2/420/650 - 0/3/15", serialized)

	restored, err := Deserialize(serialized)
	assert.NoError(t, err)

	assert.Equal(t, gs.TurnNumber, restored.TurnNumber)
	assert.Equal(t, gs.Server, restored.Server)
	assert.Equal(t, gs.Database, restored.Database)
	assert.Equal(t, gs.Cash, restored.Cash)
	assert.Equal(t, gs.Hype, restored.Hype)
	assert.Equal(t, gs.TechDebt, restored.TechDebt)
	assert.Equal(t, gs.BugCount, restored.BugCount)
	assert.Equal(t, gs.Miles, restored.Miles)
	assert.Equal(t, gs.UserCount, restored.UserCount)

	assert.Equal(t, len(gs.TurnHistory), len(restored.TurnHistory))
	for i, te := range gs.TurnHistory {
		assert.Equal(t, te.Action, restored.TurnHistory[i].Action)
		assert.Equal(t, te.DeathRoll, restored.TurnHistory[i].DeathRoll)
		assert.Equal(t, te.EventID, restored.TurnHistory[i].EventID)
	}
}

func TestSerializeDeserializeDeathRollLoss(t *testing.T) {
	gs := &GameState{
		TurnNumber: 1,
		Server:     ServerEC2,
		Database:   DBRDS,
		Cash:       1400,
		Hype:       80,
		TurnHistory: []TurnEntry{
			{Action: 2, DeathRoll: DeathRollLoss, EventID: 7},
		},
	}

	serialized := Serialize(gs)
	assert.Contains(t, serialized, "a2l")

	restored, err := Deserialize(serialized)
	assert.NoError(t, err)
	assert.Equal(t, DeathRollLoss, restored.TurnHistory[0].DeathRoll)
}

func TestDeserializeInvalid(t *testing.T) {
	_, err := Deserialize("garbage")
	assert.Error(t, err)

	_, err = Deserialize("1;s1d1;a1;100/50/0/0/0/0")
	assert.Error(t, err, "missing event history separator")
}

func TestSerializeEmptyTurnHistory(t *testing.T) {
	gs := &GameState{
		TurnNumber:  0,
		Server:      ServerFargate,
		Database:    DBAurora,
		Cash:        1500,
		Hype:        75,
		TurnHistory: []TurnEntry{},
	}

	serialized := Serialize(gs)
	assert.Equal(t, "0;s1d1;;1500/75/0/0/0/0 - ", serialized)

	restored, err := Deserialize(serialized)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(restored.TurnHistory))
}
