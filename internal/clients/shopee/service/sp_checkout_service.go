package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
)

type ShopeeCheckoutService interface {
	GetAll(cookie string)
	GetByRangeTime(cookie string, startTime time.Time, endTime time.Time) (*params.ShopeeCheckoutResponse, error)
}
