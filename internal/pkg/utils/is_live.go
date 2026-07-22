package utils

import "time"

// LiveTolerance is how far behind "now" a session's reported end may fall while
// still counting as on air. Shopee's realtime list is polled every few seconds
// and its duration advances with the stream, so a live session lands within a
// few seconds of now; the margin absorbs poll lag and clock skew.
const LiveTolerance = 5 * time.Minute

// IsLive reports whether a live session is still on air, given its start time in
// Unix epoch milliseconds and Shopee's reported duration in milliseconds.
//
// Shopee's `status` field cannot be used for this — it reads 2 for both ongoing
// and ended sessions. What does distinguish them is that duration keeps
// advancing while a stream runs, so startTime+duration tracks now until it ends
// and then freezes.
func IsLive(startTime, duration int64) bool {
	if startTime <= 0 {
		return false
	}

	endsAt := time.UnixMilli(startTime + duration)
	now := time.Now()

	// A session whose end is in the future is still running; one that ended
	// within the tolerance is treated as live to absorb polling lag.
	return now.Sub(endsAt) <= LiveTolerance
}
