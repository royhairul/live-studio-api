package entity

import (
	"gorm.io/gorm"

	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	attendanceentity "github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
)

type Accountsession struct {
	gorm.Model
	AccountID     uint
	Account       accountentity.Account `gorm:"foreignKey:AccountID;references:ID"`
	AttendanceID  uint
	Attendance    attendanceentity.Attendance `gorm:"foreignKey:AttendanceID;references:ID"`
	GMVSalesStart uint
	GMVSalesEnd   uint
	GMVPaidStart  uint
	GMVPaidEnd    uint
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
