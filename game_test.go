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
