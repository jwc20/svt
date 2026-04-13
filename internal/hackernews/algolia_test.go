package hackernews

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcBonusHype(t *testing.T) {
	tests := []struct {
		name     string
		karma    int
		nbHits   int
		expected int
	}{
		{"zero karma, zero hits", 0, 0, 0},
		{"high karma and hits caps at 100", 1000000, 50, 100},
		{"karma 1, 0 hits", 1, 0, 0},  // floor(log10(2)) * 5 = 0*5 = 0
		{"karma 0, 5 hits", 0, 5, 10}, // floor(log10(1)) * 5 + 5*2 = 0 + 10
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalcBonusHype(tt.karma, tt.nbHits)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFetchBonusHype_EmptyUsername(t *testing.T) {
	assert.Equal(t, 0, FetchBonusHype(""))
}
