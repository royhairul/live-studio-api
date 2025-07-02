package entity

import (
	"time"

	"github.com/google/uuid"
)

type Live struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	SessionID       int64     `gorm:"unique;not null"`
	Title           string
	StartTime       time.Time
	Duration        int
	Views           int
	PeakViewers     int
	AvgViewDuration float64
	Comments        int
	Likes           int
	FollowersGrowth int
	PlacedOrders    int
	PlacedSales     float64
	ConfirmedOrders int
	ConfirmedSales  float64
	HostID          uuid.UUID
	StudioID        uuid.UUID
	AccountID       uuid.UUID
}
