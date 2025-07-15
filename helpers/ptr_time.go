package helpers

import "time"

func TimeNow() *time.Time {
	now := time.Now()
	return &now
}
