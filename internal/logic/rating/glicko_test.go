package rating

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlicko(t *testing.T) {
	player := &PlayerStats{
		PreviousPeriodRating: 1500,
		RatingDeviationOld:   200,
	}

	opponent1 := PlayerStats{
		PreviousPeriodRating: 1400,
		RatingDeviationOld:   30,
	}
	opponent2 := PlayerStats{
		PreviousPeriodRating: 1550,
		RatingDeviationOld:   100,
	}
	opponent3 := PlayerStats{
		PreviousPeriodRating: 1700,
		RatingDeviationOld:   300,
	}

	player.addMatch(&opponent1, Win)
	player.addMatch(&opponent2, Loss)
	player.addMatch(&opponent3, Loss)

	assert.True(t, math.Abs(player.CurrentRating-1464) <= 1, "Wrong rating, expected=~1464 actual=%v", player.CurrentRating)
}
