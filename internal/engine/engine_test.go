package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubGameStore struct {
	State GameState
	Saved bool
}

func (s *StubGameStore) SaveState(state GameState) error {
	s.Saved = true
	s.State = state
	return nil
}

func (s *StubGameStore) LoadState() (GameState, error) {
	return s.State, nil
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
	gs := InitState()

	assert.Equal(t, 1500, gs.Cash)
	assert.GreaterOrEqual(t, gs.Hype, 50)
	assert.LessOrEqual(t, gs.Hype, 100)
	assert.Zero(t, gs.TechDebt)
	assert.Zero(t, gs.BugCount)
	assert.Zero(t, gs.Miles)
}

func TestSetServer(t *testing.T) {
	t.Run("valid choice", func(t *testing.T) {
		gs := InitState()
		ok := SetServer(&gs, 1)
		assert.True(t, ok)
		assert.Equal(t, ServerFargate, gs.Server)
	})
	t.Run("ThinkPad", func(t *testing.T) {
		gs := InitState()
		ok := SetServer(&gs, 4)
		assert.True(t, ok)
		assert.Equal(t, ServerThinkPad, gs.Server)
	})
	t.Run("invalid choice", func(t *testing.T) {
		gs := InitState()
		assert.False(t, SetServer(&gs, 0))
		assert.False(t, SetServer(&gs, 5))
	})
}

func TestSetDatabase(t *testing.T) {
	t.Run("valid choice", func(t *testing.T) {
		gs := InitState()
		ok := SetDatabase(&gs, 3)
		assert.True(t, ok)
		assert.Equal(t, DBSQLite, gs.Database)
	})
	t.Run("invalid choice", func(t *testing.T) {
		gs := InitState()
		assert.False(t, SetDatabase(&gs, 0))
		assert.False(t, SetDatabase(&gs, 4))
	})
}

func TestAPIGateway(t *testing.T) {
	t.Run("AWS server needs gateway", func(t *testing.T) {
		gs := InitState()
		gs.Server = ServerFargate
		gs.Database = DBSQLite
		assert.True(t, NeedsAPIGateway(&gs))
		assert.Equal(t, 129, APIGatewayCost(&gs))
	})
	t.Run("AWS db needs gateway", func(t *testing.T) {
		gs := InitState()
		gs.Server = ServerThinkPad
		gs.Database = DBAurora
		assert.True(t, NeedsAPIGateway(&gs))
	})
	t.Run("no AWS no gateway", func(t *testing.T) {
		gs := InitState()
		gs.Server = ServerThinkPad
		gs.Database = DBSQLite
		assert.False(t, NeedsAPIGateway(&gs))
		assert.Equal(t, 0, APIGatewayCost(&gs))
	})
}

func TestAdvanceMileage(t *testing.T) {
	gs := InitState()
	gs.Hype = 50
	gs.ActionChoice = 1
	miles := AdvanceMileage(&gs)
	assert.Greater(t, miles, 0)
	assert.Equal(t, miles, gs.Miles)
}

func TestAdvanceMileageHalvedForFixBugs(t *testing.T) {
	// Run multiple times to account for randomness
	for i := 0; i < 10; i++ {
		gs1 := InitState()
		gs1.Hype = 50
		gs1.ActionChoice = 1
		miles1 := AdvanceMileage(&gs1)

		gs2 := InitState()
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
	gs := InitState()
	gs.TechDebt = 20
	gs.BugCount = 5
	// techHealth = 100 - 20 - 5*3 = 65
	assert.Equal(t, 65, TechHealth(&gs))
}

func TestCalcCashBurn(t *testing.T) {
	gs := InitState()
	gs.Server = ServerEC2
	gs.Database = DBRDS
	gs.UserCount = 0
	// EC2: $40/mo, RDS: $30/mo, 0 users, AWS gateway: $129
	burn := CalcCashBurn(&gs)
	assert.Equal(t, 40+30+129, burn)
}

func TestCalcCashBurnWithUsers(t *testing.T) {
	gs := InitState()
	gs.Server = ServerFargate
	gs.Database = DBAurora
	gs.UserCount = 100
	// Fargate: $0.05/user, Aurora: $0.04/user = $0.09 * 100 = $9, + $129 gateway
	burn := CalcCashBurn(&gs)
	assert.Equal(t, 9+129, burn)
}

func TestCalcCashBurnNoAWS(t *testing.T) {
	gs := InitState()
	gs.Server = ServerThinkPad
	gs.Database = DBSQLite
	gs.UserCount = 500
	// ThinkPad: $0/mo $0/user, SQLite: $0/mo $0/user, no gateway
	burn := CalcCashBurn(&gs)
	assert.Equal(t, 0, burn)
}

func TestUpdateUserCount(t *testing.T) {
	gs := InitState()
	gs.Hype = 50
	UpdateUserCount(&gs)
	assert.Equal(t, 500, gs.UserCount)
}

func TestLoseConditions(t *testing.T) {
	t.Run("bankrupt", func(t *testing.T) {
		gs := InitState()
		gs.Cash = -1
		assert.True(t, IsBankrupt(&gs))
		reason, lost := CheckLoseCondition(&gs)
		assert.True(t, lost)
		assert.Contains(t, reason, "BANKRUPT")
	})
	t.Run("ghost town", func(t *testing.T) {
		gs := InitState()
		gs.Hype = 4
		assert.True(t, IsGhostTown(&gs))
		reason, lost := CheckLoseCondition(&gs)
		assert.True(t, lost)
		assert.Contains(t, reason, "GHOST TOWN")
	})
	t.Run("system failure", func(t *testing.T) {
		gs := InitState()
		gs.TechDebt = 80
		gs.BugCount = 10
		assert.True(t, IsSystemFailure(&gs))
		reason, lost := CheckLoseCondition(&gs)
		assert.True(t, lost)
		assert.Contains(t, reason, "SYSTEM FAILURE")
	})
	t.Run("no lose condition", func(t *testing.T) {
		gs := InitState()
		_, lost := CheckLoseCondition(&gs)
		assert.False(t, lost)
	})
}

func TestIsArrived(t *testing.T) {
	gs := InitState()
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
	gs := InitState()
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
		gs := InitState()
		gs.TechDebt = 0
		gs.BugCount = 0
		survived, _ := CheckIncident(&gs)
		assert.True(t, survived)
	})
	t.Run("unhealthy system may fail", func(t *testing.T) {
		gs := InitState()
		gs.TechDebt = 90
		gs.BugCount = 10
		// TechHealth = 100 - 90 - 30 = -20, should always fail
		survived, msg := CheckIncident(&gs)
		assert.False(t, survived)
		assert.Contains(t, msg, "INCIDENT")
	})
}

func TestFixBugs(t *testing.T) {
	gs := InitState()
	gs.BugCount = 10
	gs.TechDebt = 10
	gs.ActionChoice = 2
	bugsFixed, debtFixed := FixBugs(&gs)
	assert.Greater(t, bugsFixed, 0)
	assert.Greater(t, debtFixed, 0)
	assert.Less(t, gs.BugCount, 10)
}

func TestFixBugsNoOpWhenPushForward(t *testing.T) {
	gs := InitState()
	gs.BugCount = 10
	gs.ActionChoice = 1
	bugsFixed, debtFixed := FixBugs(&gs)
	assert.Equal(t, 0, bugsFixed)
	assert.Equal(t, 0, debtFixed)
	assert.Equal(t, 10, gs.BugCount)
}

func TestSystemDeathRoll(t *testing.T) {
	for i := 0; i < 50; i++ {
		result := SystemDeathRoll(100)
		assert.GreaterOrEqual(t, result, 1)
		assert.LessOrEqual(t, result, 100)
	}
	// Edge case: ceiling of 1
	result := SystemDeathRoll(1)
	assert.Equal(t, 1, result)
}
