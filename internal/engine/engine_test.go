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
		assert.Equal(t, 2040, TotalRequiredMileage, "total mileage is 2040")
	})
	t.Run("initial cash is 700", func(t *testing.T) {
		assert.Equal(t, 700, InitialCash, "initial cash is 700")
	})
}

func TestInitState(t *testing.T) {
	gs := InitState()

	assert.Equal(t, 700, gs.Player.Cash, "initial cash should be 700")
	assert.True(t, gs.Trip.FortAvailable, "initial fort should be available")
	assert.Zero(t, gs.Trip.Mileage, "initial mileage should be 0")
	assert.False(t, gs.Flags.Injured, "initial injury flag should be false")
}

func TestSetShootingLevel(t *testing.T) {
	t.Run("shooting level is valid", func(t *testing.T) {
		gs := InitState()
		ok := SetShootingLevel(&gs, 3)
		assert.True(t, ok, "expected ok")
		assert.Equal(t, 3, gs.Player.ShootingLevel, "shooting level should be 3")
	})
	t.Run("shooting level too low", func(t *testing.T) {
		gs := InitState()
		assert.False(t, SetShootingLevel(&gs, 0), "expected false")
	})
	t.Run("shooting level is too high", func(t *testing.T) {
		gs := InitState()
		assert.False(t, SetShootingLevel(&gs, 6), "expected false")
	})
}

func TestPurchaseItem(t *testing.T) {
	t.Run("valid oxen purchase", func(t *testing.T) {
		gs := InitState()
		ok, _ := PurchaseItem(&gs, PhasePurchaseOxen, 250)
		assert.True(t, ok, "expected ok")
		assert.Equal(t, 250, gs.Inventory.Oxen, "oxen should be 250")
	})

	t.Run("valid ammo purchase", func(t *testing.T) {
		gs := InitState()
		ok, _ := PurchaseItem(&gs, PhasePurchaseAmmo, 50)
		assert.True(t, ok, "expected ok")
		assert.Equal(t, 50, gs.Inventory.Ammo, "ammo should be 50")
	})

	t.Run("valid clothing purchase", func(t *testing.T) {
		gs := InitState()
		ok, _ := PurchaseItem(&gs, PhasePurchaseClothing, 55)
		assert.True(t, ok, "expected ok")
		assert.Equal(t, 55, gs.Inventory.Clothing, "clothing should be 55")
	})

	t.Run("valid food purchase", func(t *testing.T) {
		gs := InitState()
		ok, _ := PurchaseItem(&gs, PhasePurchaseFood, 199)
		assert.True(t, ok, "expected ok")
		assert.Equal(t, 199, gs.Inventory.Food, "food should be 199")
	})
}

func TestFinalizePurchase(t *testing.T) {
	t.Run("valid total", func(t *testing.T) {
		gs := InitState()
		gs.Inventory.Oxen = 200
		gs.Inventory.Food = 100
		gs.Inventory.Ammo = 50
		gs.Inventory.Clothing = 50
		gs.Inventory.Miscellaneous = 50
		ok, remaining := FinalizePurchases(&gs)

		assert.True(t, ok, "expected ok")
		assert.Equal(t, 250, remaining, "remaining cash should be 250")
		assert.Equal(t, 250, gs.Player.Cash, "player cash should be 250")
	})
	t.Run("overspent", func(t *testing.T) {
		gs := InitState()
		gs.Inventory.Oxen = 300
		gs.Inventory.Food = 200
		gs.Inventory.Ammo = 100
		gs.Inventory.Clothing = 100
		gs.Inventory.Miscellaneous = 100
		ok, _ := FinalizePurchases(&gs)

		assert.False(t, ok, "expected false for overspent")
	})
}

func TestApplyEating(t *testing.T) {
	gs := InitState()
	gs.Inventory.Food = 100
	ApplyEating(&gs, 2) // Moderately (2): 18 food eaten

	assert.Equal(t, 2, gs.Trip.EatingChoice, "eating choice should be 2")
	assert.Equal(t, 82, gs.Inventory.Food, "food should be 82")
}

func TestAdvanceMileage(t *testing.T) {
	gs := InitState()
	gs.Inventory.Oxen = 250
	gs.Trip.ActionChoice = 1
	AdvanceMileage(&gs)
	if gs.Trip.Mileage <= 0 {
		t.Errorf("mileage should be > 0, got %d", gs.Trip.Mileage)
	}
}

func TestGenerateEvent(t *testing.T) {
	gs := InitState()
	gs.Inventory.Food = 100
	gs.Inventory.Ammo = 50
	gs.Inventory.Miscellaneous = 30
	gs.Inventory.Clothing = 20
	gs.Trip.Mileage = 500

	originalFood := gs.Inventory.Food
	originalMileage := gs.Trip.Mileage
	changed := false
	for i := 0; i < 50; i++ {
		GenerateEvent(&gs)
		if gs.Inventory.Food != originalFood ||
			gs.Trip.Mileage != originalMileage ||
			gs.Flags.Injured || gs.Flags.Ill {
			changed = true
			break
		}
	}
	if !changed {
		t.Error("expected at least one event to change game state")
	}
}

func TestHandleAilment(t *testing.T) {
	t.Run("dies when no supplies", func(t *testing.T) {
		gs := InitState()
		gs.Flags.Ill = true
		gs.Inventory.Miscellaneous = 2
		survived, msg := HandleAilment(&gs)
		if survived {
			t.Error("expected death")
		}
		if msg != "YOU DIED OF PNEUMONIA." {
			t.Errorf("got %q", msg)
		}
	})
	t.Run("survives with supplies", func(t *testing.T) {
		gs := InitState()
		gs.Flags.Injured = true
		gs.Inventory.Miscellaneous = 30
		survived, _ := HandleAilment(&gs)
		if !survived {
			t.Error("expected survival")
		}
		if gs.Flags.Injured {
			t.Error("Injured should be cleared")
		}
	})
}

func TestDateName(t *testing.T) {
	got := DateName(1)
	if got != "SATURDAY, MARCH 29, 1847" {
		t.Errorf("got %q", got)
	}
	got = DateName(21)
	if got != "WINTER" {
		t.Errorf("got %q, want WINTER", got)
	}
}

func TestIsStarved(t *testing.T) {
	gs := InitState()
	gs.Inventory.Food = -1
	if !IsStarved(&gs) {
		t.Error("expected starved")
	}
}

func TestIsArrived(t *testing.T) {
	gs := InitState()
	gs.Trip.Mileage = 2040
	if !IsArrived(&gs) {
		t.Error("expected arrived")
	}
}
