package entity

import (
	"gorm.io/gorm"

	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	attendanceentity "github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

type Accountsession struct {
	gorm.Model
	GMVSalesStart uint
	GMVSalesEnd   uint
	GMVPaidStart  uint
	GMVPaidEnd    uint
	AccountID     uint
	Account       accountentity.Account `gorm:"foreignKey:AccountID;references:ID"`
	AttendanceID  uint
	Attendance    attendanceentity.Attendance `gorm:"foreignKey:AttendanceID;references:ID"`
	StudioID      uint
	Studio        studioentity.Studio `gorm:"foreignKey:StudioID;references:ID"`

	tenantdb.TenantBase
}

func (a *Accountsession) TotalPaid() uint {
	if a.GMVPaidEnd > a.GMVPaidStart {
		return a.GMVPaidEnd - a.GMVPaidStart
	}
	return 0
}

func (a Accountsession) TotalSales() uint {
	if a.GMVSalesEnd > a.GMVSalesStart {
		return a.GMVSalesEnd - a.GMVSalesStart
	}
	return 0
}
