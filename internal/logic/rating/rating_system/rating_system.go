package rating_system

type RatingSystem interface {
	// Returns the updates to apply to each player
	// If winnerIndex is -1, the game is considered a draw
	RatingsDiffs(ratings []int, winnerIndex int) ([]int, error)

	DefaultRating() int
}
