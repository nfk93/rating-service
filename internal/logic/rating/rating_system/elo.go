package rating_system

import (
	"errors"
	"math"
)

const (
	// The closer to 100% likelihood of one player winning, the closer to K the
	// rating change will be in case of an upset result
	K float64 = 40 // TODO: perhaps vary based on rating

	// Every T rating difference, the amount of games it takes for the lower rated
	// player to have a 50% probability of winning one game is increased by about a
	// factor 10. At equal rating, the number is 1.
	T float64 = 400
)

type eloRatingSystem struct {
	defaultRating int

	// The closer to 100% likelihood of one player winning, the closer to k the
	// rating change will be in case of an upset result
	k float64

	// Every t rating difference, the amount of games it takes for the lower rated
	// player to have a 50% probability of winning one game is increased by about a
	// factor 10. At equal rating, the number is 1.
	t float64
}

func NewEloRatingSystem() RatingSystem {
	return &eloRatingSystem{
		1000,
		K,
		T,
	}
}

// If winnerIndex is -1, it is a draw
func (r *eloRatingSystem) RatingsDiffs(ratings []int, winnerIndex int) ([]int, error) {
	if len(ratings) != 2 {
		return nil, errors.New("elo ratingsystem only supports 1v1 matches")
	}
	if winnerIndex < -1 || winnerIndex > 1 {
		return nil, errors.New("invalid winner index")
	}

	result := float64(1 - winnerIndex)
	if winnerIndex == -1 {
		result = 0.5
	}

	r0, r1 := ratings[0], ratings[1]
	diff := r.calculateEloDiffs(r0, r1, result)
	return []int{r0 + diff, r1 + diff}, nil
}

func (r *eloRatingSystem) DefaultRating() int {
	return r.defaultRating
}

// Result must be one of the following:
// - 1.0 if player A won
// - 0.0 if player B won
// - 0.5 for a draw
// Any other value will result in an erroneous result
func (r *eloRatingSystem) calculateUpdatedEloRatings(ratingA, ratingB int, result float64) (int, int) {
	ratingDiff := r.calculateEloDiffs(ratingA, ratingB, result)

	ratingA += ratingDiff
	ratingB -= ratingDiff

	return ratingA, ratingB
}

func (r *eloRatingSystem) calculateEloDiffs(ratingA, ratingB int, result float64) int {
	E_A := 1 / (1 + math.Pow(10, float64(ratingB-ratingA)/r.t))
	ratingDiff := int(r.k * (result - E_A))

	return ratingDiff
}
