package timehandler

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/pkg/constants"
)

func DateNow() *string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc).Format(constants.LayoutYYMMDD)
	return &now
}
