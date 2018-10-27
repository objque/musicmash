package cron

import (
	"math"
	"time"
)

func calcDiffHours(lastFetchDate time.Time) float64 {
	now := time.Now().UTC().Add(-time.Hour)

	// calc diff between hours
	num := now.Sub(lastFetchDate).Hours()
	diff, _ := math.Modf(num)
	return diff
}
