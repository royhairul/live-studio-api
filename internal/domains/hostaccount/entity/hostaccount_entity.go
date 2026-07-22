package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	hostentity "github.com/royhairul/live-studio-api/internal/domains/host/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

// HostAccount assigns a Shopee account to a host for a span of time.
//
// It replaces attendance as the host↔account bridge for reporting. Attendance
// says who was physically present; this says who a session's numbers belong to,
// which is what the performa pages actually need.
//
// The assignment is versioned rather than a plain column on accounts: ValidTo is
// nil while an assignment is current, and is set when the account moves to
// another host. Past reports therefore keep pointing at whoever held the account
// at the time, instead of silently changing on every reassignment.
//
// Many-to-many in both directions: one host can hold several accounts, and one
// account can pass between hosts over time. Two hosts holding the same account
// in overlapping periods is not blocked by the schema — the service rejects it,
// since a live session cannot belong to two hosts at once.
type HostAccount struct {
	gorm.Model

	HostID uuid.UUID       `gorm:"type:uuid;not null;index"`
	Host   hostentity.Host `gorm:"foreignKey:HostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	AccountID uint                  `gorm:"not null;index"`
	Account   accountentity.Account `gorm:"foreignKey:AccountID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// ValidFrom is inclusive. ValidTo is exclusive and nil means "still active",
	// so an open assignment needs no maintenance as time passes.
	ValidFrom time.Time  `gorm:"not null;index"`
	ValidTo   *time.Time `gorm:"index"`

	Note string `gorm:"type:text"`

	tenantdb.TenantBase
}

// IsActiveAt reports whether the assignment covers t.
func (h *HostAccount) IsActiveAt(t time.Time) bool {
	if t.Before(h.ValidFrom) {
		return false
	}
	return h.ValidTo == nil || t.Before(*h.ValidTo)
}

// Overlaps reports whether the assignment intersects [start, end]. Used both to
// pick assignments for a report window and to reject conflicting ones.
func (h *HostAccount) Overlaps(start, end time.Time) bool {
	if h.ValidTo != nil && !h.ValidTo.After(start) {
		return false
	}
	return !h.ValidFrom.After(end)
}
