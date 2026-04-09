package svt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	// TODO: test if too low or too expensive to purchase

	t.Run("valid oxen purchase", func(t *testing.T) {
		gs := InitState()
		ok, _ := PurchaseItem(&gs, PhasePurchaseOxen, 250)

		assert.True(t, ok, "expected ok")
		assert.Equal(t, 250, gs.Inventory.Oxen, "oxen should be 250")
	})

	t.Run("valid ammo purchase", func(t *testing.T) {})

	t.Run("valid clothing purchase", func(t *testing.T) {})

	t.Run("valid food purchase", func(t *testing.T) {})

}
