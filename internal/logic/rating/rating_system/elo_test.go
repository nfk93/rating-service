package rating_system

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateUpdatedEloRatings(t *testing.T) {
	r := eloRatingSystem{
		t: 400,
		k: 40,
	}

	tests := []struct {
		name string

		ratingDiff int
		draw       bool
		upset      bool

		expectedRatingChange int
	}{
		{
			name: "equal rating, draw",

			ratingDiff:           0,
			draw:                 true,
			expectedRatingChange: 0,
		},
		{
			name: "equal rating gives K/2 rating change",

			ratingDiff:           0,
			expectedRatingChange: int(r.k / 2),
		},
		{
			name: "T rating diff, expected result",

			ratingDiff:           int(r.t),
			expectedRatingChange: int(math.Floor(r.k / 11)),
		},
		{
			name: "T rating diff, upset result",

			ratingDiff:           int(r.t),
			expectedRatingChange: int(math.Floor(10 * r.k / 11)),
			upset:                true,
		},
		{
			name:                 "T rating diff, draw",
			ratingDiff:           int(r.t),
			draw:                 true,
			expectedRatingChange: int(math.Floor(4.5 * r.k / 11)),
		},
		{
			name: "2T rating diff, expected result",

			ratingDiff:           2 * int(r.t),
			expectedRatingChange: int(math.Floor(1 * r.k / 101)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				// A is the winning player
				oldRatingB := rand.Intn(4000)
				oldRatingA := oldRatingB + tt.ratingDiff
				if tt.upset {
					// swap ratings in case of upset
					temp := oldRatingA
					oldRatingA = oldRatingB
					oldRatingB = temp
				}

				var newRatingA, newRatingB int
				if tt.draw {
					newRatingA, newRatingB = r.calculateUpdatedEloRatings(oldRatingA, oldRatingB, 0.5)
				} else {
					newRatingA, newRatingB = r.calculateUpdatedEloRatings(oldRatingA, oldRatingB, 1.0)
					newRatingB_, newRatingA_ := r.calculateUpdatedEloRatings(oldRatingB, oldRatingA, 0.0)

					assert.Equal(t, newRatingA, newRatingA_, "failed rating(a, b, r)_0 == rating(b, a, 1-r)_1 check")
					assert.Equal(t, newRatingB, newRatingB_, "failed rating(a, b, r)_1 == rating(b, a, 1-r)_0 check")
				}
				assert.Equal(t, tt.expectedRatingChange, int(math.Abs(float64(oldRatingA-newRatingA))), "rating difference is wrong for player A")
				assert.True(t, newRatingA+newRatingB == oldRatingB+oldRatingA,
					"new ratings sum to different value than old ratings. New ratings sum: %v, old ratings sum: %v", newRatingA+newRatingB, oldRatingB+oldRatingA)
			}
		})
	}
}
