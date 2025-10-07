package timehandler

import "time"

func TimeNow() *time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	return &now
}
