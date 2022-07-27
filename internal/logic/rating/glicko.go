// Based on http://www.glicko.net/glicko/glicko.pdf

package rating

import (
	"math"
	"time"
)

const (
	StartRating          int = 1500
	StartRatingDeviation int = 350

	RatingPeriod time.Duration = 7 * 24 * time.Hour

	// The c value in the rating deviation calculation
	// 350 = sqrt(50^2 + 52*c^2)
	RatingDeviationFactor float64 = 48.038

	q   float64 = 0.0057565 // ln(10)/400
	q2  float64 = q * q
	pi2 float64 = math.Pi * math.Pi
)

type PlayerStats struct {
	CurrentRating        float64
	PreviousPeriodRating float64
	RatingDeviationOld   float64 // Rating deviation from last period

	// current value of (d^2)^-1
	d2powm1 float64

	// current value of sum_{j=1}^m ( g(RD_j)*(s_j-E[s|r,r_j,RD_j]) )
	gameOutcomeSum float64

	LastMatchTime time.Time
}

type MatchResult uint8

const (
	Loss MatchResult = 0
	Win  MatchResult = 1
	Draw MatchResult = 2
)

func (ps *PlayerStats) addMatch(opponentStats *PlayerStats, result MatchResult) {
	g_RD_j := g_glicko(opponentStats.RatingDeviationOld)
	g_RD_j_pow2 := g_RD_j * g_RD_j

	r_0 := ps.PreviousPeriodRating
	r_j := opponentStats.PreviousPeriodRating
	RDpow2 := ps.RatingDeviationOld * ps.RatingDeviationOld

	resultProbability := e_glicko(r_0, r_j, g_RD_j)

	d2powm1Summand := q2 * g_RD_j_pow2 * resultProbability * (1 - resultProbability)
	gameOutcomeSummand := g_RD_j * (matchResultToSValue(result) - resultProbability)

	ps.d2powm1 += d2powm1Summand
	ps.gameOutcomeSum += gameOutcomeSummand

	ps.CurrentRating = r_0 + (q/(ps.d2powm1+1/RDpow2))*ps.gameOutcomeSum
}

// g(x) 1/sqrt(1 + 3*q^2*x^2/pi^2)
func g_glicko(rd float64) float64 {
	rd2 := rd * rd

	return 1 / math.Sqrt(1+(3*q2*rd2/pi2))
}

func e_glicko(r, r_j, g_RD_j float64) float64 {
	return 1 / (1 + math.Pow(10, -1*g_RD_j*(r-r_j)/400))
}

func matchResultToSValue(m MatchResult) float64 {
	switch m {
	case Loss:
		return 0
	case Win:
		return 1
	default: // draw
		return 0.5
	}
}
