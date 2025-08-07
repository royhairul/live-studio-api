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
