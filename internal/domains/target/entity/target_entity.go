package entity

import (
	"time"

	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"gorm.io/gorm"
)

type Target struct {
	gorm.Model
	Date         time.Time `gorm:"uniqueIndex:idx_studio_month"`
	TargetGMV    int64
	TargetIncome int64
	StudioID     uint                `gorm:"uniqueIndex:idx_studio_month"`
	Studio       studioentity.Studio `gorm:"foreignKey:StudioID;references:ID"`

	tenantdb.TenantBase
}

func (t *Target) BeforeSave(tx *gorm.DB) (err error) {
	year, month, _ := t.Date.Date()
	t.Date = time.Date(year, month, 1, 0, 0, 0, 0, t.Date.Location())
	return
}
