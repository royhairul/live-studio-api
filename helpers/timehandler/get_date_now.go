package timehandler

import "time"

func DateNow() *string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc).Format("2006-01-02")
	return &now
}
